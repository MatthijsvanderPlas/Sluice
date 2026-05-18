package buffer

type Buffer interface {
	Add(line string)
	Snapshot() []string
}

type RingBuffer struct {
	data     []string
	writePos int
	count    int
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		data: make([]string, capacity),
	}
}

func (rb *RingBuffer) Add(line string) {
	rb.data[rb.writePos] = line
	rb.writePos = (rb.writePos + 1) % len(rb.data)

	if rb.count < len(rb.data) {
		rb.count++
	}
}

func (rb *RingBuffer) Snapshot() []string {
	out := make([]string, rb.count)
	if rb.count < len(rb.data) {
		copy(out, rb.data[:rb.count])
		return out
	}
	for i := 0; i < rb.count; i++ {
		out[i] = rb.data[(rb.writePos+i)%len(rb.data)]
	}
	return out
}
