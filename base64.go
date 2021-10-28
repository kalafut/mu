package mu

import "encoding/base64"

var std = base64.StdEncoding

func Base64Encode(src []byte) []byte {
	dst := make([]byte, std.EncodedLen(len(src)))
	std.Encode(dst, src)
	return dst
}

func Base64Decode(src []byte) ([]byte, error) {
	dst := make([]byte, std.DecodedLen(len(src)))
	n, err := std.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}
