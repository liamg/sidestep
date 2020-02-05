package sidestep

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestBase63EncodeDecode(t *testing.T) {
	cases := []string{
		"banana",
		"a",
		"",
		"0123456789",
		string([]byte{0, 0, 0, 0, 0}),
	}

	for _, plaintext := range cases {
		t.Run(plaintext, func(t *testing.T) {
			decoded, err := decodeBase63(encodeBase63([]byte(plaintext)))
			require.NoError(t, err)
			assert.Equal(t, []byte(plaintext), decoded)
		})
	}
}
