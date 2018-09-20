package engine

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tlmiller/bookie/pkg/domain"
	"github.com/tlmiller/bookie/pkg/k8/controller"
)

type dummyController struct {
}

func (c *dummyController) Run(stopCh <-chan struct{}) {
	<-stopCh
}

func (c *dummyController) SetEventHandler(_ controller.ActionEventHandler) {}

type dummyExecutor struct {
	RAuthority domain.Domain
}

func (e *dummyExecutor) ApplyAction(action *controller.Action) error {
	return errors.New("dummyExecutor does not implement ApplyAction not implemented")
}

func (e *dummyExecutor) Authority() domain.Domain {
	return e.RAuthority
}

func (e *dummyExecutor) Service() string {
	return "dummyExecutor"
}

func TestAddExcutorDuplication(t *testing.T) {
	e := NewEngine()
	ex1 := &dummyExecutor{domain.Domain("www.example.com")}
	err := e.AddExecutors(ex1)
	if err != nil {
		t.Fatalf("unexpected error for AddExecutors: %v", err)
	}
	ex2 := &dummyExecutor{domain.Domain("www.example.com")}
	err = e.AddExecutors(ex2)
	if err == nil {
		t.Fatalf("expected error for duplicate executor authority")
	}
}

func TestSafeRunStop(t *testing.T) {
	e := NewEngine()
	e.AddControllers(&dummyController{}, &dummyController{})
	e.Run()
	e.Stop()
}
