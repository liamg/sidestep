package sidestep

import (
	"fmt"
	"strings"
)

const base63Set = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-"
const bitsPerInputByte = 5

func encodeBase63(input []byte) string {

	stream := newBitStream(input)

	var encoded string

	for {
		index, err := stream.Read(bitsPerInputByte)
		if err != nil {
			break
		}
		encoded += base63Set[index : index+1]
	}

	padding := "0"
	if len(input)%2 > 0 {
		padding = "1"
	}
	return string(append([]byte(encoded), padding[0]))
}

func decodeBase63(input string) ([]byte, error) {

	padding := []byte(input)[len(input)-1]

	stream := newBitStream(nil)

	for i := range input {
		if i == len(input)-1 {
			break
		}
		index := strings.Index(base63Set, input[i:i+1])
		if index == -1 {
			return nil, fmt.Errorf("unexpected character in base63 data: %s", input[i:i+1])
		}
		if err := stream.Write(byte(index&0xff), bitsPerInputByte); err != nil {
			return nil, err
		}
	}

	data := stream.Data()

	expectedPadding := "0"
	if len(data)%2 > 0 {
		expectedPadding = "1"
	}

	if expectedPadding[0] != padding {
		data = data[:len(data)-1]
	}

	return data, nil
}
