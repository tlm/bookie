package engine

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/tlmiller/bookie/pkg/domain"
	"github.com/tlmiller/bookie/pkg/executor"
	"github.com/tlmiller/bookie/pkg/k8/controller"
)

type Engine struct {
	controllers []controller.Controller
	executors   domain.DomainMap
	stopCh      chan struct{}
	waitGroup   *sync.WaitGroup
}

func (e *Engine) actionHandler(action *controller.Action) {
	le := log.WithFields(log.Fields{
		"domain":    action.FQDN,
		"domain_id": action.ID,
		"action":    action.Action,
	})
	le.Infof("finding executor for %s action", action.Action)

	obj, exists := e.executors.GetNearestRoot(action.FQDN)
	if !exists {
		le.Warn("failed to find executor for new action")
	}
	ex := obj.(executor.Executor)

	le = le.WithFields(log.Fields{
		"executor": ex.Authority(),
	})
	le.Info("found executor to handle action")

	err := ex.ApplyAction(action)
	if err != nil {
		le.Errorf("executor failed applying actiong: %v", err)
		return
	}
	le.Info("executor finished applying action")
}

func (e *Engine) AddControllers(controllers ...controller.Controller) {
	for _, c := range controllers {
		c.SetEventHandler(controller.ActionEventHandlerFunc(e.actionHandler))
	}
	e.controllers = append(e.controllers, controllers...)
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
		controllers: make([]controller.Controller, 0),
		executors:   domain.NewDomainMap(),
	}
}

func (e *Engine) Run() {
	e.stopCh = make(chan struct{})
	e.waitGroup = &sync.WaitGroup{}
	for _, c := range e.controllers {
		e.waitGroup.Add(1)
		go func() {
			c.Run(e.stopCh)
			e.waitGroup.Done()
		}()
	}
}

func (e *Engine) Stop() {
	close(e.stopCh)
	e.waitGroup.Wait()
}
