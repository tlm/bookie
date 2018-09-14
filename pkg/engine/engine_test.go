package engine

import (
	"errors"
	"testing"

	"github.com/tlmiller/bookie/pkg/domain"
	"github.com/tlmiller/bookie/pkg/k8/controller"
)

type dummyExecutor struct {
	RAuthority domain.Domain
}

func (d *dummyExecutor) ApplyAction(action *controller.Action) error {
	return errors.New("dummyExecutor does not implement ApplyAction not implemented")
}

func (d *dummyExecutor) Authority() domain.Domain {
	return d.RAuthority
}

func (d *dummyExecutor) Service() string {
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
