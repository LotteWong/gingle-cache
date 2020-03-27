package ginglecache

type ByteView struct {
	bytes []byte
}

func (v ByteView) Len() int {
	return len(v.bytes)
}

func (v ByteView) Bytes() []byte {
	return cloneByteView(v.bytes)
}

func (v ByteView) String() string {
	return string(v.bytes)
}

func cloneByteView(bytes []byte) []byte {
	copys := make([]byte, len(bytes))
	copy(copys, bytes)
	return copys
}
