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
	if _, ok := e.Authority().IsSubDomain(action.FQDN); !ok {
		return errors.Errorf("action fqdn '%s' is not a subdomain", action.FQDN)
	}

	switch action.Action {
	case controller.METHOD_ADD, controller.METHOD_UPDATE:
		return e.upsertAction(action)
	case controller.METHOD_CHECK:
		return e.checkAction(action)
	case controller.METHOD_DELETE:
		return e.deleteAction(action)
	default:
		return errors.Errorf("unsupported action %s", action.Action)
	}
}

func (e *Executor) Authority() domain.Domain {
	return DOMAIN.Prepend("service", e.datacenter)
}

func (e *Executor) checkAction(action *controller.Action) error {
	var (
		ok  bool
		sub domain.Domain
	)
	if sub, ok = e.Authority().IsSubDomain(action.FQDN); !ok {
		return errors.Errorf("action fqdn '%s' is not a subdomain", action.FQDN)
	}

	services, _, err := e.client.Catalog().Service(string(sub), "", &api.QueryOptions{
		Datacenter: e.datacenter,
	})
	if err != nil {
		return errors.Wrapf(err, "checking service status for sub domain '%s'", sub)
	}

	found := false
	serviceID := actionServiceID(action)
	for _, s := range services {
		if s.ID != serviceID {
			continue
		}
		found = true
		if s.Address != action.Value ||
			s.Node != string(sub) ||
			s.ServiceAddress != action.Value {
			e.upsertAction(action)
		}
	}

	if !found {
		e.upsertAction(action)
	}
	return nil
}

func (e *Executor) Datacenter() string {
	return e.datacenter
}

func (e *Executor) deleteAction(action *controller.Action) error {
	var (
		ok  bool
		sub domain.Domain
	)
	if sub, ok = e.Authority().IsSubDomain(action.FQDN); !ok {
		return errors.Errorf("action fqdn '%s' is not a subdomain", action.FQDN)
	}

	_, err := e.client.Catalog().Deregister(&api.CatalogDeregistration{
		Datacenter: e.datacenter,
		ServiceID:  actionServiceID(action),
	}, &api.WriteOptions{})

	if err != nil {
		return errors.Wrapf(err, "deregistering previous sub domain '%s'", sub)
	}
	return nil
}

func NewExecutor(conf *api.Config) (*Executor, error) {
	if conf.Datacenter == "" {
		return nil, errors.New("datacenter must be defined for consul")
	}
	ex := &Executor{
		datacenter: conf.Datacenter,
	}
	ex.datacenter = conf.Datacenter

	var err error
	ex.client, err = api.NewClient(conf)
	if err != nil {
		return nil, errors.Wrap(err, "creating executor consul client")
	}
	return ex, nil
}

func (e *Executor) Service() string {
	return "consul"
}

func (e *Executor) upsertAction(action *controller.Action) error {
	var (
		ok  bool
		sub domain.Domain
	)
	if sub, ok = e.Authority().IsSubDomain(action.FQDN); !ok {
		return errors.Errorf("action fqdn '%s' is not a subdomain", action.FQDN)
	}

	_, err := e.client.Catalog().Register(&api.CatalogRegistration{
		Address:    action.Value,
		Datacenter: e.datacenter,
		Node:       string(sub),
		Service: &api.AgentService{
			Address: action.Value,
			ID:      actionServiceID(action),
			Service: string(sub),
		},
		SkipNodeUpdate: false,
	}, &api.WriteOptions{})

	if err != nil {
		return errors.Wrapf(err, "upserting consul service record for '%s'", sub)
	}
	return nil
}
