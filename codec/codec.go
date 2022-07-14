package codec

import (
	b64 "encoding/base64"
)

func Encode(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

func Decode(data string) (string, error) {
	decoded, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
