package object

import (
	"fmt"
	"os"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewFile(
	url string,
	mode string,
	f *os.File,
	sz int64,
) Object {
	obj := &FileObject{url: url, mode: mode, f: f, sz: sz}
	obj.fns = objectBuiltins{
		FnNot:   obj.builtinNot,
		"close": obj.builtinClose,
		"seek":  obj.builtinSeek,
		"size":  obj.builtinSize,
		"read":  obj.builtinRead,
		"write": obj.builtinWrite,
	}
	return obj
}

// FileObject : implement Object
type FileObject struct {
	url  string
	mode string
	f    *os.File
	sz   int64
	fns  objectBuiltins
}

func (this *FileObject) String() string {
	return fmt.Sprintf("<open file '%v', mode '%v' at %v>", this.url, this.mode, this.f)
}

func (this *FileObject) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *FileObject) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *FileObject) Calc(op *token.Token, right Object) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *FileObject) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *FileObject) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *FileObject) True() bool {
	return false
}

func (this *FileObject) Return() (bool, Object) {
	return false, nil
}

func (this *FileObject) Break() (bool, int) {
	return false, 0
}

func (this *FileObject) getType() ObjectType {
	return objectTypeFile
}

func (this *FileObject) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *FileObject) equal(other Object) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalInteger(other *Integer) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalString(other *String) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalBoolean(other *Boolean) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalNull(other *Null) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalArray(other *Array) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalHash(other *Hash) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalBuiltin(other *Builtin) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalFunction(other *Function) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalObjectFunc(other *ObjectFunc) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalArrayIter(other *ArrayIterator) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) equalHashIter(other *HashIterator) error {
	return unsupported(function.GetFunc(), this)
}

func (this *FileObject) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcString(op *token.Token, left *String) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcNull(op *token.Token, left *Null) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcArray(op *token.Token, left *Array) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcHash(op *token.Token, left *Hash) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}
func (this *FileObject) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcFunction(op *token.Token, left *Function) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *FileObject) calcHashIter(op *token.Token, left *HashIterator) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

// functions
type seekArgs struct {
	offset int64
	// 0 means relative to the origin of the file
	// 1 means relative to the current offset
	// 2 means relative to the end.
	whence int
}

func (this *FileObject) newSeekArgs(args Objects) (*seekArgs, error) {
	argc := len(args)
	if argc < 1 {
		return nil, fmt.Errorf("seek() takes at least 1 argument (%v given)", argc)
	}
	if argc > 2 {
		return nil, fmt.Errorf("seek() takes at most 2 arguments (%v given)", argc)
	}
	offset, err := args[0].asInteger()
	if nil != err {
		return nil, function.NewError(err)
	}
	whence := int(0)
	if argc == 2 {
		v, err := args[1].asInteger()
		if nil != err {
			return nil, function.NewError(err)
		}
		whence = int(v)
	}
	return &seekArgs{offset: offset, whence: whence}, nil
}

// builtin
func (this *FileObject) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}

func (this *FileObject) builtinClose(args Objects) (Object, error) {
	if nil == this.f {
		return Nil, nil
	}
	this.f.Close()
	this.f = nil
	return Nil, nil
}

func (this *FileObject) builtinSeek(args Objects) (Object, error) {
	r, err := this.newSeekArgs(args)
	if nil != err {
		return False, function.NewError(err)
	}
	if nil == this.f {
		return False, nil
	}
	if _, err := this.f.Seek(r.offset, r.whence); nil != err {
		return False, function.NewError(err)
	}
	return True, nil
}

func (this *FileObject) builtinSize(args Objects) (Object, error) {
	if nil == this.f {
		return NewInteger(-1), nil
	}
	return NewInteger(this.sz), nil
}

func (this *FileObject) builtinRead(args Objects) (Object, error) {
	if nil == this.f {
		return NewString(""), nil
	}
	if this.sz < 1 {
		return NewString(""), nil
	}
	b := make([]byte, this.sz)
	if _, err := this.f.Read(b); nil != err {
		return nil, err
	}
	return NewString(function.BytesToString(b)), nil
}

func (this *FileObject) builtinWrite(args Objects) (Object, error) {
	if nil == this.f {
		return False, nil
	}
	argc := len(args)
	if argc != 1 {
		return False, fmt.Errorf("write() takes one argument (%v given)", argc)
	}
	if !IsString(args[0]) {
		return False, fmt.Errorf("write() the first argument should be string (%v given)", Typeof(args[0]))
	}
	if _, err := this.f.WriteString(args[0].String()); nil != err {
		return False, function.NewError(err)
	}
	return True, nil
}
