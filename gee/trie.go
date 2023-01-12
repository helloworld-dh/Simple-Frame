package gee

import "strings"

// 此处应该还有优先匹配具体的路由 以及每个pattern
type node struct {
	pattern  string  // 待匹配路由 /login/:dd/tt
	part     string  // 路由中的一部分 例如 tt
	children []*node // 子节点
	isWild   bool    // 是否精确匹配 如果是 * :dd 都是true
}

// matchChild 查找第一个路由树中第一个匹配到的节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// insert pattern==url part为url的各层路由 height为那一层
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 匹配所有成功的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0, 1)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 查找节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
