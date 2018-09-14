package engine

import (
	"fmt"

	"github.com/tlmiller/bookie/pkg/domain"
	"github.com/tlmiller/bookie/pkg/executor"
)

type Engine struct {
	executors domain.DomainMap
}

func (e *Engine) AddExecutors(executors ...executor.Executor) error {
	for _, ex := range executors {
		if _, exists := e.executors.Get(ex.Authority()); exists {
			return fmt.Errorf("executor for domain %s already exists in engine",
				ex.Authority())
		}
		e.executors.Set(ex.Authority(), ex)
	}
	return nil
}

func NewEngine() *Engine {
	return &Engine{
		executors: domain.NewDomainMap(),
	}
}
