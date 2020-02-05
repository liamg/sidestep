package sidestep

import (
	"fmt"
	"io"
)

type bitStream struct {
	pointer uint64
	data    []byte
	length  uint64
}

func newBitStream(data []byte) *bitStream {
	if data == nil {
		data = []byte{}
	}
	return &bitStream{data: data}
}

func (s *bitStream) Data() []byte {
	return s.data
}

func (s *bitStream) Write(data byte, bitCount uint8) error {
	if bitCount > 8 {
		return fmt.Errorf("cannot write more than 8 bits of data per byte")
	}

	s.length += uint64(bitCount)

	value := data << (8 - bitCount)

	bytePointer := s.pointer / 8

	for bytePointer >= uint64(len(s.data)) {
		s.data = append(s.data, 0)
	}

	rawByte := s.data[bytePointer]
	bitOffset := uint8(s.pointer % 8)

	s.data[bytePointer] = rawByte ^ (value >> bitOffset)

	s.pointer += uint64(bitCount)

	if bitCount+bitOffset <= 8 {
		return nil
	}

	bytePointer++

	if bytePointer >= uint64(len(s.data)) {
		s.data = append(s.data, 0)
	}

	trailingBits := (bitOffset + bitCount) - 8

	nextByte := s.data[bytePointer]

	s.data[bytePointer] = nextByte ^ (data << (8 - trailingBits))

	return nil
}

func (s *bitStream) Read(bitCount uint8) (byte, error) {

	if bitCount > 8 {
		return 0, fmt.Errorf("cannot read more than 1 byte of data per read")
	}

	bytePointer := s.pointer / 8

	if bytePointer >= uint64(len(s.data)) {
		return 0, io.EOF
	}

	rawByte := s.data[bytePointer]

	bitOffset := uint8(s.pointer % 8)

	rawByte = rawByte << bitOffset

	if bitOffset+bitCount <= 8 {
		s.pointer += uint64(bitCount)
		return rawByte >> (8 - bitCount), nil
	}

	if bytePointer+1 >= uint64(len(s.data)) {
		s.pointer += uint64(bitCount)
		return rawByte >> (8 - bitCount), nil
	}

	bitOffset = 8 - bitOffset

	nextByte := s.data[bytePointer+1] >> bitOffset

	s.pointer += uint64(bitCount)

	return (rawByte ^ nextByte) >> (8 - bitCount), nil
}
