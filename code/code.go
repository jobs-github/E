package code

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jobs-github/escript/token"
)

var (
	errUnsupportedWidth = errors.New("unsupported width")
)

type Opcode byte

const (
	OpUndefined Opcode = iota
	OpClosure
	OpConst
	OpArray
	OpHash
	OpJumpWhenFalse
	OpJump
	OpGetGlobal
	OpSetGlobal
	OpGetBuiltin
	OpGetObjectFn
	OpGetLocal
	OpSetLocal
	OpIncLocal
	OpSetLocalIdx
	OpCall
	OpGetFree
	OpGetLambda
	OpReturn
	OpPop
	OpTrue
	OpFalse
	OpNull
	OpNot
	OpNeg
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpLt
	OpGt
	OpEq
	OpNeq
	OpLeq
	OpGeq
	OpAnd
	OpOr
	OpIndex
	OpLen
	OpNewArray
	OpPlaceholder
)

var (
	definitions = map[Opcode]*Definition{
		OpClosure:       {"OpClosure", []int{2, 1}},
		OpConst:         {"OpConst", []int{2}},
		OpArray:         {"OpArray", []int{2}},
		OpHash:          {"OpHash", []int{2}},
		OpJumpWhenFalse: {"OpJumpWhenFalse", []int{2}},
		OpJump:          {"OpJump", []int{2}},
		OpGetGlobal:     {"OpGetGlobal", []int{2}},
		OpSetGlobal:     {"OpSetGlobal", []int{2}},
		OpGetBuiltin:    {"OpGetBuiltin", []int{1}},
		OpGetObjectFn:   {"OpGetObjectFn", []int{1}},
		OpGetLocal:      {"OpGetLocal", []int{1}},
		OpSetLocal:      {"OpSetLocal", []int{1}},
		OpIncLocal:      {"OpIncLocal", []int{1}},
		OpSetLocalIdx:   {"OpSetLocalIdx", []int{1, 1}},
		OpCall:          {"OpCall", []int{1}},
		OpGetFree:       {"OpGetFree", []int{1}},
		OpGetLambda:     {"OpCurrentClosure", []int{1}},
		OpReturn:        {"OpReturn", []int{}},
		OpPop:           {"OpPop", []int{}},
		OpTrue:          {"OpTrue", []int{}},
		OpFalse:         {"OpFalse", []int{}},
		OpNull:          {"OpNull", []int{}},
		OpNot:           {"OpNot", []int{}},
		OpNeg:           {"OpNeg", []int{}},
		OpAdd:           {"OpAdd", []int{}},
		OpSub:           {"OpSub", []int{}},
		OpMul:           {"OpMul", []int{}},
		OpDiv:           {"OpDiv", []int{}},
		OpMod:           {"OpMod", []int{}},
		OpLt:            {"OpLt", []int{}},
		OpGt:            {"OpGt", []int{}},
		OpEq:            {"OpEq", []int{}},
		OpNeq:           {"OpNeq", []int{}},
		OpLeq:           {"OpLeq", []int{}},
		OpGeq:           {"OpGeq", []int{}},
		OpAnd:           {"OpAnd", []int{}},
		OpOr:            {"OpOr", []int{}},
		OpIndex:         {"OpIndex", []int{}},
		OpLen:           {"OpLen", []int{}},
		OpNewArray:      {"OpNewArray", []int{}},
		OpPlaceholder:   {"OpPlaceholder", []int{}},
	}
	prefixCodePairs = tokenCodePairs{
		{token.Not, OpNot},
		{token.Neg, OpNeg},
	}
	infixCodePairs = tokenCodePairs{
		{token.Add, OpAdd},
		{token.Sub, OpSub},
		{token.Mul, OpMul},
		{token.Div, OpDiv},
		{token.Mod, OpMod},
		{token.Lt, OpLt},
		{token.Gt, OpGt},
		{token.Eq, OpEq},
		{token.Neq, OpNeq},
		{token.Leq, OpLeq},
		{token.Geq, OpGeq},
		{token.And, OpAnd},
		{token.Or, OpOr},
	}
	prefixCodeMap = prefixCodePairs.newMap()
	infixCodeMap  = infixCodePairs.newMap()
)

func PrefixCode(t token.TokenType) (Opcode, error) {
	return prefixCodeMap.opCode(t)
}

func InfixToken(c Opcode) (*token.Token, error) {
	return infixCodeMap.token(c)
}

