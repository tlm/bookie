package consul

import (
	"github.com/hashicorp/consul/api"

	"github.com/pkg/errors"

	"github.com/tlmiller/bookie/pkg/domain"
	"github.com/tlmiller/bookie/pkg/k8/controller"
)

type Executor struct {
	client     *api.Client
	datacenter string
}

const (
	DOMAIN domain.Domain = "consul"
)

func (e *Executor) ApplyAction(action *controller.Action) error {
	if action.Type != domain.A && action.Type != domain.AAAA {
		return errors.Errorf("unsupported RR type %s", action.Type)
	}
	var (
		ok  bool
		sub domain.Domain
	)
	if sub, ok = e.Authority().IsSubDomain(action.FQDN); !ok {
		return errors.Errorf("action fqdn '%s' is not a subdomain", action.FQDN)
	}

	err := e.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Address: action.Value,
		ID:      action.ID,
		Kind:    api.ServiceKindTypical,
		Name:    string(sub),
	})
	if err != nil {
		return errors.Wrapf(err, "applying consul service record for '%s'", sub)
	}
	return nil
}

func (e *Executor) Authority() domain.Domain {
	return DOMAIN.Prepend("service", e.datacenter)
}

func NewExecutor(conf *api.Config) (*Executor, error) {
	ex := &Executor{}
	var err error

	if conf.Datacenter == "" {
		errors.New("datacenter must be defined for consul")
	}
	ex.datacenter = conf.Datacenter

	ex.client, err = api.NewClient(conf)
	if err != nil {
		return nil, errors.Wrap(err, "creating executor consul client")
	}
	return ex, nil
}

func (e *Executor) Service() string {
	return "consul"
}
