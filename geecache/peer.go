package geecache

// PeerPicker 根据Key选择相应节点的PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 从对应的group中获取缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
