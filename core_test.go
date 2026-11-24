package core_test

import (
	"errors"
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/gcottom/audiometa-core/v2"
	"github.com/gcottom/audiometa-core/v2/contracttest"
)

func TestRegistryValueEncoding(t *testing.T) {
	seen := map[core.Field]bool{}
	for _, d := range core.Fields() {
		t.Run(string(d.ID), func(t *testing.T) {
			if seen[d.ID] {
				t.Fatal("duplicate ID")
			}
			seen[d.ID] = true
			lookup, ok := core.Lookup(d.ID)
			if !ok || lookup != d || d.Name == "" || d.Meaning == "" || d.Constraints == "" {
				t.Fatal("incomplete registry definition")
			}
			for _, v := range contracttest.Sample(d) {
				raw, err := core.Encode(d.ID, v)
				if err != nil {
					t.Fatal(err)
				}
				decoded, err := core.Decode(d.ID, raw)
				if err != nil {
					t.Fatal(err)
				}
				if !decoded.Equal(v) {
					t.Fatal("encoding did not round trip")
				}
			}
		})
	}
}

func TestEditsAndOwnership(t *testing.T) {
	var last core.Changes
	tag := core.NewEditor(nil, nil, func(_ io.Writer, c core.Changes) error { last = c; return nil })
	if _, ok := tag.Get(core.Title); ok {
		t.Fatal("absent title present")
	}
	if err := tag.SetField(core.Title, core.TextValue("")); err != nil {
		t.Fatal(err)
	}
	if v, ok := tag.GetField(core.Title); !ok || !v.Equal(core.TextValue("")) {
		t.Fatal("empty title missing")
	}
	if err := tag.SetField(core.Rating, core.IntValue(0)); err != nil {
		t.Fatal(err)
	}
	if v, ok := tag.GetField(core.Rating); !ok || !v.Equal(core.IntValue(0)) {
		t.Fatal("zero rating missing")
	}
	if err := tag.SetField(core.Rating, core.IntValue(101)); !errors.Is(err, core.ErrInvalidValue) {
		t.Fatal(err)
	}
	if err := tag.SetField(core.ReplayGainTrackGain, core.DecimalValue(math.NaN())); !errors.Is(err, core.ErrInvalidValue) {
		t.Fatal(err)
	}
	if err := tag.SetField(core.R128AlbumGain, core.IntValue(-32769)); !errors.Is(err, core.ErrInvalidValue) {
		t.Fatal(err)
	}
	if err := tag.SetField(core.AudioDuration, core.IntValue(1000)); !errors.Is(err, core.ErrReadOnly) {
		t.Fatal(err)
	}
	if err := tag.SetField(core.Title, core.TextValue("x"), core.TextValue("y")); !errors.Is(err, core.ErrInvalidValue) {
		t.Fatal(err)
	}
	if err := tag.SetCustom("pretty:key", "x"); !errors.Is(err, core.ErrInvalidCustomKey) {
		t.Fatal(err)
	}
	if err := tag.SetCustom("Key", "a", "b", "a"); err != nil {
		t.Fatal(err)
	}
	if err := tag.RemoveCustomValue("key", "a"); err != nil {
		t.Fatal(err)
	}
	copied, _ := tag.GetCustom("key")
	copied[0] = "caller edit"
	if got := tag.Custom()["KEY"]; !reflect.DeepEqual(got, []string{"b", "a"}) {
		t.Fatal(got)
	}

	data := []byte("original")
	picture := core.PictureValue(core.Picture{MIME: "image/jpeg", Data: data})
	data[0] = '!'
	got, _ := picture.Picture()
	got.Data[0] = '!'
	again, _ := picture.Picture()
	if string(again.Data) != "original" {
		t.Fatal("picture ownership leaked")
	}
	if err := tag.SetField(core.Artwork, picture, picture); err != nil {
		t.Fatal(err)
	}
	if err := tag.RemoveValue(core.Artwork, core.PictureValue(again)); err != nil {
		t.Fatal(err)
	}
	if vals, _ := tag.Get(core.Artwork); len(vals) != 1 {
		t.Fatal("picture removal did not use semantic equality")
	}
	tag.ClearAll(true)
	if err := tag.Save(io.Discard); err != nil {
		t.Fatal(err)
	}
	if _, ok := last.Fields[core.Title]; !ok || last.Fields[core.Title] != nil || last.ClearUnknown || len(last.Custom) != 1 {
		t.Fatalf("unexpected changes %+v", last)
	}
}

