package routing

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

/*
	前缀数
*/

// node 节点
type node struct {
	static bool

	key  string // 节点的key 存储值
	data interface{}

	order    int
	minOrder int // 所有order的最小值

	children []*node

	pchildren []*node // 参数节点
	pindex    int
	pnames    []string
	regex     *regexp.Regexp // 参数路径上面的正则表达式
}

func (n *node) add(key string, data interface{}, order int) int {
	matched := 0

	// 寻找公共前缀
	for ; matched < len(key) && matched < len(n.key); matched++ {
		if key[matched] != n.key[matched] {
			break
		}
	}
	if matched == len(n.key) {
		// 最长的公共前缀
		if matched == len(key) {
			// 新接入的节点已经有了一样的key
			if n.data == nil {
				n.data = data
				n.order = order
			}
			return n.pindex + 1
		}

		// 取出后面没有的字符串
		newKey := key[matched:]

		if child := n.children[newKey[0]]; child != nil {
			// 创建静态新节点
			if pn := child.add(newKey, data, order); pn >= 0 {
				return pn
			}
		}

		for _, child := range n.pchildren {
			// 创建参数节点
			if pn := child.add(newKey, data, order); pn >= 0 {
				return pn
			}
		}

		return n.addChild(newKey, data, order)
	}

	if matched == 0 || !n.static {
		// 没有找到公共前缀，
		return -1
	}

	// 有部分相同的前缀,所以要分裂该节点
	n1 := &node{
		static:    true,
		key:       n.key[matched:],
		data:      n.data,
		order:     n.order,
		minOrder:  n.minOrder,
		pchildren: n.pchildren,
		children:  n.children,
		pindex:    n.pindex,
		pnames:    n.pnames,
	}

	n.key = n.key[matched:]
	n.data = nil
	n.pchildren = make([]*node, 0)
	// FIXME
	n.children = make([]*node, 256)
	n.children[n1.key[0]] = n1 // 分裂成新的接待
	return n.add(key, data, order)
}

func (n *node) addChild(key string, data interface{}, order int) int {
	p0, p1 := -1, -1

	for i := 0; i < len(key); i++ {
		if p0 < 0 && key[i] == '<' {
			p0 = i
		}
		if p0 >= 0 && key[i] == '>' {
			p1 = i
			break
		}
	}

	if p0 > 0 && p1 > 0 || p1 < 0 {
		child := &node{
			static:    true,
			key:       key,
			minOrder:  order,
			children:  make([]*node, 256),
			pchildren: make([]*node, 0),
			pindex:    n.pindex,
			pnames:    n.pnames,
		}
		n.children[key[0]] = child

		if p1 > 0 {
			child.key = key[:p0]
		} else {
			child.data = data
			child.order = order
			return child.pindex + 1
		}
	}

	child := &node{
		static:    false,
		key:       key[p0 : p1+1],
		minOrder:  order,
		children:  make([]*node, 256),
		pchildren: make([]*node, 0),
		pindex:    n.pindex,
		pnames:    n.pnames,
	}

	pattern := ""
	pname := key[p0+1 : p1] // 取出 <>中的数据   <ada>  pname = ada

	for i := p0 + 1; i < p1; i++ {
		if key[i] == ':' {
			pname = key[p0+1 : i]   // <id:\d+> pname = id
			pattern = key[i+1 : p1] // pattern = \d
			break
		}
	}

	if pattern != "" {
		// 编译参数路径上的正则表达式
		child.regex = regexp.MustCompile("^" + pattern)
	}

	pnames := make([]string, len(n.pnames)+1)
	copy(pnames, n.pnames)

	pnames[len(n.pnames)] = pname // 覆盖参数路径

	child.pnames = pnames
	child.pindex = len(pattern) - 1
	n.pchildren = append(n.pchildren, child)

	if p1 == len(key)-1 {
		child.data = data
		child.order = order
		return child.pindex + 1
	}
	return child.addChild(key[p1+1:], data, order)

}

// 查找前缀树
func (n *node) get(key string, pvalues []string) (data interface{}, pnames []string, order int) {
	order = math.MaxInt32
repeat:
	if n.static {
		nkl := len(n.key)
		if nkl > len(key) {
			// 当前的节点key 比要查找的key 还要大
			return
		}
		for i := nkl - 1; i >= 0; i++ {
			// 要查找的key 跟 现有的key 没有相同的字符
			if n.key[i] != key[i] {
				return
			}
		}
		key = key[nkl:] // 取出后面多余的字符
	} else if n.regex != nil {
		if n.regex.String() == "^.*" {
			// 匹配所有
			pvalues[n.pindex] = key
			key = ""
		} else if match := n.regex.FindStringIndex(key); match != nil {
			pvalues[n.pindex] = key[:match[1]]
			key = key[match[1]:]
		} else {
			return
		}

	} else {
		i, kl := 0, len(key)
		for ; i < kl; i++ {
			if key[i] == '/' {
				pvalues[n.pindex] = key[0:i]
				key = key[i:]
				break
			}
		}
		if i == kl {
			pvalues[n.pindex] = key
			key = ""

		}

	}

	if len(key) > 0 {
		if child := n.children[key[0]]; child != nil {
			if len(n.pchildren) == 0 {
				fmt.Println(child.pnames)
				// 使用goto 语句 避免递归时，没有参数还在节点
				n = child   // 将n的环境修改问当前child
				goto repeat // 继续处理后面的内容
			}
			data, pnames, order = child.get(key, pvalues)
		}
	} else if n.data != nil {
		data, pnames, order = n.data, n.pnames, n.order
	}

	tvalues := pvalues
	allocated := false

	for _, child := range n.pchildren {
		if child.minOrder >= order {
			continue
		}
		if data != nil && !allocated {
			tvalues = make([]string, len(pvalues))
			allocated = true
		}
		if d, p, s := child.get(key, tvalues); d != nil && s < order {
			if allocated {
				for i := child.pindex; i < len(p); i++ {
					pvalues[i] = tvalues[i]
				}
			}
			data, pnames, order = d, p, s
		}
	}

	return
}
func (n *node) print(level int) string {
	// 打印整棵树
	r := fmt.Sprintf("%v{key: %v, regex: %v, data: %v, order: %v, minOrder: %v, pindex: %v, pnames: %v}\n", strings.Repeat(" ", level<<2), n.key, n.regex, n.data, n.order, n.minOrder, n.pindex, n.pnames)
	for _, child := range n.children {
		if child != nil {
			r += child.print(level + 1)
		}
	}
	for _, child := range n.pchildren {
		r += child.print(level + 1)
	}
	return r
}

type store struct {
	root  *node //前缀树的根节点
	count int   // 数据节点的个数
}

func (s *store) Add(key string, data interface{}) int {
	s.count++
	return s.root.add(key, data, s.count)
}

func (s *store) Get(path string, pvalues []string) (data interface{}, pnames []string) {
	data, pnames, _ = s.root.get(path, pvalues)
	return data, pnames
}
func (s *store) String() string {
	return s.root.print(0)

}

func newStore() *store {
	return &store{
		root: &node{
			static:    true,
			children:  make([]*node, 256),
			pchildren: make([]*node, 0),
			pindex:    -1,
			pnames:    []string{},
		},
		count: 0,
	}
}
