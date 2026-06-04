package cassandra

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real Cassandra v4 OPTIONS request frame
var optionsFrameBytes = []byte{
	// --- FIXED 9-BYTE HEADER ---
	0x04,     				// Version: 4 (Request)
	0x00,      				// Flags: 0
	0x00, 0x02, 			// Stream ID: 2
	0x05,                   // OpCode: 5 (OPTIONS)
	0x00, 0x00, 0x00, 0x05, // Body Length: 5 bytes! (No body follows)
}

func TestParsePacketHeader(t *testing.T) {
	reader := bytes.NewReader(optionsFrameBytes)

	frameHeader, err := BuildFrameHeader(reader)

	if err != nil {
		t.Fatalf("Failed to parse real packet header: %v", err)
	}

	assert := assert.New(t)

	assert.Equal("request", frameHeader.RequestMetadata.RequestType())
	assert.Equal(uint8(4), frameHeader.RequestMetadata.ProtocolVersion())
	assert.Equal(uint8(5), frameHeader.OpCode)
	assert.Equal(uint32(5), frameHeader.BodyLength)
}
