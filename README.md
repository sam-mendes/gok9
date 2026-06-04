# gok9

## Comparing heap and stack allocation

`cd lab && go test -bench=. -benchmem`
|BENCHMARK NAME                                  |TOTAL OPS (b.N)          |SPEED PER OP          |HEAP BYTES/OP     |ALLOCS/OP  |
| :--- | :---: | :---: | :---: | :---: |
|Benchmark_Strategy1_BinaryRead_Heap-12          |22442587  (2 * 10^7)     |54.37 ns/op           |16 B/op           |1 allocs/op|
|Benchmark_Strategy2_Slicing_Stack-12            |358340168 (3 * 10^8)     |3.189 ns/op           | 0 B/op           |0 allocs/op|

### The Trade-off: Convinience with bit of magic Over Speed
* **Strategy 1 (Heap Allocation via Reflection):** Fewer lines of code, but significantly slower and impacts the garbage collector.
* **Strategy 2 (Zero Heap Allocation):** Gives us **~17x better performance** (over an order of magnitude increase in throughput)!

#### So the choice is clearly zero heap allocation, right?
**NO! Classic "It depends!"** 

However, after digging deep into the Cassandra v4 binary protocol specifications, manually slicing the frame buffer provides a distinct benefit: you can explicitly see exactly which raw incoming bytes map directly to each struct field. It makes the protocol implementation highly readable.

### The Implementations

```
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
````

Does `gok9` need that 17x performance gain? I don't know yet. But when building a tool to sniff a live TCP socket, minimizing operational overhead and 
skipping garbage collection cycles is obviously a win.
