package sidestep

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBitStreamReading(t *testing.T) {

	cases := []struct {
		input        []byte
		bitCount     uint8
		firstOutput  byte
		secondOutput byte
	}{
		{
			input:        []byte{0b01010000},
			bitCount:     3,
			firstOutput:  0b10,
			secondOutput: 0b100,
		},
		{
			input:        []byte{0b01110111, 0b01110111},
			bitCount:     6,
			firstOutput:  0b011101,
			secondOutput: 0b110111,
		},
	}

	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			stream := newBitStream(test.input)
			firstOutput, err := stream.Read(test.bitCount)
			require.NoError(t, err)
			assert.Equal(t, test.firstOutput, firstOutput)
			secondOutput, err := stream.Read(test.bitCount)
			require.NoError(t, err)
			assert.Equal(t, test.secondOutput, secondOutput)
		})
	}

}

func TestBitStreamWriting(t *testing.T) {

	cases := []struct {
		bitCount    uint8
		firstWrite  byte
		secondWrite byte
		output      []byte
	}{
		{
			firstWrite:  0b11111111,
			secondWrite: 0b00000001,
			bitCount:    5,
			output:      []byte{0b11111000, 0b01000000},
		},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("%b_%d", test.output, test.bitCount), func(t *testing.T) {
			stream := newBitStream(nil)
			require.NoError(t, stream.Write(test.firstWrite, test.bitCount))
			require.NoError(t, stream.Write(test.secondWrite, test.bitCount))
			assert.Equal(t, test.output, stream.Data())
		})
	}

}
