package atticuscache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_atticuscache"

type HTTPPool struct {
	self     string
	basePath string
}

// NewHTTPPool init
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServerHTTP(w http.ResponseWriter, r http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPHTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	// 将数据写进去
	w.Write(view.ByteSlice())
}
