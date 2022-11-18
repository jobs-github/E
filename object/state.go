package object

import (
	"fmt"
)

func NewState(quit bool, v Object) *State {
	obj := &State{
		Quit:  quit,
		Value: v,
	}
	obj.fns = objectBuiltins{
		FnValue: obj.builtinValue,
	}
	return obj
}

// State : implement Object
type State struct {
	defaultObject
	Quit  bool
	Value Object
}

func (this *State) String() string {
	return this.Value.String()
}

func (this *State) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *State) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *State) AsState() (*State, error) {
	return this, nil
}

func (this *State) getType() ObjectType {
	return objectTypeState
}

// builtin
func (this *State) builtinValue(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("value() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return this.Value, nil
}
