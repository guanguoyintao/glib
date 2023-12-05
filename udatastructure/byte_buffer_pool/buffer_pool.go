package ubytebufferpool

import (
	"io"
)

type ByteBufferPool struct {
	ByteBuffer []byte
	EOF        bool
}

func (b *ByteBufferPool) Len() int {
	return len(b.ByteBuffer)
}

func (b *ByteBufferPool) ReadFrom(r io.Reader) (int64, error) {
	bb := b.ByteBuffer
	start := int64(len(bb))
	max := int64(cap(bb))
	n := start
	if max == 0 {
		max = 64
		bb = make([]byte, max)
	} else {
		bb = bb[:max]
	}
	for {
		if n == max {
			max *= 2
			bNew := make([]byte, max)
			copy(bNew, bb)
			bb = bNew
		}
		nn, err := r.Read(bb[n:])
		n += int64(nn)
		if err != nil {
			b.ByteBuffer = bb[:n]
			n -= start
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

func (b *ByteBufferPool) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.ByteBuffer)
	return int64(n), err
}

func (b *ByteBufferPool) Bytes() []byte {
	return b.ByteBuffer
}

func (b *ByteBufferPool) Write(bb []byte) (int, error) {
	b.ByteBuffer = append(b.ByteBuffer, bb...)
	return len(bb), nil
}

func (b *ByteBufferPool) WriteByte(c byte) error {
	b.ByteBuffer = append(b.ByteBuffer, c)
	return nil
}

func (b *ByteBufferPool) WriteString(s string) (int, error) {
	b.ByteBuffer = append(b.ByteBuffer, s...)
	return len(s), nil
}

func (b *ByteBufferPool) Set(bb []byte) {
	b.ByteBuffer = append(b.ByteBuffer[:0], bb...)
}

func (b *ByteBufferPool) SetString(s string) {
	b.ByteBuffer = append(b.ByteBuffer[:0], s...)
}

func (b *ByteBufferPool) String() string {
	return string(b.ByteBuffer)
}

func (b *ByteBufferPool) Reset() {
	b.ByteBuffer = b.ByteBuffer[:0]
}

func (b *ByteBufferPool) WriteEOF() {
	b.EOF = true
}

func (b *ByteBufferPool) ReadNBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, nil
	}
	if len(b.ByteBuffer) == 0 && b.EOF {
		return nil, io.EOF
	}
	if n > len(b.ByteBuffer) {
		n = len(b.ByteBuffer)
	}
	bytesToRead := b.ByteBuffer[:n]
	b.ByteBuffer = b.ByteBuffer[n:]
	return bytesToRead, nil
}
