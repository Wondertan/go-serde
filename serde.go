package serde

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Message interface {
	Size() int
	MarshalTo([]byte) (int, error)
	Unmarshal([]byte) error
}

func Marshal(msg Message, buf []byte) (int, error) {
	n := binary.PutUvarint(buf, uint64(msg.Size()))
	nn, err := msg.MarshalTo(buf[n:])
	n += nn
	if err != nil {
		return n, err
	}

	return n, nil
}

func Unmarshal(msg Message, data []byte) (int, error) {
	vint, n := binary.Uvarint(data)
	if n < 0 {
		return 0, fmt.Errorf("serde: varint overflow")
	}

	nn := n + int(vint)
	err := msg.Unmarshal(data[n:nn])
	if err != nil {
		return 0, err
	}

	return nn, nil
}

func Write(w io.Writer, msg Message) (int, error) {
	buf := Get(binary.MaxVarintLen64 + msg.Size())
	defer Put(buf)

	n, err := Marshal(msg, buf)
	if err != nil {
		return 0, err
	}

	return w.Write(buf[:n])
}

func Read(r io.Reader, msg Message) (n int, err error) {
	size, err := binary.ReadUvarint(&byteCounter{NewByteReader(r), &n})
	if err != nil {
		return
	}

	buf := Get(int(size))
	nn, err := readWith(r, msg, buf)
	n += nn
	Put(buf)
	return
}

func readWith(r io.Reader, msg Message, buf []byte) (int, error) {
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}

	return n, msg.Unmarshal(buf)
}

type ByteReader struct {
	io.Reader

	b [1]byte
}

func NewByteReader(r io.Reader) *ByteReader {
	return &ByteReader{Reader: r}
}

func (b *ByteReader) ReadByte() (byte, error) {
	_, err := io.ReadFull(b.Reader, b.b[:])
	if err != nil {
		return 0, err
	}

	return b.b[0], nil
}

type byteCounter struct {
	br io.ByteReader
	i  *int
}

func (bc *byteCounter) ReadByte() (byte, error) {
	b, err := bc.br.ReadByte()
	if err != nil {
		return 0, err
	}

	*bc.i++
	return b, nil
}
