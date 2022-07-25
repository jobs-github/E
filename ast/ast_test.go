package ast

import (
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Stmts: StatementSlice{
			&VarStmt{
				Name: &Identifier{
					Value: "testVar1",
				},
				Value: &Identifier{
					Value: "testVar2",
				},
			},
		},
	}
	if program.String() != "var testVar1 = testVar2;" {
		t.Errorf("program.String() wrong, got: %v", program.String())
	}
}
