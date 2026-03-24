package common

import (
	"testing"

	"github.com/go-faster/jx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToJx_SimpleString(t *testing.T) {
	result := StringToJx("ReadOnly")
	assert.Equal(t, jx.Raw(`"ReadOnly"`), result)
}

func TestStringToJx_StringContainingJSON(t *testing.T) {
	input := `[{ "name": "debian-web", "host": "10.0.5.42", "user": "deploy", "port": "2222" }]`
	result := StringToJx(input)

	// The result must be valid JSON (a JSON string value)
	d := jx.DecodeStr(string(result))
	decoded, err := d.Str()
	require.NoError(t, err, "StringToJx output must be valid JSON, got: %s", string(result))
	assert.Equal(t, input, decoded)
}

func TestStringToJx_EmptyString(t *testing.T) {
	result := StringToJx("")
	assert.Equal(t, jx.Raw(`""`), result)
}

func TestRoundTrip_SimpleString(t *testing.T) {
	original := "localhost"
	raw := StringToJx(original)
	result, err := JxToString(raw)
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestRoundTrip_StringContainingJSON(t *testing.T) {
	original := `[{ "name": "centos-db", "host": "192.168.1.50" }]`
	raw := StringToJx(original)
	result, err := JxToString(raw)
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestRoundTrip_StringWithSpecialChars(t *testing.T) {
	original := "path\\to\\file\nnewline"
	raw := StringToJx(original)
	result, err := JxToString(raw)
	require.NoError(t, err)
	assert.Equal(t, original, result)
}
