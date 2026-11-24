package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Kind string

const (
	KindText         Kind = "text"
	KindInt          Kind = "integer"
	KindDecimal      Kind = "decimal"
	KindBool         Kind = "boolean"
	KindDate         Kind = "date"
	KindURL          Kind = "URL"
	KindUUID         Kind = "UUID"
	KindPicture      Kind = "picture"
	KindBinary       Kind = "binary"
	KindCredit       Kind = "credit"
	KindChapter      Kind = "chapter"
	KindSyncedLyrics Kind = "synchronized lyrics"
	KindTimedEvent   Kind = "timed event"
)

var kindConstraints = map[Kind]string{
	KindText:         "UTF-8 without U+0000; an empty string is a value.",
	KindInt:          "Nonnegative signed 64-bit integer.",
	KindDecimal:      "Finite 64-bit decimal value.",
	KindBool:         "True or false.",
	KindDate:         "Calendar-valid ISO year, month, date or timestamp at supplied precision.",
	KindURL:          "Absolute URI with scheme; HTTP and HTTPS require a host.",
	KindUUID:         "Hyphenated hexadecimal UUID.",
	KindPicture:      "Picture role 0–20, MIME type or linked image marker, and original bytes.",
	KindBinary:       "Opaque nonempty byte sequence.",
	KindCredit:       "Nonempty person name and role or instrument.",
	KindChapter:      "Chapter ID, title and millisecond start, with an optional millisecond end.",
	KindSyncedLyrics: "Language, description, and ordered millisecond lyric cues.",
	KindTimedEvent:   "Millisecond start, content type and opaque nonempty bytes.",
}

var (
	ErrInvalidField     = errors.New("invalid metadata field")
	ErrInvalidValue     = errors.New("invalid metadata value")
	ErrReadOnly         = errors.New("read-only metadata field")
	ErrUnsupported      = errors.New("unsupported metadata operation")
	ErrInvalidCustomKey = errors.New("invalid custom metadata key")
)

// Date preserves the original precision and timezone of an ISO timestamp.
type Date string

func ParseDate(s string) (Date, error) {
	for _, layout := range []string{
		"2006", "2006-01", "2006-01-02", "2006-01-02T15", "2006-01-02T15:04",
		"2006-01-02T15:04:05", time.RFC3339Nano,
	} {
		if len(s) != len(layout) && len(layout) < len("2006-01-02T15:04:05") {
			continue
		}
		if _, err := time.Parse(layout, s); err == nil {
			return Date(s), nil
		}
	}
	return "", fmt.Errorf("%w: invalid ISO date %q", ErrInvalidValue, s)
}

// Picture contains the original encoded image or, for MIME "-->", an absolute
// image URI. Its role uses ID3 APIC and FLAC PICTURE role numbers.
type Picture struct {
	Role                         uint32
	MIME                         string
	Description                  string
	Width, Height, Depth, Colors uint32
	Data                         []byte
}

type Credit struct {
	Name string
	Role string
}

type Chapter struct {
	ID      string
	Title   string
	StartMS int64
	EndMS   int64 // zero indicates an unspecified end
}

type LyricCue struct {
	TimeMS int64
	Text   string
}

type SyncedLyrics struct {
	Language    string
	Description string
	Cues        []LyricCue
}

type TimedEvent struct {
	StartMS     int64
	ContentType string
	Data        []byte
}

// Value is an immutable tagged metadata value. Constructors and getters copy
// mutable data. The zero value has no type and is invalid for a field.
type Value struct {
	kind    Kind
	text    string
	integer int64
	decimal float64
	payload string
	picture *Picture
}

