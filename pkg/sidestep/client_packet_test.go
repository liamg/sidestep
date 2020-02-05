package sidestep

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientPacketEncodingPreservesData(t *testing.T) {

	domain := "example.org"

	cases := []clientPacket{
		{
			operation:    OpOpen,
			transmission: 7,
			sequence:     99,
			baseSize:     uint8(len(domain)),
			data:         []byte("127.0.0.1:9999"),
		},
		{
			operation:    OpReceive,
			transmission: 1,
			sequence:     0,
			baseSize:     uint8(len(domain)),
			data:         []byte{78, 32, 6, 3, 56},
		},
		{
			operation:    OpSend,
			transmission: 0,
			sequence:     1,
			baseSize:     uint8(len(domain)),
			data:         []byte{1, 2, 3},
		},
	}

	for _, packet := range cases {
		t.Run(fmt.Sprintf("op_%d", packet.operation), func(t *testing.T) {
			dns := packet.ToDNS(domain)
			output, err := DecodePacket(dns)
			require.NoError(t, err)
			assert.Equal(t, packet, *output)
		})
	}
}
