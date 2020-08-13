package app

import (
	"bytes"
	"io"
)

type BytesUploader struct {
	name string
	b    []byte
}

func NewBytesUploader(name string, b []byte) *BytesUploader {
	return &BytesUploader{
		name: name,
		b:    b,
	}
}

func (b *BytesUploader) Name() string {
	return b.name
}

func (b *BytesUploader) Reader() (io.Reader, error) {
	return bytes.NewReader(b.b), nil
}

func (b *BytesUploader) Size() int64 {
	return int64(len(b.b))
}
