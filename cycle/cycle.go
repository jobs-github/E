package cycle

import (
	"fmt"

	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
)

func Detect(
	module string,
	loadAst func(module string) (ast.Node, error),
	nonrecursive bool,
) error {
	d := newDetector(loadAst)
	return d.detect(module, nonrecursive)
}

type Detector interface {
	detect(module string, nonrecursive bool) error
}

// detector : implement Detector
type detector struct {
	loadAst func(module string) (ast.Node, error)
	visited function.StringSet
	cycle   cycleStack
}

func newDetector(loadAst func(module string) (ast.Node, error)) Detector {
	return &detector{
		loadAst: loadAst,
		visited: function.StringSet{},
		cycle:   cycleStack{},
	}
}

func (this *detector) detect(module string, nonrecursive bool) error {
	if nonrecursive {
		return this.nrdetect(module)
	} else {
		return this.rdetect(0, module)
	}
}

func (this *detector) cyclic(depth int) func(next string) error {
	return func(next string) error {
		if !this.visited.Find(next) {
			if err := this.rdetect(depth+1, next); nil != err {
				return err
			}
		} else {
			if this.cycle.find(next) {
				return fmt.Errorf("cycle detected : %v -> `%v`", this.cycle.string(), next)
			}
		}
		return nil
	}
}

func (this *detector) nrcyclic(s *recursionStack, depth int) func(next string) (bool, error) {
	return func(next string) (bool, error) {
		if !this.visited.Find(next) {
			this.push(s, depth+1, next)
			return false, nil
		}
		if this.cycle.find(next) {
			return true, fmt.Errorf("cycle detected : %v -> `%v`", this.cycle.string(), next)
		}
		return true, nil
	}
}

func (this *detector) walk(
	module string,
	cb func(next string) error,
) error {
	modules, err := ast.Imports(module, this.loadAst)
	if nil != err {
		return function.NewError(err)
	}
	for _, m := range modules {
		if err := cb(m); nil != err {
			return err
		}
	}
	return nil
}

func (this *detector) nrwalk(
	module string,
	cb func(next string) (bool, error),
) (bool, error) {
	modules, err := ast.Imports(module, this.loadAst)
	if nil != err {
		return false, function.NewError(err)
	}
	for _, m := range modules {
		if visited, err := cb(m); !visited || nil != err {
			return false, err
		}
	}
	return true, nil
}

// recursive
func (this *detector) rdetect(depth int, module string) error {
	if !this.visited.Find(module) {
		this.visited.Add(module)
		this.cycle.push(module, depth)

		if err := this.walk(module, this.cyclic(depth)); nil != err {
			return err
		}
	}
	this.cycle.pop(module)
	return nil
}

// nonrecursive
func (this *detector) nrdetect(module string) error {
	s := &recursionStack{}
	this.push(s, 0, module)

	for {
		if s.empty() {
			break
		}
		node := s.top()

		done, err := this.nrwalk(node.key, this.nrcyclic(s, node.depth))

		if nil != err {
			return err
		}

		if done {
			this.pop(s)
		}
	}
	return nil
}

func (this *detector) push(s *recursionStack, depth int, key string) {
	s.push(&cycleNode{key, depth})
	this.cycle.push(key, depth)
	this.visited.Add(key)
}

func (this *detector) pop(s *recursionStack) {
	node := s.pop()
	if nil != node {
		this.cycle.pop(node.key)
	}
}
