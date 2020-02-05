package sidestep

import (
	"fmt"
	"strings"
)

type clientPacket struct {
	transmission uint8 // transmission number (wraps)
	sequence     uint8 // sequence number (wraps)
	operation    byte
	baseSize     uint8
	data         []byte
}

const (
	OpOpen byte = 1 << iota
	OpSend
	OpReceive
)

func (p *clientPacket) ToDNS(domain string) string {
	rawBytes := append([]byte{
		p.transmission,
		p.sequence,
		p.operation,
		p.baseSize,
	}, p.data...)

	encoded := encodeBase63(rawBytes)
	name := domain
	for len(encoded) > 0 {
		if len(encoded) > 63 {
			name = fmt.Sprintf("%s.%s", encoded[len(encoded)-63:], name)
			encoded = encoded[:len(encoded)-63]
		} else {
			name = fmt.Sprintf("%s.%s", encoded, name)
			encoded = ""
		}
	}

	return name
}

func DecodePacket(name string) (*clientPacket, error) {

	name = strings.ReplaceAll(name, ".", "")

	data, err := decodeBase63(name)
	if err != nil {
		return nil, err
	}

	baseSize := data[3]

	encoded := name[:len(name)-int(baseSize-1)]

	safeData, err := decodeBase63(encoded)
	if err != nil {
		return nil, err
	}

	return &clientPacket{
		transmission: safeData[0],
		sequence:     safeData[1],
		operation:    safeData[2],
		baseSize:     safeData[3],
		data:         safeData[4:],
	}, nil
}
