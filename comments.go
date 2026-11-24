package core

import (
	"fmt"
	"strconv"
	"strings"
)

// ReadComments interprets an ordered FLAC, Vorbis or Opus comment list. Native
// pictures are decoded by the adapter. Entries that cannot be interpreted
// remain part of the source and can be preserved by EditComments.
func ReadComments(format Format, entries []string) (map[Field][]Value, map[string][]string) {
	index := make(map[string][]string)
	custom := make(map[string][]string)
	for _, entry := range entries {
		key, s, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		key = strings.ToUpper(key)
		index[key] = append(index[key], s)
		if k, ok := commentCustomKey(key); ok && validText(s) {
			custom[k] = append(custom[k], s)
		}
	}
	fields := make(map[Field][]Value)
	for _, plan := range commentLayouts[format].ordered {
		d := plan.definition
		for _, key := range plan.keys {
			for _, s := range index[key] {
				if v, err := decodeComment(d.ID, key, s); err == nil {
					fields[d.ID] = append(fields[d.ID], v)
					if !d.Multiple {
						break
					}
				}
			}
			if len(fields[d.ID]) > 0 {
				break
			}
		}
	}
	for _, pair := range commentPairs {
		if len(fields[pair.total]) != 0 {
			continue
		}
		for _, s := range index[pair.key] {
			if _, right, ok := strings.Cut(s, "/"); ok {
				if v, err := Decode(pair.total, right); err == nil {
					fields[pair.total] = []Value{v}
					break
				}
			}
		}
	}
	return fields, custom
}

type commentPair struct {
	number, total Field
	key           string
}

var commentPairs = []commentPair{{TrackNumber, TrackTotal, "TRACKNUMBER"}, {DiscNumber, DiscTotal, "DISCNUMBER"}}

type commentField struct {
	definition Definition
	mapping    Representation
	keys       []string
}

type commentLayout struct {
	ordered []commentField
	byField map[Field]commentField
	byKey   map[string]Field
}

var commentLayouts = func() map[Format]commentLayout {
	layouts := make(map[Format]commentLayout, 3)
	for _, format := range []Format{FLAC, Vorbis, Opus} {
		layout := commentLayout{byField: make(map[Field]commentField), byKey: make(map[string]Field)}
		for _, d := range definitions {
			m := d.Mapping(format)
			if d.ID == Artwork || m.Key == "" || m.Class == ReadOnly {
				continue
			}
			p := commentField{definition: d, mapping: m, keys: commentKeys(m, d.ID)}
			layout.ordered = append(layout.ordered, p)
			layout.byField[d.ID] = p
			for _, key := range p.keys {
				layout.byKey[key] = d.ID
			}
		}
		layouts[format] = layout
	}
	return layouts
}()

func commentCustomKey(key string) (string, bool) {
	if !strings.HasPrefix(key, commentCustomPrefix) {
		return "", false
	}
	custom, err := NormalizeCustomKey(strings.TrimPrefix(key, commentCustomPrefix))
	return custom, err == nil
}

func commentKeys(m Representation, field Field) []string {
	keys := append([]string{m.Key}, m.Aliases...)
	fallback := CommentFallbackName(field)
	if m.Key != fallback {
		keys = append(keys, fallback)
	}
	return keys
}

func decodeComment(field Field, key, s string) (Value, error) {
	switch field {
	case TrackNumber, DiscNumber:
		if key == "TRACKNUMBER" || key == "DISCNUMBER" {
			s, _, _ = strings.Cut(s, "/")
		}
	case ReplayGainTrackGain, ReplayGainAlbumGain, ReplayGainReferenceLoudness:
		s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), " dB"))
	case R128TrackGain, R128AlbumGain:
		if !strings.HasPrefix(key, "AUDIOMETA_") && len(s) > 6 {
			return Value{}, fmt.Errorf("%w: Opus R128 integer exceeds six characters", ErrInvalidValue)
		}
	case Rating:
		if key == "RATING" {
			return decodeCommentRating(s)
		}
	}
	return Decode(field, s)
}

func decodeCommentRating(s string) (Value, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "+")
	whole, fraction, _ := strings.Cut(s, ".")
	whole = strings.TrimLeft(whole, "0")
	fraction = strings.TrimRight(fraction, "0")
	if s == "" || s == "." || len(fraction) > 2 {
		return Value{}, ErrInvalidValue
	}
	if whole == "1" && fraction == "" {
		return IntValue(100), nil
	}
	if whole != "" {
		return Value{}, ErrInvalidValue
	}
	var score int64
	for i, digit := range fraction {
		if digit < '0' || digit > '9' {
			return Value{}, ErrInvalidValue
		}
		place := int64(10)
		if i == 1 {
			place = 1
		}
		score += int64(digit-'0') * place
	}
	return IntValue(score), nil
}

func encodeComment(field Field, v Value, convention bool) (string, error) {
	s, err := Encode(field, v)
	if err != nil || !convention {
		return s, err
	}
	switch field {
	case ReplayGainTrackGain, ReplayGainAlbumGain, ReplayGainReferenceLoudness:
		s += " dB"
	case Rating:
		rating, _ := v.Int()
		s = strconv.FormatFloat(float64(rating)/100, 'f', 2, 64)
	default:
		if b, ok := v.Bool(); ok {
			if b {
				s = "1"
			} else {
				s = "0"
			}
		}
	}
	return s, nil
}

