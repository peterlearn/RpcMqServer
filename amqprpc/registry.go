package amqprpc

import "fmt"

type Registry struct {
	name    string
	methods map[string]Method
}

func NewRegistry(name string) *Registry {
	return &Registry{name: name, methods: make(map[string]Method)}
}

func (r *Registry) AddMethod(meth Method) {
	name := fmt.Sprintf("%s__%s", r.name, meth.GetName())
	r.methods[name] = meth
}

func (r *Registry) GetMethods() map[string]Method {
	return r.methods
}
