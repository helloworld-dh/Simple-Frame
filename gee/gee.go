package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(ctx *Context)

// Engine 和RouterGroup中互相包含的是指针，与gin框架并未不一致
type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix     string
	middleware []HandleFunc
	parent     *RouterGroup
	engine     *Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(logger(), Recovery())
	return engine
}

// Group 路由组中的engine都是同一个
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middleware = append(group.middleware, middlewares...)
}

func (group *RouterGroup) addRouter(method string, comp string, handler HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.addRouter(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRouter("POST", pattern, handler)
}

func (engine *Engine) addRouter(method, pattern string, handler HandleFunc) {
	engine.router.addRouter(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandleFunc) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandleFunc) {
	engine.addRouter("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandleFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleware...)
		}
	}
	ctx := newContext(w, req)
	ctx.handler = middlewares
	engine.router.handler(ctx)
}

func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}
