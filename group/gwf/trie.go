package gwf

import "strings"

//使用trie树实现动态路由前缀匹配
type node struct {
	pattern  string  //匹配的路由
	part     string  //路由中的一部分
	children []*node //子节点
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func (node *node) matchChild(part string) *node {
	for _, child := range node.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//注册插入路由结点
func (n *node) insert(pattern string, parts []string, height int) {
	//若部分路由数与高度相等，说明已遍历到最后一个part
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		//若匹配不到对应子节点，则新建一个
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		//加入当前节点的子节点中
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//路由查询
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]

	//匹配子节点集
	children := n.matchChildren(part)

	for _, child := range children {
		//符合条件的子节点各自分别递归继续查询
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
