package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

func (router *router) addRouter(method, pattern string, handler HandleFunc) {
	parts := parsePattern(pattern)
	if _, ok := router.roots[method]; !ok {
		router.roots[method] = &node{}
	}
	router.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	if _, ok := router.handlers[key]; ok {
		log.Fatal("the router have already added")
		return
	}
	router.handlers[key] = handler
}

// 对: * 路由转换
func (router *router) getRouter(method, path string) (*node, map[string]string) {
	searchRouter := parsePattern(path)
	params := make(map[string]string)
	if _, ok := router.roots[method]; !ok {
		return nil, nil
	}
	targetNode := router.roots[method].search(searchRouter, 0)
	if targetNode != nil {
		parts := parsePattern(targetNode.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchRouter[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchRouter[index:], "/")
				break
			}
		}
		return targetNode, params
	}
	return nil, nil
}

func (router *router) handler(ctx *Context) {
	n, params := router.getRouter(ctx.Method, ctx.Path)
	if n != nil {
		ctx.Params = params
		key := ctx.Method + "-" + ctx.Path
		router.handlers[key](ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
	}
}

// 解析pattern /login/dd/ii 变成数组
func parsePattern(pattern string) []string {
	partsBySp := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, part := range partsBySp {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}
