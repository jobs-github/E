package object

func NewState(v Object, quit bool) *State {
	obj := &State{
		Value: v,
		Quit:  quit,
	}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
	}
	return obj
}

// State : implement Object
type State struct {
	defaultObject
	Value Object
	Quit  bool
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
func (this *State) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