func TestPartialDates(t *testing.T) {
	for _, candidate := range []string{"2024", "2024-02", "2024-02-29", "2024-02-29T14", "2024-02-29T14:30", "2024-02-29T14:30:00+05:30"} {
		if date, err := core.ParseDate(candidate); err != nil || string(date) != candidate {
			t.Fatalf("date %q: %q %v", candidate, date, err)
		}
	}
	for _, candidate := range []string{"2023-02-29", "2020-13", "2020-00", "2020-1", "196", "2024-02-29T25:00:00"} {
		if _, err := core.ParseDate(candidate); !errors.Is(err, core.ErrInvalidValue) {
			t.Fatalf("accepted invalid date %q", candidate)
		}
	}
}

func TestCommentConflictAndPairs(t *testing.T) {
	entries := []string{"tItLe=old", "Foo=α", "Foo=β", "ARTIST=first", "ARTIST=second", "ALBUM ARTIST=alias", "TRACKNUMBER=2/11", "WAVEFORMATEXTENSIBLE_CHANNEL_MASK=0x3", "TITLESORT=first", "TITLESORT=second"}
	fields, _ := core.ReadComments(core.FLAC, entries)
	contracttest.AssertValues(t, fields[core.Artist], []core.Value{core.TextValue("first"), core.TextValue("second")})
	contracttest.AssertValues(t, fields[core.TrackTotal], []core.Value{core.IntValue(11)})
	changes := core.Changes{Fields: map[core.Field][]core.Value{core.Title: {core.TextValue("")}, core.AlbumArtist: {core.TextValue("canonical")}, core.TrackNumber: nil}}
	out, err := core.EditComments(core.FLAC, entries, changes)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"TITLE=", "Foo=α", "Foo=β", "ARTIST=first", "ARTIST=second", "ALBUM ARTIST=alias", "WAVEFORMATEXTENSIBLE_CHANNEL_MASK=0x3", "TITLESORT=first", "TITLESORT=second", "ALBUMARTIST=canonical", "TRACKTOTAL=11"}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("edits = %q", out)
	}
	out, err = core.EditComments(core.FLAC, entries, core.Changes{Fields: map[core.Field][]core.Value{core.TrackTotal: nil}, ClearUnknown: true})
	if err != nil {
		t.Fatal(err)
	}
	fields, _ = core.ReadComments(core.FLAC, out)
	if _, ok := fields[core.TrackTotal]; ok {
		t.Fatal("implicit total survived removal")
	}
	contracttest.AssertValues(t, fields[core.TrackNumber], []core.Value{core.IntValue(2)})
}

func FuzzCanonicalDecode(f *testing.F) {
	for _, d := range core.Fields() {
		s, err := core.Encode(d.ID, contracttest.Sample(d)[0])
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(d.ID), s)
	}
	f.Fuzz(func(t *testing.T, field, input string) {
		v, err := core.Decode(core.Field(field), input)
		if err != nil {
			return
		}
		s, err := core.Encode(core.Field(field), v)
		if err != nil {
			t.Fatal(err)
		}
		back, err := core.Decode(core.Field(field), s)
		if err != nil || !back.Equal(v) {
			t.Fatalf("canonical decode unstable: %v", err)
		}
	})
}

func TestMalformedStructuredTextAndPairedTotals(t *testing.T) {
	invalid := core.CreditValue(core.Credit{Name: string([]byte{255}), Role: "piano"})
	if err := core.Validate(core.PerformerCredits, invalid); !errors.Is(err, core.ErrInvalidValue) {
		t.Fatalf("invalid UTF-8 was normalized: %v", err)
	}
	entries := []string{"TRACKNUMBER=2/9", "TRACKTOTAL=10"}
	change := core.Changes{Fields: map[core.Field][]core.Value{core.TrackNumber: {core.IntValue(3)}}}
	if _, err := core.EditComments(core.FLAC, entries, change); !errors.Is(err, core.ErrUnsupported) {
		t.Fatalf("conflicting untouched total discarded: %v", err)
	}
	both := core.Changes{Fields: map[core.Field][]core.Value{core.TrackNumber: nil, core.TrackTotal: nil}}
	out, err := core.EditComments(core.FLAC, []string{"TRACKNUMBER=bad/9"}, both)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, []string{"TRACKNUMBER=bad"}) {
		t.Fatalf("opaque left component: %q", out)
	}
	fields, _ := core.ReadComments(core.FLAC, out)
	if _, ok := fields[core.TrackTotal]; ok {
		t.Fatal("cleared paired total reappeared")
	}
}
