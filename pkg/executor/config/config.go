package config

import (
	"github.com/pkg/errors"

	"github.com/spf13/viper"

	"github.com/tlmiller/bookie/pkg/executor"
)

type executorConfig struct {
	Type string
	Conf map[string]string
}

const (
	EXECUTORS_KEY string = "executors"
)

func ExecutorsForConfig() ([]executor.Executor, error) {
	confs := make([]executorConfig, 0)
	err := viper.UnmarshalKey(EXECUTORS_KEY, &confs)
	if err != nil {
		return nil, errors.Wrap(err, "reading executors config")
	}

	executors := make([]executor.Executor, len(confs))
	for i, conf := range confs {
		ex, err := makeExecutor(&conf)
		if err != nil {
			return nil, errors.Wrap(err, "making executors from configuration")
		}
		executors[i] = ex
	}
	return executors, nil
}
