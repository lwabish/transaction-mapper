package bank

import (
	"fmt"
	"github.com/samber/lo"
)

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
		return nil, fmt.Errorf("unknown bank: %s", name)
	}

	return c(), nil
}

func (r *registry) register(name string, c constructor) {
	r.constructors[name] = c
}

func (r *registry) List() []string {
	return lo.Keys(r.constructors)
}