func TextValue(s string) Value     { return Value{kind: KindText, text: s} }
func IntValue(n int64) Value       { return Value{kind: KindInt, integer: n} }
func DecimalValue(n float64) Value { return Value{kind: KindDecimal, decimal: n} }
func BoolValue(b bool) Value {
	if b {
		return Value{kind: KindBool, integer: 1}
	}
	return Value{kind: KindBool}
}
func DateValue(d Date) Value     { return Value{kind: KindDate, text: string(d)} }
func URLValue(s string) Value    { return Value{kind: KindURL, text: s} }
func UUIDValue(s string) Value   { return Value{kind: KindUUID, text: strings.ToLower(s)} }
func BinaryValue(b []byte) Value { return Value{kind: KindBinary, payload: string(b)} }
func PictureValue(p Picture) Value {
	p.Data = append([]byte(nil), p.Data...)
	return Value{kind: KindPicture, picture: &p}
}
func CreditValue(c Credit) Value {
	if !validText(c.Name) || !validText(c.Role) {
		return Value{kind: KindCredit}
	}
	return structuredValue(KindCredit, c)
}

func ChapterValue(c Chapter) Value {
	if !validText(c.ID) || !validText(c.Title) {
		return Value{kind: KindChapter}
	}
	return structuredValue(KindChapter, c)
}

func SyncedLyricsValue(l SyncedLyrics) Value {
	if !validText(l.Language) || !validText(l.Description) {
		return Value{kind: KindSyncedLyrics}
	}
	for _, cue := range l.Cues {
		if !validText(cue.Text) {
			return Value{kind: KindSyncedLyrics}
		}
	}
	return structuredValue(KindSyncedLyrics, l)
}

func TimedEventValue(e TimedEvent) Value {
	if !validText(e.ContentType) {
		return Value{kind: KindTimedEvent}
	}
	return structuredValue(KindTimedEvent, e)
}

type structured interface {
	Credit | Chapter | SyncedLyrics | TimedEvent
}

func structuredValue[T structured](kind Kind, data T) Value {
	// All five concrete types contain only JSON-compatible members.
	raw, _ := json.Marshal(data)
	return Value{kind: kind, payload: string(raw)}
}

// Equal reports whether two values represent the same canonical value.
func (v Value) Equal(other Value) bool {
	if v.picture == nil || other.picture == nil {
		return v == other
	}
	first, second := *v.picture, *other.picture
	return first.Role == second.Role && first.MIME == second.MIME &&
		first.Description == second.Description && first.Width == second.Width &&
		first.Height == second.Height && first.Colors == second.Colors &&
		first.Depth == second.Depth && bytes.Equal(first.Data, second.Data)
}

func (v Value) Kind() Kind               { return v.kind }
func (v Value) Text() (string, bool)     { return v.text, v.kind == KindText }
func (v Value) Int() (int64, bool)       { return v.integer, v.kind == KindInt }
func (v Value) Decimal() (float64, bool) { return v.decimal, v.kind == KindDecimal }
func (v Value) Bool() (bool, bool)       { return v.integer != 0, v.kind == KindBool }
func (v Value) Date() (Date, bool)       { return Date(v.text), v.kind == KindDate }
func (v Value) URL() (string, bool)      { return v.text, v.kind == KindURL }
func (v Value) UUID() (string, bool)     { return v.text, v.kind == KindUUID }
func (v Value) Binary() ([]byte, bool) {
	if v.kind != KindBinary {
		return nil, false
	}
	return []byte(v.payload), true
}
func (v Value) Picture() (Picture, bool) {
	if v.kind != KindPicture || v.picture == nil {
		return Picture{}, false
	}
	p := *v.picture
	p.Data = append([]byte(nil), p.Data...)
	return p, true
}
func (v Value) Credit() (Credit, bool)   { return getStructured[Credit](v, KindCredit) }
func (v Value) Chapter() (Chapter, bool) { return getStructured[Chapter](v, KindChapter) }
func (v Value) SyncedLyrics() (SyncedLyrics, bool) {
	return getStructured[SyncedLyrics](v, KindSyncedLyrics)
}
func (v Value) TimedEvent() (TimedEvent, bool) { return getStructured[TimedEvent](v, KindTimedEvent) }

