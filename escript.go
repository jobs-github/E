package escript

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/object"

	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/parser"
	"github.com/jobs-github/escript/vm"
)

type RunnableType uint

const (
	RunnableTypeInterpreter RunnableType = iota
	RunnableTypeVM
)

type Runnable interface {
	Type() RunnableType
	Ast() ast.Node
	Run() (object.Object, error) // TODO: env
}

func NewInterpreter(code string) (Runnable, error) {
	node, err := LoadAst(code)
	if nil != err {
		return nil, function.NewError(err)
	}
	return &interpreter{node: node}, nil
}

func NewState(code string) (Runnable, error) {
	node, err := LoadAst(code)
	if nil != err {
		return nil, function.NewError(err)
	}
	consts := object.Objects{}
	globals := vm.NewGlobals()
	st := compiler.NewSymbolTable(nil)
	c := compiler.Make(st, consts)
	if err := c.Compile(node); nil != err {
		return nil, function.NewError(err)
	}
	s := vm.Make(c.Bytecode(), c.Constants(), globals)
	if nil != err {
		return nil, function.NewError(err)
	}
	return &virtualMachine{node: node, state: s}, nil
}

// interpreter : implement Runnable
type interpreter struct {
	node ast.Node
}

func (this *interpreter) Type() RunnableType {
	return RunnableTypeInterpreter
}

func (this *interpreter) Ast() ast.Node {
	return this.node
}

func (this *interpreter) Run() (object.Object, error) {
	return this.node.Eval(object.NewEnv())
}

// virtualMachine : implement Runnable
type virtualMachine struct {
	node  ast.Node
	state vm.VM
}

func (this *virtualMachine) Type() RunnableType {
	return RunnableTypeVM
}

func (this *virtualMachine) Ast() ast.Node {
	return this.node
}

func (this *virtualMachine) Run() (object.Object, error) {
	if err := this.state.Run(); nil != err {
		return nil, err
	}
	return this.state.LastPopped(), nil
}

func LoadAst(code string) (ast.Node, error) {
	p, err := parser.New(code)
	if nil != err {
		return nil, function.NewError(err)
	}
	return p.ParseProgram()
}
