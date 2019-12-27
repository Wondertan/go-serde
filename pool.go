package serde

import (
	"github.com/libp2p/go-buffer-pool"
)

var bpool = new(pool.BufferPool)

func Get(len int) []byte {
	return bpool.Get(len)
}

func Put(buf []byte) {
	bpool.Put(buf)
}

func Extend(old []byte, len int) []byte {
	nw := Get(len)
	copy(nw, old)
	Put(old)
	return nw
}