// EditComments preserves ordering, key spelling, opaque entries and untreated
// aliases. Setting a field replaces its preferred key. Removing a field
// removes its recognized encodings. Structural comment keys remain untouched.
func EditComments(format Format, entries []string, changes Changes) ([]string, error) {
	if changes.Empty() {
		return append([]string(nil), entries...), nil
	}
	layout, supported := commentLayouts[format]
	if !supported {
		return nil, fmt.Errorf("%w: comment format %s", ErrUnsupported, format)
	}

	fields := cloneMap(changes.Fields)
	// Common players write TRACKNUMBER=number/total. Separate the components
	// when an edit would otherwise lose the untouched total or retain a removed one.
	out := append([]string(nil), entries...)
	for _, pair := range commentPairs {
		_, editNumber := fields[pair.number]
		total, editTotal := fields[pair.total]
		if editNumber && !editTotal {
			original, _ := ReadComments(format, entries)
			existing := original[pair.total]
			totalKeys := layout.byField[pair.total].keys
			hasStandalone := false
			needsMigration := false
			for _, entry := range entries {
				key, raw, ok := strings.Cut(entry, "=")
				if !ok {
					continue
				}
				for _, wanted := range totalKeys {
					if strings.EqualFold(wanted, key) {
						if _, err := Decode(pair.total, raw); err == nil {
							hasStandalone = true
						}
					}
				}
				if !strings.EqualFold(key, pair.key) {
					continue
				}
				if _, err := decodeComment(pair.number, pair.key, raw); err != nil {
					continue
				}
				_, suffix, paired := strings.Cut(raw, "/")
				if !paired {
					continue
				}
				component, err := Decode(pair.total, suffix)
				if err != nil || len(existing) == 0 || !component.Equal(existing[0]) {
					return nil, fmt.Errorf("%w: edit %s and %s together to resolve a conflicting or opaque paired total", ErrUnsupported, pair.number, pair.total)
				}
				needsMigration = true
			}
			if !hasStandalone && needsMigration {
				fields[pair.total] = existing
			}
		}
		if editTotal && len(total) == 0 {
			for i, entry := range out {
				key, s, _ := strings.Cut(entry, "=")
				if !strings.EqualFold(key, pair.key) {
					continue
				}
				if left, right, ok := strings.Cut(s, "/"); ok {
					if _, err := Decode(pair.total, right); err == nil {
						out[i] = key + "=" + left
					}
				}
			}
		}
	}

	replacements := make(map[string][]string)
	ordered := make([]string, 0, len(fields)+len(changes.Custom))
	for _, p := range layout.ordered {
		d, m := p.definition, p.mapping
		keys := p.keys
		values, edited := fields[d.ID]
		if !edited {
			continue
		}
		if len(values) == 0 {
			for _, key := range keys {
				replacements[key] = nil
			}
			continue
		}
		var items []string
		for _, v := range values {
			s, err := encodeComment(d.ID, v, m.Class != Fallback)
			if err != nil {
				return nil, err
			}
			items = append(items, m.Key+"="+s)
		}
		replacements[m.Key] = items
		ordered = append(ordered, m.Key)
	}
	for _, k := range SortedCustomKeys(changes.Custom) {
		key := commentCustomPrefix + k
		var items []string
		for _, s := range changes.Custom[k] {
			items = append(items, key+"="+s)
		}
		replacements[key] = items
		ordered = append(ordered, key)
	}
	seen := make(map[string]bool)
	result := make([]string, 0, len(out)+len(replacements))
	for _, entry := range out {
		key, s, ok := strings.Cut(entry, "=")
		key = strings.ToUpper(key)
		if !ok || key == "" {
			if !changes.ClearUnknown {
				result = append(result, entry)
			}
			continue
		}
		replacement, edited := replacements[key]
		if !edited && !changes.ClearUnknown {
			result = append(result, entry)
			continue
		}
		field, isKnown := layout.byKey[key]
		if isKnown {
			_, err := decodeComment(field, key, s)
			if err != nil && edited && len(replacement) > 0 && !changes.ClearUnknown &&
				(key == "R128_TRACK_GAIN" || key == "R128_ALBUM_GAIN") {
				return nil, fmt.Errorf("%w: an opaque value occupies the single-valued %s key", ErrUnsupported, key)
			}
			isKnown = err == nil
		}
		_, isCustom := commentCustomKey(key)
		isCustom = isCustom && validText(s)
		if edited && (isKnown || isCustom) {
			if !seen[key] {
				result = append(result, replacement...)
				seen[key] = true
			}
			continue
		}
		if !changes.ClearUnknown || isKnown || isCustom ||
			key == "WAVEFORMATEXTENSIBLE_CHANNEL_MASK" || (format != FLAC && key == "METADATA_BLOCK_PICTURE") {
			result = append(result, entry)
		}
	}
	for _, key := range ordered {
		if !seen[key] {
			result = append(result, replacements[key]...)
		}
	}
	return result, nil
}
