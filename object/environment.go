package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func (e *Environment) Get(i string) (Object, bool) {
	o, ok := e.store[i]
	return o, ok
}

func (e *Environment) Set(i string, o Object) {
	e.store[i] = o
}
