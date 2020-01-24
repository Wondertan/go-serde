package serde

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ReadUvarint(r io.Reader, buf []byte) (uint64, int, error) {
	if len(buf) < binary.MaxVarintLen64 {
		return 0, 0, fmt.Errorf("serde: small buffer")
	}

	nn, err := r.Read(buf)
	if err != nil {
		return 0, nn, err
	}

	vint, vn := binary.Uvarint(buf)
	if vn < 0 {
		return 0, nn, fmt.Errorf("serde: varint overflow")
	}

	return vint, vn, nil
}
