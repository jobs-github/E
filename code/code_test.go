package code

import (
	"testing"
)

func newCode(op Opcode, operands ...int) Instructions {
	r, err := Make(op, operands...)
	if nil != err {
		return Instructions{}
	}
	return r
}

func TestMake(t *testing.T) {
	tests := []struct {
		name     string
		op       Opcode
		operands []int
		want     Instructions
	}{
		{"case_1", OpConst, []int{65534}, Instructions{byte(OpConst), 255, 254}},
		{"case_2", OpAdd, []int{}, Instructions{byte(OpAdd)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instrcution, err := Make(tt.op, tt.operands...)
			if nil != err {
				t.Fatal(err)
			}
			if len(instrcution) != len(tt.want) {
				t.Errorf("size mismatch %v, want %v", len(instrcution), len(tt.want))
				return
			}
			for i, b := range tt.want {
				if instrcution[i] != tt.want[i] {
					t.Errorf("wrong byte at pos %v, want %v, got %v", i, b, instrcution[i])
					return
				}
			}
		})
	}
}

func TestInstructionsString(t *testing.T) {
	ins := []Instructions{
		newCode(OpAdd),
		newCode(OpConst, 2),
		newCode(OpConst, 65535),
	}
	want := `0000 OpAdd
0001 OpConst 2
0004 OpConst 65535
`
	concatted := Instructions{}
	for _, i := range ins {
		concatted = append(concatted, i...)
	}
	got := concatted.String()
	if got != want {
		t.Errorf("wrong instruction\nwant: %v\ngot: %v", want, got)
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		name      string
		op        Opcode
		operands  []int
		readBytes int
	}{
		{"case_1", OpConst, []int{65535}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins, err := Make(tt.op, tt.operands...)
			if nil != err {
				t.Fatal(err)
			}
			d, err := Lookup(tt.op)
			if nil != err {
				t.Fatal(err)
			}
			r, err := DecodeOperands(d, ins[1:])
			if nil != err {
				t.Fatal(err)
			}
			if r.Pos != tt.readBytes {
				t.Fatalf("n wrong, want: %v, got: %v", tt.readBytes, r.Pos)
			}
			for i, want := range tt.operands {
				if r.Value[i] != want {
					t.Errorf("operand wrong, want: %v, got: %v", want, r.Value[i])
				}
			}
		})
	}
}
