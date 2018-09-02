package controller

import (
	"github.com/tlmiller/bookie/pkg/domain"
)

type ActionMethod string

type Action struct {
	domain.ResourceRecord
	Action ActionMethod
}

const (
	METHOD_ADD    ActionMethod = "add"
	METHOD_CHECK  ActionMethod = "check"
	METHOD_DELETE ActionMethod = "delete"
	METHOD_UPDATE ActionMethod = "update"
)
