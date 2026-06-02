package cassandra

import (
	"encoding/binary"
	"io"
)

type RequestMetadata uint8

type FrameHeader struct {
	RequestMetadata RequestMetadata
	Flags           uint8
	StreamId        uint16
	OpCode          uint8
	BodyLength      uint32
}

func (r RequestMetadata) ProtocolVersion() uint8 {
	const bitmaskProtocolVersion = 0b0111_1111
	return uint8(r & bitmaskProtocolVersion)
}

func (r RequestMetadata) RequestType() string {
	const bitmaskReqOrRes = 0b1000_0000
	if (r & bitmaskReqOrRes) == 0 {
		return "request"
	} else {
		return "response"
	}
}

func BuildFrameHeader(reader io.Reader) (FrameHeader, error) {
	// Represent the structured binary data
	var wireHeader FrameHeader

	if err := binary.Read(reader, binary.BigEndian, &wireHeader); err != nil {
		return wireHeader, err
	}

	return wireHeader, nil
}
