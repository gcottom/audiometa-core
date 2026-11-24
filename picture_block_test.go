package core

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestPictureBlockRoundTrip(t *testing.T) {
	for _, image := range []Picture{
		{
			Role: 4, MIME: "image/png", Description: "裏表紙",
			Width: 512, Height: 256, Depth: 8, Colors: 32, Data: []byte{1, 2, 3, 4},
		},
		{Role: 3, MIME: "-->", Data: []byte("https://example.com/cover.png")},
		{Role: 3, MIME: "-->", Data: []byte("cover.png")},
	} {
		wire, err := MarshalPictureBlock(image)
		if err != nil {
			t.Fatal(err)
		}
		actual, err := ParsePictureBlock(wire)
		if err != nil {
			t.Fatal(err)
		}
		wire[len(wire)-1] ^= 0xff
		if !PictureValue(image).Equal(PictureValue(actual)) {
			t.Fatal("picture metadata or data changed after decode")
		}
	}
}

func TestPictureBlockRejectsInvalidSizes(t *testing.T) {
	wire, err := MarshalPictureBlock(Picture{Role: 3, MIME: "image/png", Data: []byte{1, 2, 3}})
	if err != nil {
		t.Fatal(err)
	}
	extra := append(append([]byte(nil), wire...), 0)
	oversized := append([]byte(nil), wire...)
	binary.BigEndian.PutUint32(oversized[4:8], ^uint32(0))
	for _, raw := range [][]byte{wire[:len(wire)-1], extra, oversized} {
		if _, err := ParsePictureBlock(raw); err == nil {
			t.Fatal("accepted invalid picture size")
		}
	}
}

func FuzzPictureBlock(f *testing.F) {
	wire, err := MarshalPictureBlock(Picture{Role: 3, MIME: "image/png", Data: []byte{1, 2, 3}})
	if err != nil {
		f.Fatal(err)
	}
	f.Add(wire)
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 1<<20 {
			t.Skip()
		}
		image, err := ParsePictureBlock(raw)
		if err != nil {
			return
		}
		got, err := MarshalPictureBlock(image)
		if err != nil || !bytes.Equal(got, raw) {
			t.Fatalf("picture changed after re-encode, err=%v", err)
		}
	})
}
