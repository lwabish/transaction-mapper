package consumer

import "fmt"

var (
	Registry = &registry{
		constructors: make(map[string]constructor),
	}
)

type constructor func() Plugin

type registry struct {
	constructors map[string]constructor
}

func (r *registry) Get(name string) (Plugin, error) {
	c, ok := r.constructors[name]
	if !ok {
		return nil, fmt.Errorf("unknown consumer: %s", name)
	}

	return c(), nil
}

func (r *registry) Register(name string, c constructor) {
	r.constructors[name] = c
}
