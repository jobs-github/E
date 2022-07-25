package function

type StringSet map[string]bool

func NewStringSet(keys ...string) StringSet {
	s := StringSet{}
	for i := range keys {
		s.Add(keys[i])
	}
	return s
}

func (this *StringSet) AddKeys(keys []string) {
	if nil == keys {
		return
	}
	for _, key := range keys {
		this.Add(key)
	}
}

func (this *StringSet) Add(key string) {
	(*this)[key] = true
}

func (this *StringSet) Remove(key string) {
	delete(*this, key)
}

func (this *StringSet) Find(key string) bool {
	_, ok := (*this)[key]
	return ok
}

func (this *StringSet) Keys() []string {
	keys := []string{}
	for k, _ := range *this {
		keys = append(keys, k)
	}
	return keys
}
