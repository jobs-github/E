package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/eval"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
	"github.com/jobs-github/escript/vm"
)

func Start(in io.Reader, out io.Writer) {
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

func intepreterMain() {
	argc := len(os.Args)
	e := eval.New()
	if argc == 1 {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		e.Repl(filepath.Dir(path), os.Stdin, os.Stdout)
	} else if argc == 2 {
		e.EvalScript(os.Args[1])
	} else {
		if os.Args[1] == "--dump" {
			if s, err := e.DumpAst(os.Args[2]); nil != err {
				fmt.Println(err)
			} else {
				fmt.Println(s)
			}
		} else if os.Args[1] == "--load" {
			if argc == 3 {
				e.EvalJson(os.Args[2])
			}
		}
	}
}

func main() {
	// Start(os.Stdin, os.Stdout)
	intepreterMain()
}