func getStructured[T structured](v Value, kind Kind) (T, bool) {
	var output T
	if v.kind != kind {
		return output, false
	}
	err := json.Unmarshal([]byte(v.payload), &output)
	return output, err == nil
}

var (
	isrcFormat    = regexp.MustCompile("^[A-Z]{2}[A-Z0-9]{3}[0-9]{7}$")
	uuidFormat    = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
	upcFormat     = regexp.MustCompile("^[0-9]{12}$")
	barcodeFormat = regexp.MustCompile("^([0-9]{8}|[0-9]{12,14})$")
	asinFormat    = regexp.MustCompile("^[A-Z0-9]{10}$")
	customFormat  = regexp.MustCompile("^[A-Z0-9][A-Z0-9_.-]{0,127}$")
)

func validText(s string) bool { return utf8.ValidString(s) && !strings.ContainsRune(s, 0) }

func validURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.IsAbs() && validText(s) &&
		((u.Scheme != "http" && u.Scheme != "https") || u.Host != "")
}

func validPictureURI(s string) bool {
	_, err := url.Parse(s)
	return s != "" && err == nil && validText(s) && !strings.ContainsRune(s, ' ')
}

func NormalizeCustomKey(name string) (string, error) {
	for i := range name {
		if name[i] >= utf8.RuneSelf {
			return "", fmt.Errorf("%w: %q", ErrInvalidCustomKey, name)
		}
	}
	key := strings.ToUpper(name)
	if !customFormat.MatchString(key) {
		return "", fmt.Errorf("%w: %q", ErrInvalidCustomKey, name)
	}
	return key, nil
}

func Validate(field Field, values ...Value) error {
	d, ok := Lookup(field)
	if !ok {
		return fmt.Errorf("%w: %s", ErrInvalidField, field)
	}
	if len(values) > 1 && !d.Multiple {
		return fmt.Errorf("%w: %s accepts one value", ErrInvalidValue, field)
	}
	for _, v := range values {
		if v.kind != d.Kind || !valueValid(field, v) {
			return fmt.Errorf("%w: %s requires %s", ErrInvalidValue, field, d.Constraints)
		}
	}
	return nil
}

func valueValid(field Field, v Value) bool {
	switch v.kind {
	case KindText:
		if !validText(v.text) {
			return false
		}
		switch field {
		case ISRC:
			return isrcFormat.MatchString(v.text)
		case UPC:
			return upcFormat.MatchString(v.text)
		case Barcode:
			return barcodeFormat.MatchString(v.text)
		case ASIN:
			return asinFormat.MatchString(v.text)
		}
		return true
	case KindInt:
		switch field {
		case R128TrackGain, R128AlbumGain:
			return v.integer >= -32768 && v.integer <= 32767
		case Rating:
			return v.integer >= 0 && v.integer <= 100
		}
		return v.integer >= 0
	case KindDecimal:
		if math.IsInf(v.decimal, 0) || math.IsNaN(v.decimal) {
			return false
		}
		switch field {
		case ReplayGainTrackPeak, ReplayGainAlbumPeak, LoudnessRange:
			return v.decimal >= 0
		}
		return true
	case KindBool:
		return true
	case KindDate:
		_, err := ParseDate(v.text)
		return err == nil
	case KindURL:
		return validURL(v.text)
	case KindUUID:
		return uuidFormat.MatchString(v.text)
	case KindBinary:
		return len(v.payload) > 0
	case KindPicture:
		p := v.picture
		return p != nil && p.Role <= 20 && validText(p.Description) &&
			((p.MIME == "-->" && validPictureURI(string(p.Data))) ||
				(strings.HasPrefix(p.MIME, "image/") && validText(p.MIME) && len(p.Data) > 0))
	case KindCredit:
		c, ok := v.Credit()
		return ok && c.Name != "" && c.Role != "" && validText(c.Name) && validText(c.Role)
	case KindChapter:
		c, ok := v.Chapter()
		return ok && c.ID != "" && validText(c.ID) && validText(c.Title) &&
			c.StartMS >= 0 && (c.EndMS == 0 || c.EndMS > c.StartMS)
	case KindSyncedLyrics:
		l, ok := v.SyncedLyrics()
		if !ok || !validText(l.Language) || !validText(l.Description) || len(l.Cues) == 0 {
			return false
		}
		var previous int64
		for _, cue := range l.Cues {
			if cue.TimeMS < previous || !validText(cue.Text) {
				return false
			}
			previous = cue.TimeMS
		}
		return true
	case KindTimedEvent:
		e, ok := v.TimedEvent()
		return ok && e.StartMS >= 0 && e.ContentType != "" && validText(e.ContentType) && len(e.Data) > 0
	}
	return false
}

