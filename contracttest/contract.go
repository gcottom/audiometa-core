// Package contracttest exercises the editor contract on real media from any adapter.
package contracttest

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"testing"

	"github.com/gcottom/audiometa-core"
)

type Open func(io.ReadSeeker) (*core.Editor, error)

// Sample provides reproducible, type-correct values for every registry field.
func Sample(d core.Definition) []core.Value {
	first := value(d, 0)
	if !d.Multiple {
		return []core.Value{first}
	}
	return []core.Value{first, value(d, 1)}
}

func value(d core.Definition, n int) core.Value {
	text := []string{"確認 — Áudio", "Segundo"}
	switch d.ID {
	case core.ISRC:
		return core.TextValue([]string{"USABC2600001", "USABC2600002"}[n])
	case core.ASIN:
		return core.TextValue("B000000001")
	case core.UPC, core.Barcode:
		return core.TextValue("012345678901")
	case core.Language:
		return core.TextValue([]string{"eng", "spa"}[n])
	case core.MusicalKey:
		return core.TextValue("Am")
	case core.R128TrackGain, core.R128AlbumGain:
		return core.IntValue(-1024)
	case core.BPM:
		return core.IntValue(123)
	}
	switch d.Kind {
	case core.KindText:
		return core.TextValue(text[n])
	case core.KindInt:
		return core.IntValue(int64(n + 7))
	case core.KindDecimal:
		return core.DecimalValue(1.25)
	case core.KindBool:
		return core.BoolValue(false)
	case core.KindDate:
		return core.DateValue("2024-02")
	case core.KindURL:
		return core.URLValue("https://example.org/" + []string{"first", "second"}[n])
	case core.KindUUID:
		return core.UUIDValue("27b771fd-f244-4341-9ff4-79c7dff5600" + []string{"1", "2"}[n])
	case core.KindBinary:
		return core.BinaryValue([]byte{0, 255, 1})
	case core.KindCredit:
		return core.CreditValue(core.Credit{Name: text[n], Role: "piano"})
	case core.KindChapter:
		return core.ChapterValue(core.Chapter{ID: []string{"one", "two"}[n], Title: text[n], StartMS: int64(n) * 100, EndMS: int64(n+1) * 100})
	case core.KindSyncedLyrics:
		return core.SyncedLyricsValue(core.SyncedLyrics{Language: "eng", Description: "verse", Cues: []core.LyricCue{{TimeMS: 100, Text: text[n]}}})
	case core.KindTimedEvent:
		return core.TimedEventValue(core.TimedEvent{StartMS: int64(n) * 10, ContentType: "application/octet-stream", Data: []byte{0, 255}})
	case core.KindPicture:
		img, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQImWNgYGD4DwABBAEAX+XDSwAAAABJRU5ErkJggg==")
		return core.PictureValue(core.Picture{Role: 3 + uint32(n), MIME: "image/png", Description: text[n], Width: 1, Height: 1, Depth: 32, Data: img})
	default:
		panic("missing sample for " + d.Kind)
	}
}

func AssertValues(t testing.TB, got, want []core.Value) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("metadata cardinality = %d, want %d", len(got), len(want))
	}
	for i, v := range want {
		if !v.Equal(got[i]) {
			t.Fatalf("value %d = %#v, want %#v", i, got[i], v)
		}
	}
}

func OpenFixture(t testing.TB, open Open, input []byte) *core.Editor {
	t.Helper()
	tag, err := open(bytes.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	return tag
}

func SaveFixture(t testing.TB, tag *core.Editor) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := tag.Save(&buf); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// Run verifies field-by-field round trips and removal plus custom values,
// no-op byte preservation and repeatable saves.
func Run(t *testing.T, open Open, input []byte) {
	t.Helper()
	t.Run("unmodified", func(t *testing.T) {
		tag := OpenFixture(t, open, input)
		if out := SaveFixture(t, tag); !bytes.Equal(out, input) {
			t.Fatal("unmodified bytes changed")
		}
	})
	for _, d := range core.Fields() {
		t.Run(string(d.ID), func(t *testing.T) {
			tag := OpenFixture(t, open, input)
			want := Sample(d)
			err := tag.SetField(d.ID, want...)
			if !d.Mutable {
				if !errors.Is(err, core.ErrReadOnly) {
					t.Fatalf("read-only set returned %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			out := SaveFixture(t, tag)
			got := OpenFixture(t, open, out)
			gotValues, _ := got.Get(d.ID)
			AssertValues(t, gotValues, want)
			if err := got.RemoveField(d.ID); err != nil {
				t.Fatal(err)
			}
			cleared := OpenFixture(t, open, SaveFixture(t, got))
			if _, ok := cleared.Get(d.ID); ok {
				t.Fatal("removed field remains present")
			}
		})
	}
	t.Run("custom_and_repeated_save", func(t *testing.T) {
		tag := OpenFixture(t, open, input)
		if err := tag.SetCustom("collection", "", "音楽", "音楽"); err != nil {
			t.Fatal(err)
		}
		if err := tag.RemoveCustomValue("COLLECTION", "音楽"); err != nil {
			t.Fatal(err)
		}
		if err := tag.AppendCustom("Collection", "second"); err != nil {
			t.Fatal(err)
		}
		first, second := SaveFixture(t, tag), SaveFixture(t, tag)
		if !bytes.Equal(first, second) {
			t.Fatal("repeated saves differ")
		}
		got := OpenFixture(t, open, first)
		values, ok := got.GetCustom("COLLECTION")
		if !ok || len(values) != 3 || values[0] != "" || values[1] != "音楽" || values[2] != "second" {
			t.Fatalf("custom round trip = %#v", values)
		}
		got.ClearCustom()
		if reopened := OpenFixture(t, open, SaveFixture(t, got)); len(reopened.Custom()) != 0 {
			t.Fatal("custom values remain")
		}
	})
}
