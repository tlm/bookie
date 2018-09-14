package config

import (
	"net/url"

	"github.com/hashicorp/consul/api"

	"github.com/pkg/errors"

	"github.com/tlmiller/bookie/pkg/consul"
	"github.com/tlmiller/bookie/pkg/executor"
)

func Maker(conf map[string]string) (executor.Executor, error) {
	var exists bool
	var host string
	if host, exists = conf["host"]; !exists {
		return nil, errors.New("consul host not defined for executor")
	}

	hostU, err := url.Parse(host)
	if err != nil {
		return nil, errors.Wrap(err, "parsing consul host")
	}

	consConf := api.Config{}
	consConf.Scheme = hostU.Scheme
	consConf.Address = hostU.Host

	var dc string
	if dc, exists = conf["datacenter"]; !exists {
		return nil, errors.New("consul datacenter not defined for executor")
	}
	consConf.Datacenter = dc
	return consul.NewExecutor(&consConf)
}
