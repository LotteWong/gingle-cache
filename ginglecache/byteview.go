package ginglecache

// ByteView is read-only bytes
type ByteView struct {
	bytes []byte
}

// Len make ByteView implement Value interface
func (v ByteView) Len() int {
	return len(v.bytes)
}

// Bytes converts ByteView to []byte
func (v ByteView) Bytes() []byte {
	return cloneByteView(v.bytes)
}

// String() converts ByteView to string
func (v ByteView) String() string {
	return string(v.bytes)
}

// cloneByteViews helps clone data because Byteview is read-only
func cloneByteView(bytes []byte) []byte {
	copys := make([]byte, len(bytes))
	copy(copys, bytes)
	return copys
}
