package consul

import (
	"fmt"

	"github.com/tlmiller/bookie/pkg/k8/controller"
)

func actionServiceID(a *controller.Action) string {
	return fmt.Sprintf("%s-%s", a.ID, string(a.Type))
}
