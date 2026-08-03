package progress

import "io"

// Progress defines the interface for task execution progress reporting.
type Progress interface {
	Start(total int64)
	Increment(n int64)
	Message(msg string)
	Finish()
}

// WriterProgress is a simple writer-based progress reporter
type WriterProgress struct {
	out   io.Writer
	total int64
	cur   int64
}

func NewWriterProgress(out io.Writer) *WriterProgress {
	return &WriterProgress{out: out}
}

func (p *WriterProgress) Start(total int64) {
	p.total = total
	p.cur = 0
}

func (p *WriterProgress) Increment(n int64) {
	p.cur += n
}

func (p *WriterProgress) Message(msg string) {
	// Simple console write representation
	_, _ = io.WriteString(p.out, msg+"\n")
}

func (p *WriterProgress) Finish() {
	p.cur = p.total
}
