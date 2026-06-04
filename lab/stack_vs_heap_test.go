package lab

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type FrameHeader struct {
	RequestMetadata uint8
	Flags           uint8
	StreamId        uint16
	OpCode          uint8
	BodyLength      uint32
}

// STRATEGY 1: Reflection-based (Forces Heap Escape)
func ParseWithBinaryRead(r *bytes.Reader) (FrameHeader, error) {
	var f FrameHeader
	// Passing &f into an interface forces it to escape to the heap
	err := binary.Read(r, binary.BigEndian, &f)
	return f, err
}

// STRATEGY 2: Manual Slicing & Return-by-Value (100% Stack Allocation)
func ParseWithSlicing(r *bytes.Reader) (FrameHeader, error) {
	var buf [9]byte
	if _, err := r.Read(buf[:]); err != nil {
		return FrameHeader{}, err
	}

	// Explicit structural decoding keeps memory local to the stack
	return FrameHeader{
		RequestMetadata: buf[0],
		Flags:           buf[1],
		StreamId:        binary.BigEndian.Uint16(buf[2:4]),
		OpCode:          buf[4],
		BodyLength:      binary.BigEndian.Uint32(buf[5:9]),
	}, nil
}

// --- BENCHMARK SUITE ---

func Benchmark_Strategy1_BinaryRead_Heap(b *testing.B) {
	// A valid 9-byte Cassandra header payload
	rawData := []byte{0x04, 0x00, 0x00, 0x02, 0x05, 0x00, 0x00, 0x00, 0x05}
	reader := bytes.NewReader(rawData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Reset(rawData) // Reset stream offset with zero allocations

		_, _ = ParseWithBinaryRead(reader)
	}
}

func Benchmark_Strategy2_Slicing_Stack(b *testing.B) {
	rawData := []byte{0x04, 0x00, 0x00, 0x02, 0x05, 0x00, 0x00, 0x00, 0x05}
	reader := bytes.NewReader(rawData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Reset(rawData) // Reset stream offset with zero allocations

		_, _ = ParseWithSlicing(reader)
	}
}
