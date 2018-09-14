package executor

import (
	"github.com/tlmiller/bookie/pkg/domain"
	"github.com/tlmiller/bookie/pkg/k8/controller"
)

type Executor interface {
	ApplyAction(action *controller.Action) error
	Authority() domain.Domain
	Service() string
}
