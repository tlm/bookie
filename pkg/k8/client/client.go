package k8

import (
	"github.com/pkg/errors"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func ClientForConfig(k8Master string, k8Conf string) (*kubernetes.Clientset, error) {
	r, err := clientcmd.BuildConfigFromFlags("", k8Conf)
	if err != nil {
		return nil, errors.Wrap(err, "geting client rest config")
	}
	c, err := kubernetes.NewForConfig(r)
	if err != nil {
		return nil, errors.Wrap(err, "getting client for rest config")
	}
	return c, nil
}
