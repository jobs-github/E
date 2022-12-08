package main

import (
	"fmt"
	"os"

	"github.com/jobs-github/escript/eval"
)

func intepreterMain() {
	argc := len(os.Args)
	e := eval.NewInterpreter()
	if argc == 1 {
		e.Repl(os.Stdin, os.Stdout)
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
	intepreterMain()
}
