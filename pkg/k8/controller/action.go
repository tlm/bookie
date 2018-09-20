package controller

import (
	"github.com/tlmiller/bookie/pkg/domain"
)

type ActionMethod string

type Action struct {
	domain.ResourceRecord
	Action ActionMethod
}

type ActionEventHandler interface {
	Handle(*Action)
}

type ActionEventHandlerFunc func(*Action)

const (
	METHOD_ADD    ActionMethod = "add"
	METHOD_CHECK  ActionMethod = "check"
	METHOD_DELETE ActionMethod = "delete"
	METHOD_UPDATE ActionMethod = "update"
)

var (
	EmptyActionHandler = ActionEventHandlerFunc(func(_ *Action) {})
)

func (f ActionEventHandlerFunc) Handle(a *Action) {
	f(a)
}

func SendActionsToHandler(handle ActionEventHandler, actions ...*Action) {
	for _, action := range actions {
		handle.Handle(action)
	}
}
