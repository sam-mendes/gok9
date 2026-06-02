package cassandra

import (
	"bytes"
	"testing"
)

// Real Cassandra v4 OPTIONS request frame
var optionsFrameBytes = []byte{
	// --- FIXED 9-BYTE HEADER ---
	0x04,       // Version: 4 (Request)
	0x00,       // Flags: 0
	0x00, 0x02, // Stream ID: 2
	0x05,                   // OpCode: 5 (OPTIONS)
	0x00, 0x00, 0x00, 0x05, // Body Length: 5 bytes! (No body follows)
}

func TestParsePacketHeader(t *testing.T) {
	reader := bytes.NewReader(optionsFrameBytes)

	wireHeader, err := BuildFrameHeader(reader)
	if err != nil {
		t.Fatalf("Failed to parse real packet header: %v", err)
	}

	if wireHeader.RequestMetadata.RequestType() != "request" {
		t.Errorf("Expected client request type, got %v", wireHeader.RequestMetadata.RequestType())
	}

	if wireHeader.RequestMetadata.ProtocolVersion() != 4 {
		t.Errorf("Expected protocol version 4, got %d", wireHeader.RequestMetadata.ProtocolVersion())
	}

	if wireHeader.OpCode != uint8(5) { // OpCode 1 is Startup
		t.Errorf("Expected OpCode 5, got %d", wireHeader.OpCode)
	}
	if wireHeader.BodyLength != 5 {
		t.Errorf("Expected BodyLength 22, got %d", wireHeader.BodyLength)
	}

	t.Logf("WireHeader: %v", wireHeader)
}