// Encode serializes a value for an adapter. Binary and structured fields use
// unpadded URL-safe base64. Structured payloads contain JSON for the field kind.
func Encode(field Field, v Value) (string, error) {
	if err := Validate(field, v); err != nil {
		return "", err
	}
	switch v.kind {
	case KindText, KindDate, KindURL, KindUUID:
		return v.text, nil
	case KindInt:
		return strconv.FormatInt(v.integer, 10), nil
	case KindBool:
		return strconv.FormatBool(v.integer != 0), nil
	case KindDecimal:
		return strconv.FormatFloat(v.decimal, 'f', -1, 64), nil
	case KindPicture:
		raw, err := json.Marshal(v.picture)
		if err != nil {
			return "", err
		}
		return base64.RawURLEncoding.EncodeToString(raw), nil
	default:
		return base64.RawURLEncoding.EncodeToString([]byte(v.payload)), nil
	}
}

// Decode parses the canonical wire value. Adapters must convert native syntax
// (for example, track pairs and ReplayGain units) before calling Decode.
func Decode(field Field, s string) (Value, error) {
	d, ok := Lookup(field)
	if !ok {
		return Value{}, fmt.Errorf("%w: %s", ErrInvalidField, field)
	}
	var v Value
	var err error
	switch d.Kind {
	case KindText:
		v = TextValue(s)
	case KindDate:
		v = DateValue(Date(s))
	case KindURL:
		v = URLValue(s)
	case KindUUID:
		v = UUIDValue(s)
	case KindInt:
		var n int64
		n, err = strconv.ParseInt(s, 10, 64)
		v = IntValue(n)
	case KindDecimal:
		var n float64
		n, err = strconv.ParseFloat(s, 64)
		v = DecimalValue(n)
	case KindBool:
		var b bool
		b, err = strconv.ParseBool(s)
		v = BoolValue(b)
	default:
		var data []byte
		data, err = base64.RawURLEncoding.DecodeString(s)
		v = Value{kind: d.Kind, payload: string(data)}
		if err == nil {
			v = normalizeStructured(v)
		}
	}
	if err != nil {
		return Value{}, fmt.Errorf("%w: %s: %v", ErrInvalidValue, field, err)
	}
	if err = Validate(field, v); err != nil {
		return Value{}, err
	}
	return v, nil
}

func normalizeStructured(v Value) Value {
	switch v.kind {
	case KindPicture:
		var p Picture
		if err := json.Unmarshal([]byte(v.payload), &p); err == nil {
			return Value{kind: KindPicture, picture: &p}
		}
	case KindCredit:
		if c, ok := v.Credit(); ok {
			return CreditValue(c)
		}
	case KindChapter:
		if c, ok := v.Chapter(); ok {
			return ChapterValue(c)
		}
	case KindSyncedLyrics:
		if l, ok := v.SyncedLyrics(); ok {
			return SyncedLyricsValue(l)
		}
	case KindTimedEvent:
		if e, ok := v.TimedEvent(); ok {
			return TimedEventValue(e)
		}
	}
	return v
}
