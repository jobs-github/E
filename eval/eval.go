package eval

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
)

// evalImpl : implement Eval
type evalImpl struct {
}

func (this evalImpl) Repl(baseDir string, in io.Reader, out io.Writer) {
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

func (this evalImpl) EvalJson(path string) {
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

func (this evalImpl) EvalScript(path string) {
	b, err := loadCode(path)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	this.EvalCode(function.BytesToString(b))
}

func (this evalImpl) EvalCode(code string) {
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

func (this evalImpl) DumpAst(path string) (string, error) {
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

func (this evalImpl) LoadJson(path string) (ast.Node, error) {
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

func (this evalImpl) LoadAst(code string) (ast.Node, error) {
	return loadAst(code)
}
