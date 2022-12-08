package eval

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
	"github.com/jobs-github/escript/vm"
)

// interpreter : implement Eval
type interpreter struct{}

func (this interpreter) Repl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()
	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p, err := parser.New(l)
		if nil != err {
			io.WriteString(out, fmt.Sprintf("\t%v\n", err))
			continue
		}

		program, err := p.ParseProgram()
		if nil != err {
			io.WriteString(out, fmt.Sprintf("\t%v\n", err))
			continue
		}
		val, err := program.Eval(env)
		if nil != err {
			io.WriteString(out, err.Error())
			io.WriteString(out, "\n")
		} else {
			if !object.IsNull(val) {
				io.WriteString(out, val.String())
				io.WriteString(out, "\n")
			}
		}
	}
}

func (this interpreter) EvalJson(path string) {
	node, err := this.LoadJson(path)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	val, err := node.Eval(object.NewEnv())
	if nil != err {
		fmt.Println(err.Error())
	} else {
		if !object.IsNull(val) {
			fmt.Print(val.String())
		}
	}
}

func (this interpreter) EvalScript(path string) {
	b, err := loadCode(path)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	this.EvalCode(function.BytesToString(b))
}

func (this interpreter) EvalCode(code string) {
	node, err := this.LoadAst(code)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	val, err := node.Eval(object.NewEnv())
	if nil != err {
		fmt.Println(err.Error())
	} else {
		if !object.IsNull(val) {
			fmt.Print(val.String())
		}
	}
}

func (this interpreter) DumpAst(path string) (string, error) {
	b, err := loadCode(path)
	if nil != err {
		return "", function.NewError(err)
	}
	program, err := this.LoadAst(function.BytesToString(b))
	if nil != err {
		return "", function.NewError(err)
	}
	b, err = json.Marshal(program.Encode())
	if nil != err {
		return "", function.NewError(err)
	}
	return function.BytesToString(b), nil
}

func (this interpreter) LoadJson(path string) (ast.Node, error) {
	if !strings.HasSuffix(path, ast.SuffixJson) {
		err := fmt.Errorf(`file "%v" not endwith ".json"`, path)
		return nil, function.NewError(err)
	}
	b, err := function.LoadFile(path)
	if nil != err {
		return nil, function.NewError(err)
	}
	return ast.Decode(b)
}

func (this interpreter) LoadAst(code string) (ast.Node, error) {
	return loadAst(code)
}

// virtualMachine : implement Eval
type virtualMachine struct{}

func (this *virtualMachine) Repl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	consts := object.Objects{}
	globals := vm.NewGlobals()
	st := compiler.NewSymbolTable(nil)

	for {
		fmt.Fprintf(out, ">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p, err := parser.New(l)
		if nil != err {
			io.WriteString(out, fmt.Sprintf("\t%v\n", err))
			continue
		}

		program, err := p.ParseProgram()
		if nil != err {
			io.WriteString(out, fmt.Sprintf("\t%v\n", err))
			continue
		}

		c := compiler.Make(st, consts)
		if err := c.Compile(program); nil != err {
			fmt.Fprintf(out, fmt.Sprintf("compile error: %v", err))
			continue
		}
		machine := vm.Make(c.Bytecode(), c.Constants(), globals)
		if err := machine.Run(); nil != err {
			fmt.Fprintf(out, fmt.Sprintf("run vm error: %v\n", err))
			continue
		}
		e := machine.LastPopped()
		if !object.IsNull(e) {
			io.WriteString(out, e.String())
			io.WriteString(out, "\n")
		}
	}
}

func (this *virtualMachine) eval(program ast.Node) (object.Object, error) {
	consts := object.Objects{}
	globals := vm.NewGlobals()
	st := compiler.NewSymbolTable(nil)
	c := compiler.Make(st, consts)
	if err := c.Compile(program); nil != err {
		return object.Nil, function.NewError(err)
	}
	machine := vm.Make(c.Bytecode(), c.Constants(), globals)
	if err := machine.Run(); nil != err {
		return object.Nil, function.NewError(err)
	}
	return machine.LastPopped(), nil
}

func (this *virtualMachine) EvalJson(path string) {
	node, err := this.LoadJson(path)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	val, err := this.eval(node)
	if nil != err {
		fmt.Println(err.Error())
	} else {
		if !object.IsNull(val) {
			fmt.Print(val.String())
		}
	}
}

func (this *virtualMachine) EvalScript(path string) {
	b, err := loadCode(path)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	this.EvalCode(function.BytesToString(b))
}

func (this *virtualMachine) EvalCode(code string) {
	node, err := this.LoadAst(code)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	val, err := this.eval(node)
	if nil != err {
		fmt.Println(err.Error())
	} else {
		if !object.IsNull(val) {
			fmt.Print(val.String())
		}
	}
}

func (this *virtualMachine) DumpAst(path string) (string, error) {
	b, err := loadCode(path)
	if nil != err {
		return "", function.NewError(err)
	}
	program, err := this.LoadAst(function.BytesToString(b))
	if nil != err {
		return "", function.NewError(err)
	}
	b, err = json.Marshal(program.Encode())
	if nil != err {
		return "", function.NewError(err)
	}
	return function.BytesToString(b), nil
}

func (this *virtualMachine) LoadJson(path string) (ast.Node, error) {
	if !strings.HasSuffix(path, ast.SuffixJson) {
		err := fmt.Errorf(`file "%v" not endwith ".json"`, path)
		return nil, function.NewError(err)
	}
	b, err := function.LoadFile(path)
	if nil != err {
		return nil, function.NewError(err)
	}
	return ast.Decode(b)
}

func (this *virtualMachine) LoadAst(code string) (ast.Node, error) {
	return loadAst(code)
}
