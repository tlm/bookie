package config

import (
	"fmt"

	consul "github.com/tlmiller/bookie/pkg/consul/config"
	"github.com/tlmiller/bookie/pkg/executor"
)

type Maker func(map[string]string) (executor.Executor, error)

var (
	makers map[string]Maker
)

func init() {
	makers = make(map[string]Maker, 0)
	RegisterMaker("consul", consul.Maker)
}

func makeExecutor(c *executorConfig) (executor.Executor, error) {
	if m, exists := makers[c.Type]; exists {
		return m(c.Conf)
	}
	return nil, fmt.Errorf("no maker for executor type '%s'", c.Type)
}

func RegisterMaker(t string, m Maker) error {
	if _, exists := makers[t]; exists {
		return fmt.Errorf("maker for type %s already exists", t)
	}
	makers[t] = m
	return nil
}
