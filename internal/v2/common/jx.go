package common

import "github.com/go-faster/jx"

// StringToJx converts a string to a jx.Raw JSON representation.
func StringToJx(s string) jx.Raw {
	var e jx.Encoder
	e.Str(s)
	return e.Bytes()
}

// JxToString converts a jx.Raw JSON representation to a string.
func JxToString(r jx.Raw) (string, error) {
	return jx.DecodeStr(r.String()).Str()
}
