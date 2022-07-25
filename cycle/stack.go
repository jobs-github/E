package cycle

import (
	"bytes"
	"sort"
)

type cycleNode struct {
	key   string
	depth int
}

type cycleSlice []*cycleNode

func (s cycleSlice) Len() int           { return len(s) }
func (s cycleSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s cycleSlice) Less(i, j int) bool { return s[i].depth < s[j].depth }

type recursionStack []*cycleNode

func (s *recursionStack) top() *cycleNode {
	return (*s)[len(*s)-1]
}

func (s *recursionStack) empty() bool {
	return len(*s) == 0
}

func (s *recursionStack) push(n *cycleNode) {
	*s = append(*s, n)
}

func (s *recursionStack) pop() *cycleNode {
	if !s.empty() {
		idx := len(*s) - 1
		top := (*s)[idx]
		*s = (*s)[:(idx)]
		return top
	}
	return nil
}

type cycleStack map[string]int

func (this *cycleStack) push(key string, depth int) {
	(*this)[key] = depth
}

func (this *cycleStack) find(key string) bool {
	_, ok := (*this)[key]
	return ok
}

func (this *cycleStack) pop(key string) {
	delete(*this, key)
}

func (this *cycleStack) string() string {
	nodes := cycleSlice{}
	for k, v := range *this {
		nodes = append(nodes, &cycleNode{k, v})
	}
	sort.Sort(nodes)
	var b bytes.Buffer
	sz := len(nodes)
	for i, v := range nodes {
		b.WriteString("`")
		b.WriteString(v.key)
		b.WriteString("`")
		if i < sz-1 {
			b.WriteString(" -> ")
		}
	}
	return b.String()
}
