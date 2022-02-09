package rlines

import (
	"bufio"
	"errors"
	"io"
)

type Reader struct {
	inp *bufio.Reader
	eof bool
}

func (r *Reader) readByte() (byte, error) {
	if r.eof {
		return 0, io.EOF
	}
	b, err := r.inp.ReadByte()
	if errors.Is(err, io.EOF) {
		r.eof = true
	}
	return b, err
}

type lineReader struct {
	r                  *Reader
	any, eof, maybeEOL bool
}

type ByteReader interface {
	io.Reader
	io.ByteReader
}

func NewReader(r io.Reader) *Reader {
	inp, ok := r.(*bufio.Reader)
	if !ok {
		inp = bufio.NewReader(r)
	}
	return &Reader{inp: inp}
}

func (r *Reader) Next() ByteReader {
	if r.eof {
		return nil
	}
	_, err := r.inp.Peek(1)
	if errors.Is(err, io.EOF) {
		r.eof = true
		return nil
	}
	return &lineReader{r: r}
}

func (r *lineReader) ReadByte() (byte, error) {
	if r.eof {
		return 0, io.EOF
	}
	if r.maybeEOL {
		return r.readByteMaybeEOL()
	}
	b, err := r.r.readByte()
	if errors.Is(err, io.EOF) {
		r.eof = true
		return 0, io.EOF
	}
	if err != nil {
		return 0, err
	}
	if b == '\r' {
		r.maybeEOL = true
		return r.readByteMaybeEOL()
	}
	if b == '\n' {
		r.eof = true
		return 0, io.EOF
	}
	return b, nil
}

func (r *lineReader) readByteMaybeEOL() (byte, error) {
	b, err := r.r.readByte()
	if errors.Is(err, io.EOF) {
		r.eof = true
		return '\r', nil
	}
	if err != nil {
		return 0, err
	}
	if b == '\n' {
		r.eof = true
		return 0, io.EOF
	}
	r.maybeEOL = (b == '\r')
	return b, nil
}

func (r *lineReader) Read(buf []byte) (int, error) {
	for i := 0; i < len(buf); i++ {
		b, err := r.ReadByte()
		if err != nil {
			return i, err
		}
		buf[i] = b
	}
	return len(buf), nil
}
