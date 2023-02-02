package geecache

import (
	"fmt"
	"gin/geecache/consistenthash"
	"gin/geecache/lru"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geeCache/"
	defaultReplicas = 50
)

// HTTPPool 各节点之间通信的数据结构
type HTTPPool struct {
	self        string // 本节点地址（主机名/Ip/端口号）
	basePath    string //节点通信地址的前缀
	mu          sync.Mutex
	peers       *consistenthash.Map    // 一致性hash的map，根据key获得节点
	httpGetters map[string]*httpGetter // 每一个远程节点对应一个httpGetter
}

func NewHttpPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// 初始化一个一致性hash算法
	h.peers = consistenthash.NewMap(defaultReplicas, nil)
	// 将节点加入一致性hash的map中
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseUrl: peer + h.basePath}
	}
}

// PickPeer 根据key获得分布式节点，而后返回节点的客户端
func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		return h.httpGetters[peer], true
	}
	return nil, false
}

// 判断是否HTTPPool是否实现了PeerPicker接口，如果没实现就会报错
var _ PeerPicker = (*HTTPPool)(nil)

func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Serve %s] %s", h.self, fmt.Sprintf(format, v...))
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool serve a unexpected path :" + r.URL.Path)
	}
	h.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "URL is wrong", http.StatusNotFound)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	ca, err := group.Get(lru.Key(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(ca.ByteSlice())
}

type httpGetter struct {
	baseUrl string // 访问的远程节点 http://example.com/_geecache/
}

// Get 从group中根据key获取返回值
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}
