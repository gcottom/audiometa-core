package core

import (
	"encoding/binary"
	"errors"
	"mime"
	"net/url"
	"strings"
	"unicode/utf8"
)

const pictureBlockFieldsSize = 32

// ParsePictureBlock reads the payload from an RFC 9639 PICTURE block or its
// decoded METADATA_BLOCK_PICTURE comment representation. The returned picture
// owns its data. A linked image is left as a URI in Data and never retrieved.
func ParsePictureBlock(data []byte) (Picture, error) {
	if len(data) < pictureBlockFieldsSize {
		return Picture{}, errors.New("audiometa: short picture block")
	}
	picture := Picture{Role: binary.BigEndian.Uint32(data[:4])}
	mediaType, data, ok := splitPictureBytes(data[4:])
	if !ok {
		return Picture{}, errors.New("audiometa: invalid picture media type length")
	}
	description, data, ok := splitPictureBytes(data)
	if !ok || len(data) < 20 {
		return Picture{}, errors.New("audiometa: invalid picture description length")
	}
	picture.MIME = string(mediaType)
	picture.Description = string(description)
	picture.Width = binary.BigEndian.Uint32(data[0:4])
	picture.Height = binary.BigEndian.Uint32(data[4:8])
	picture.Depth = binary.BigEndian.Uint32(data[8:12])
	picture.Colors = binary.BigEndian.Uint32(data[12:16])
	imageData, rest, ok := splitPictureBytes(data[16:])
	if !ok || len(rest) != 0 {
		return Picture{}, errors.New("audiometa: invalid picture data length")
	}
	picture.Data = imageData
	if err := validatePictureBlock(picture); err != nil {
		return Picture{}, err
	}
	picture.Data = append([]byte(nil), imageData...)
	return picture, nil
}

// MarshalPictureBlock serializes a picture without decoding or re-encoding it.
// The caller applies the maximum payload size of its container.
func MarshalPictureBlock(picture Picture) ([]byte, error) {
	if err := validatePictureBlock(picture); err != nil {
		return nil, err
	}
	size := uint64(pictureBlockFieldsSize) + uint64(len(picture.MIME)) +
		uint64(len(picture.Description)) + uint64(len(picture.Data))
	if size > uint64(^uint32(0)) || size > uint64(^uint(0)>>1) {
		return nil, errors.New("audiometa: picture block is too large")
	}
	out := make([]byte, int(size))
	binary.BigEndian.PutUint32(out, picture.Role)
	data := out[4:]
	binary.BigEndian.PutUint32(data, uint32(len(picture.MIME)))
	n := copy(data[4:], picture.MIME)
	data = data[4+n:]
	binary.BigEndian.PutUint32(data, uint32(len(picture.Description)))
	n = copy(data[4:], picture.Description)
	data = data[4+n:]
	binary.BigEndian.PutUint32(data[0:4], picture.Width)
	binary.BigEndian.PutUint32(data[4:8], picture.Height)
	binary.BigEndian.PutUint32(data[8:12], picture.Depth)
	binary.BigEndian.PutUint32(data[12:16], picture.Colors)
	binary.BigEndian.PutUint32(data[16:20], uint32(len(picture.Data)))
	copy(data[20:], picture.Data)
	return out, nil
}

func splitPictureBytes(data []byte) (value, rest []byte, ok bool) {
	if len(data) < 4 {
		return nil, nil, false
	}
	size := binary.BigEndian.Uint32(data[:4])
	if uint64(size) > uint64(len(data)-4) {
		return nil, nil, false
	}
	end := 4 + int(size)
	return data[4:end], data[end:], true
}

func validatePictureBlock(picture Picture) error {
	if picture.Role > 20 || len(picture.Data) == 0 || !utf8.ValidString(picture.Description) {
		return errors.New("audiometa: invalid picture block")
	}
	for i := range picture.MIME {
		if picture.MIME[i] < 0x20 || picture.MIME[i] > 0x7e {
			return errors.New("audiometa: invalid picture media type")
		}
	}
	if picture.MIME == "-->" {
		text := string(picture.Data)
		if _, err := url.Parse(text); err != nil || strings.ContainsRune(text, ' ') || !utf8.ValidString(text) {
			return errors.New("audiometa: invalid linked picture URI")
		}
		return nil
	}
	mediaType, _, err := mime.ParseMediaType(picture.MIME)
	if err != nil || !strings.ContainsRune(mediaType, '/') {
		return errors.New("audiometa: invalid picture media type")
	}
	if picture.Role == 1 && (mediaType != "image/png" ||
		(picture.Width != 0 && picture.Width != 32) || (picture.Height != 0 && picture.Height != 32)) {
		return errors.New("audiometa: a file icon must be a 32x32 PNG")
	}
	return nil
}
