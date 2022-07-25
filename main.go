package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jobs-github/escript/eval"
)

func main() {
	argc := len(os.Args)
	e := eval.New(false)
	if argc == 1 {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		e.Repl(filepath.Dir(path), os.Stdin, os.Stdout)
	} else if argc == 2 {
		e.EvalScript(os.Args[1], nil)
	} else {
		if os.Args[1] == "--dump" {
			if s, err := e.DumpAst(os.Args[2]); nil != err {
				fmt.Println(err)
			} else {
				fmt.Println(s)
			}
		} else if os.Args[1] == "--load" {
			if argc == 3 {
				e.EvalJson(os.Args[2], nil)
			} else {
				e.EvalJson(os.Args[2], os.Args[3:])
			}
		} else {
			e.EvalScript(os.Args[1], os.Args[2:])
		}
	}
}
