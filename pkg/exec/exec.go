package exec

import (
	"encoding/base64"
)

func DecodeURL(enc string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func EncodeURL(s string) string {
	src := []byte(s)
	return base64.URLEncoding.EncodeToString(src)
}
