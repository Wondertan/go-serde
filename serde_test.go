package serde

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshal(t *testing.T) {
	in := &fakeMsg{data: []byte("test")}
	buf := make([]byte, 100)

	n, err := Marshal(in, buf)
	require.Nil(t, err)
	assert.Greater(t, n, in.Size())

	out := &fakeMsg{}
	nn, err := Unmarshal(out, buf)
	require.Nil(t, err)
	assert.Equal(t, n, nn)

	assert.Equal(t, in, out)
}

func TestWriteReadByteReader(t *testing.T) {
	in := &fakeMsg{data: []byte("test")}
	rw := &simpleRW{data: make([]byte, 100)}

	n, err := Write(rw, in)
	require.Nil(t, err)
	assert.NotEqual(t, n, in.Size())

	out := &fakeMsg{}
	nn, err := Read(NewByteReader(rw), out)
	require.Nil(t, err)
	assert.Equal(t, n, nn)
	assert.Equal(t, in, out)
}

type simpleRW struct {
	data []byte
	r, w int
}

func (rw *simpleRW) Write(b []byte) (n int, err error) {
	if len(rw.data) == rw.w {
		data := rw.data
		rw.data = make([]byte, len(data)*2)
		copy(rw.data, data)
	}
	n = copy(rw.data[rw.w:], b)
	rw.w += n
	return
}

func (rw *simpleRW) Read(b []byte) (n int, err error) {
	if len(rw.data) == rw.r {
		return 0, io.EOF
	}
	n = copy(b, rw.data[rw.r:])
	rw.r += n
	return
}

type fakeMsg struct {
	data []byte
}

func (f *fakeMsg) Size() int {
	return len(f.data)
}

func (f *fakeMsg) MarshalTo(buf []byte) (int, error) {
	return copy(buf, f.data), nil
}

func (f *fakeMsg) Unmarshal(data []byte) error {
	f.data = make([]byte, len(data))
	copy(f.data, data)
	return nil
}
