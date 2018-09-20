package ingress

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type ListWatch struct {
	ListFunc  func(metav1.ListOptions) (runtime.Object, error)
	WatchFunc func(metav1.ListOptions) (watch.Interface, error)
}

func NewListWatch(c *kubernetes.Clientset) *ListWatch {
	return &ListWatch{
		ListFunc: func(o metav1.ListOptions) (runtime.Object, error) {
			return c.Extensions().Ingresses(corev1.NamespaceAll).List(o)
		},
		WatchFunc: func(o metav1.ListOptions) (watch.Interface, error) {
			return c.Extensions().Ingresses(corev1.NamespaceAll).Watch(o)
		},
	}
}

func (l *ListWatch) List(o metav1.ListOptions) (runtime.Object, error) {
	return l.ListFunc(o)
}

func (l *ListWatch) Watch(o metav1.ListOptions) (watch.Interface, error) {
	return l.WatchFunc(o)
}
