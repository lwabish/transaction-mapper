package consumer

import (
	"fmt"
	"github.com/lwabish/transaction-mapper/pkg/config"
)

var (
	Registry = &registry{
		constructors: make(map[string]constructor),
	}
)

type constructor func(*config.Config) Plugin

type registry struct {
	constructors map[string]constructor
}

func (r *registry) Get(name string, config *config.Config) (Plugin, error) {
	c, ok := r.constructors[name]
	if !ok {
		return nil, fmt.Errorf("unknown consumer: %s", name)
	}

	return c(config), nil
}

func (r *registry) Register(name string, c constructor) {
	r.constructors[name] = c
}
