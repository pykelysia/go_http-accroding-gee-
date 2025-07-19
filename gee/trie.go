package gee

import "strings"

type node struct {
	// whole url path
	pattern string
	// part of url path
	part string
	// all of children in the trie tree
	children []*node
	// flag if it need to check strictly
	isWild bool
}

// 查找所有的子节点
// 如果有匹配的子节点则返回正确的子节点
// 如果没有匹配的子节点，则返回nil
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 查找所有子节点
// 返回所有可以匹配的子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 在合适的长度内深度插入所需要的节点
func (n *node) insert(pattern string, parts []string, height int) {
	// trie 树的深度与路径长度相等时将路径赋值给该子节点
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 未到达与路径长度相符合的树深
	// 则再次递归寻找正确的路径

	// 寻找正确的下一个子节点
	part := parts[height]
	child := n.matchChild(part)
	// 如果没有找到正确的下一个子节点
	// 则在本节点下新建一个正确的子节点
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 递归寻找终点节点
	child.insert(pattern, parts, height+1)
}

// 查询节点
func (n *node) search(parts []string, height int) *node {
	// trie 树的深度与路径长度等长时，说明应该在该节点找到对应的路径
	// 或者在该节点具有统配通配功能 ‘*’
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果该节点没有 pattern 说明该节点不具有实际功能
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 未到达与路径长度相符合的树深
	// 则再次递归寻找正确的路径

	// 寻找所有正确的下一个子节点
	part := parts[height+1]
	children := n.matchChildren(part)
	// 遍历递归所有正确的子节点，如果找到正确的子节点则返回该节点
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// 获取所有的实际存在的子节点
// travel 旅行，全走一遍来获取全部
func (n *node) travel(list *([]*node)) {
	// 该节点实际存在，则记录该节点
	if n.pattern != "" {
		*list = append(*list, n)
	}
	// 不论该节点是否实际存在，都遍历递归所有子节点
	for _, child := range n.children {
		child.travel(list)
	}
}