func InfixCode(t token.TokenType) (Opcode, error) {
	return infixCodeMap.opCode(t)
}

type tokenCodePair struct {
	t *token.Token
	c Opcode
}

type tokenCodePairs []*tokenCodePair

func (this *tokenCodePairs) newMap() *tokenCodeMap {
	opcodes := map[token.TokenType]Opcode{}
	tokens := map[Opcode]*token.Token{}
	for _, v := range *this {
		opcodes[v.t.Type] = v.c
		tokens[v.c] = v.t
	}
	return &tokenCodeMap{opcodes: opcodes, tokens: tokens}
}

type tokenCodeMap struct {
	opcodes map[token.TokenType]Opcode
	tokens  map[Opcode]*token.Token
}

func (this *tokenCodeMap) opCode(t token.TokenType) (Opcode, error) {
	v, ok := this.opcodes[t]
	if !ok {
		return OpUndefined, fmt.Errorf("no matched opcode, token: %v", token.ToString(t))
	}
	return v, nil
}

func (this *tokenCodeMap) token(c Opcode) (*token.Token, error) {
	v, ok := this.tokens[c]
	if !ok {
		return nil, fmt.Errorf("no matched token, opcode: %v", c)
	}
	return v, nil
}

type Instructions []byte

func (this *Instructions) String() string {

	var out bytes.Buffer
	sz := len(*this)
	for i := 0; i < sz; {
		d, err := Lookup(Opcode((*this)[i]))
		if nil != err {
			fmt.Fprintf(&out, "ERROR: %v\n", err)
			continue
		}
		r, err := DecodeOperands(d, (*this)[i+1:])
		if nil != err {
			fmt.Fprintf(&out, "ERROR: %v\n", err)
			continue
		}
		fmt.Fprintf(&out, "%04d %s\n", i, this.format(d, r.Value))
		i = i + 1 + r.Pos
	}
	return out.String()
}

func (this *Instructions) format(d *Definition, operands []int) string {
	sz := len(d.OperandWidths)
	if len(operands) != sz {
		return fmt.Sprintf("ERROR: operand len, want %v, got %v\n", sz, len(operands))
	}
	switch sz {
	case 0:
		return d.Name
	case 1:
		return fmt.Sprintf("%s %d", d.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", d.Name, operands[0], operands[1])
	}
	return fmt.Sprintf("ERROR: unsupport format for %s\n", d.Name)
}

type Operands struct {
	Value []int
	Pos   int
}

func (this *Operands) Add(i int, v int, w int) {
	this.Value[i] = v
	this.Pos += w
}

func DecodeOperands(d *Definition, ins Instructions) (*Operands, error) {
	r := &Operands{Value: make([]int, len(d.OperandWidths)), Pos: 0}
	for i, width := range d.OperandWidths {
		v, err := decodeOperand(width, ins[r.Pos:])
		if nil != err {
			return nil, err
		}
		r.Add(i, v, width)
	}
	return r, nil
}

type Definition struct {
	Name          string
	OperandWidths []int
}

func Lookup(op Opcode) (*Definition, error) {
	v, ok := definitions[op]
	if !ok {
		return nil, fmt.Errorf("undefined opcode: %v", op)
	}
	return v, nil
}

func Make(op Opcode, operands ...int) (Instructions, error) {
	v, err := Lookup(op)
	if nil != err {
		return nil, err
	}

	sz := 1
	for _, w := range v.OperandWidths {
		sz = sz + w
	}

	instruction := make(Instructions, sz)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := v.OperandWidths[i]
		err := encodeOperand(o, v.OperandWidths[i], instruction[offset:])
		if nil != err {
			return nil, err
		}
		offset += width
	}

	return instruction, nil
}

func encodeOperand(operand int, width int, b []byte) error {
	switch width {
	case 2:
		binary.BigEndian.PutUint16(b, uint16(operand))
		return nil
	case 1:
		b[0] = byte(operand)
		return nil
	default:
		return fmt.Errorf("unsupported width: %v", width)
	}
}

func DecodeUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func DecodeUint8(b []byte) uint8 {
	return uint8(b[0])
}

func decodeOperand(width int, b []byte) (int, error) {
	switch width {
	case 2:
		return int(DecodeUint16(b)), nil
	case 1:
		return int(DecodeUint8(b)), nil
	default:
		return -1, fmt.Errorf("unsupported width: %v", width)
	}
}
