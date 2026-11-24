package core

import (
	"fmt"
	"testing"
)

func TestCommentRatingPrecisionAndPreservation(t *testing.T) {
	original := []string{"RATING=0.125", "RATING:listener@example.org=0.8", "RATING=0.750"}
	for _, format := range []Format{FLAC, Vorbis, Opus} {
		t.Run(string(format), func(t *testing.T) {
			values, _ := ReadComments(format, original)
			if got := values[Rating]; len(got) != 1 || !got[0].Equal(IntValue(75)) {
				t.Fatalf("rating got %v", got)
			}
			out, err := EditComments(format, original, Changes{Fields: map[Field][]Value{Rating: {IntValue(1)}}})
			if err != nil {
				t.Fatal(err)
			}
			if len(out) != len(original) || out[0] != original[0] || out[1] != original[1] || out[2] != "RATING=0.01" {
				t.Fatalf("encoded comments %q", out)
			}
			cleared, err := EditComments(format, out, Changes{Fields: map[Field][]Value{Rating: nil}})
			if err != nil || len(cleared) != 2 || cleared[0] != original[0] || cleared[1] != original[1] {
				t.Fatalf("clear left %q, err=%v", cleared, err)
			}
			values, _ = ReadComments(format, []string{CommentFallbackName(Rating) + "=43"})
			if got := values[Rating]; len(got) != 1 || !got[0].Equal(IntValue(43)) {
				t.Fatalf("fallback rating got %v", got)
			}
		})
	}
}

func TestCommentRatingAllPercentValues(t *testing.T) {
	for n := int64(0); n <= 100; n++ {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			encoded, err := encodeComment(Rating, IntValue(n), true)
			if err != nil {
				t.Fatal(err)
			}
			got, err := decodeComment(Rating, "RATING", encoded)
			if err != nil || !got.Equal(IntValue(n)) {
				t.Fatalf("round trip %q: %v, %v", encoded, got, err)
			}
		})
	}
}
