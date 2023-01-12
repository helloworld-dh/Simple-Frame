package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context 和当前请求强相关的信息
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handler []HandleFunc
	index   int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (context *Context) Next() {
	n := len(context.handler)
	context.index++
	for ; context.index < n; context.index++ {
		context.handler[context.index](context)
	}
}

func (context *Context) Param(key string) string {
	val, _ := context.Params[key]
	return val
}

// PostForm 获取req中key对应的val
func (context *Context) PostForm(Key string) string {
	return context.Req.FormValue(Key)
}

// Query 获取url的查询参数
func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

// Status 设置status code
func (context *Context) Status(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

// SetHeader 设置response的header
func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Set(key, value)
}

func (context *Context) String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.Status(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (context *Context) JSON(code int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.Status(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

func (context *Context) Data(code int, data []byte) {
	context.Status(code)
	context.Writer.Write(data)
}

func (context *Context) HTML(code int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.Data(code, []byte(html))
}
