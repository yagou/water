package mysql

import (
	"io"
)

/**
 * 默认缓存字节大小
 * @type {Number}
 */
const defaultBufSize = 4096

type buffer struct {
	buf    []byte
	rd     io.Reader
	idx    int
	length int
}

func newBuffer(rd io.Reader) *buffer {
	var b [defaultBufSize]byte
	return &buffer{
		buf: b[:],
		rd:  rd,
	}
}

func (b *buffer) fill(need int) error {
	if b.length > 0 && b.idx > 0 {
		copy(b.buf[0:b.length], b.buf[b.idx:])
	}

	if need > len(b.buf) {
		newBuf := make([]byte, ((need/defaultBufSize)+1)*defaultBufSize)
		copy(newBuf, b.buf)
		b.buf = newBuf
	}

	// 重置read position
	b.idx = 0

	for {
		n, err := b.rd.Read(b.buf[b.length:])
		b.length += n
		if err == nil {
			if b.length < need {
				continue
			}
			return nil
		}
		if b.length >= need && err == io.EOF {
			return nil
		}
		return err
	}
}

func (b *buffer) readNext(need int) ([]byte, error) {
	if b.length < need {

		if err := b.fill(need); err != nil {
			return nil, err
		}
	}

	offset := b.idx
	b.idx += need
	b.length -= need
	return b.buf[offset:b.idx], nil
}
