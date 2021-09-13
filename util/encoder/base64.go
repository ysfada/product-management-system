package encoder

import "encoding/base64"

type Base64 struct{}

var _ Encoder = (*Base64)(nil)

func (endec *Base64) Encode(str string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(str))
}

func (endec *Base64) Decode(str string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(str)
}
