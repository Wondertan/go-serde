package serde

import (
	"encoding/binary"
	"io"

	pool "github.com/libp2p/go-buffer-pool"
	"github.com/libp2p/go-msgio"
)

type msg interface {
	Unmarshal([]byte) error
	MarshalTo([]byte) (int, error)
	Size() int
}

func WriteMessage(w io.Writer, msg msg) error {
	size := msg.Size()
	buf := pool.Get(size + binary.MaxVarintLen64)
	defer pool.Put(buf)

	n, err := MarshalMessage(msg, buf)
	if err != nil {
		return err
	}

	_, err = w.Write(buf[:n])
	return err
}

func MarshalMessage(msg msg, buf []byte) (int, error) {
	size := msg.Size()

	n := binary.PutUvarint(buf, uint64(size))
	n2, err := msg.MarshalTo(buf[n:])
	if err != nil {
		return 0, err
	}
	n += n2

	return n, nil
}

func ReadMessage(r io.Reader, msg msg) error {
	mr := msgio.NewVarintReader(r)
	b, err := mr.ReadMsg()
	if err != nil {
		return err
	}

	err = UnmarshalMessage(msg, b)
	mr.ReleaseMsg(b)
	if err != nil {
		return err
	}

	return nil
}

func UnmarshalMessage(msg msg, b []byte) error {
	err := msg.Unmarshal(b)
	if err != nil {
		return err
	}

	return nil
}
