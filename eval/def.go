package eval

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/lexer"
	"github.com/jobs-github/Q/object"
	"github.com/jobs-github/Q/parser"
)

type Eval interface {
	Repl(baseDir string, in io.Reader, out io.Writer)
	EvalJson(path string, args []string)
	EvalScript(path string, args []string)
	EvalCode(importer object.Importer, code string, args []string)

	EvalAst(node ast.Node, env object.Env) (object.Object, error)

	NewImporter(baseDir string, suffix string) object.Importer
	NewEnv(importer object.Importer, args []string) object.Env

	DumpAst(path string) (string, error)
	LoadJson(path string) (ast.Node, error)
	LoadAst(code string) (ast.Node, error)
}

func New(nonrecursive bool) Eval { return evalImpl{nonrecursive: nonrecursive} }

func loadCode(path string) ([]byte, error) {
	if !strings.HasSuffix(path, ast.SuffixQs) {
		err := fmt.Errorf(`file "%v" not endwith ".qs"`, path)
		return nil, function.NewError(err)
	}
	b, err := function.LoadFile(path)
	if nil != err {
		return nil, function.NewError(err)
	}
	return b, nil
}

func getModuleName(path string) string {
	return strings.Split(filepath.Base(path), ".")[0]
}

func loadAst(code string) (ast.Node, error) {
	l := lexer.New(code)
	p, err := parser.New(l)
	if nil != err {
		return nil, function.NewError(err)
	}
	return p.ParseProgram()
}

// importer : implement Importer
type importer struct {
	baseDir string
	suffix  string
}

func (this *importer) Load(module string) (object.Env, error) {
	node, err := ast.LoadAst(this.baseDir, this.suffix, loadAst)(module)
	if nil != err {
		return nil, function.NewError(err)
	}
	env := object.NewEnv(this)
	if _, err := node.Eval(env, false); nil != err {
		return nil, function.NewError(err)
	}
	return env, nil
}
