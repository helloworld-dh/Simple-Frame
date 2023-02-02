package geecache

// ByteView 存储真实的缓存值且只读
type ByteView struct {
	b []byte
}

func (byteView *ByteView) ByteSlice() []byte {
	return cloneBytes(byteView.b)
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return b
	}
	res := make([]byte, len(b))
	copy(res, b)
	return res
}

func (byteView *ByteView) String() string {
	return string(byteView.b)
}
