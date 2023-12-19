package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	nameToValueMap := make(map[string]Object)
	return &Environment{
		store: nameToValueMap,
	}
}

func (self *Environment) Get(name string) (Object, bool) {
	value, ok := self.store[name]
	return value, ok
}

func (self *Environment) Set(name string, value Object) Object {
	self.store[name] = value
	return value
}
