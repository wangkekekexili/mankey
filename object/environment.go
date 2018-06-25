package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: outer,
	}
}

func (e *Environment) Get(i string) (Object, bool) {
	o, ok := e.store[i]
	if !ok && e.outer != nil {
		o, ok = e.outer.Get(i)
	}
	return o, ok
}

func (e *Environment) Set(i string, o Object) {
	e.store[i] = o
}
