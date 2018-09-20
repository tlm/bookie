package config

import (
	"time"

	"github.com/pkg/errors"

	"github.com/spf13/viper"

	"github.com/tlmiller/bookie/pkg/k8/controller"
	"github.com/tlmiller/bookie/pkg/k8/controller/ingress"

	"k8s.io/client-go/kubernetes"
)

func ControllersForConfig(client *kubernetes.Clientset) ([]controller.Controller, error) {
	if viper.GetBool("ingress.enabled") {
		return IngressControllerForConfig(client)
	}
	return []controller.Controller{}, nil
}

func IngressControllerForConfig(client *kubernetes.Clientset) ([]controller.Controller, error) {
	resync, err := time.ParseDuration(viper.GetString("ingress.resync"))
	if err != nil {
		return nil, errors.Wrap(err, "parsing ingress.resync duration")
	}

	return []controller.Controller{
		ingress.NewController(ingress.NewListWatch(client), resync),
	}, nil
}

func init() {
	viper.SetDefault("ingress.enabled", false)
	viper.SetDefault("ingress.resync", "60s")
}
