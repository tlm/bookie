package ingress

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tlmiller/bookie/pkg/domain"
	"github.com/tlmiller/bookie/pkg/k8/controller"

	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	controller   cache.Controller
	eventHandler controller.ActionEventHandler
	store        cache.Store
}

func actionsFromIngress(
	ingObj *extv1beta1.Ingress, meth controller.ActionMethod) ([]*controller.Action, error) {

	endpoints := make([]net.IP, 0)
	for _, ilb := range ingObj.Status.LoadBalancer.Ingress {
		if ilb.IP == "" {
			continue
		}
		r := net.ParseIP(ilb.IP)
		if r == nil {
			return nil, fmt.Errorf("ingress status loadbalancer address %s is not valid", ilb.IP)
		}
		endpoints = append(endpoints, r)
	}

	if len(endpoints) == 0 {
		// no endpoints nothing todo
		return []*controller.Action{}, nil
	}

	actions := make([]*controller.Action, 0)
	for _, rule := range ingObj.Spec.Rules {
		for _, endp := range endpoints {
			rrType := domain.A
			if endp.To4() == nil {
				rrType = domain.AAAA
			}
			actions = append(actions, &controller.Action{
				domain.ResourceRecord{
					FQDN:  domain.Domain(rule.Host),
					ID:    string(ingObj.UID),
					Type:  rrType,
					Value: endp.String(),
				},
				meth,
			})
		}
	}
	return actions, nil
}

func (c *Controller) ingressAdd(obj interface{}) {
	ingObj, ok := obj.(*extv1beta1.Ingress)
	if !ok {
		log.Error("unknown ingress object for informer add event")
	}

	le := log.WithFields(log.Fields{
		"k8_event":        "add",
		"k8_uid":          ingObj.UID,
		"k8_ingress_name": ingObj.Name,
		"k8_namespace":    ingObj.Namespace,
	})
	le.Info("new ingress object")

	actions, err := actionsFromIngress(ingObj, controller.METHOD_ADD)
	if err != nil {
		le.Errorf("failed creating actions from ingress object: %v", err)
		return
	}

	if len(actions) == 0 {
		le.Info("no actions produced from ingress object, doing nothing")
		return
	}

	controller.SendActionsToHandler(c.eventHandler, actions...)
}

func (c *Controller) ingressDelete(obj interface{}) {
	ingObj, ok := obj.(*extv1beta1.Ingress)
	if !ok {
		log.Error("unknown ingress object for informer delete event")
	}

	le := log.WithFields(log.Fields{
		"k8_event":        "delete",
		"k8_uid":          ingObj.UID,
		"k8_ingress_name": ingObj.Name,
		"k8_namespace":    ingObj.Namespace,
	})
	le.Info("delete ingress object")

	actions, err := actionsFromIngress(ingObj, controller.METHOD_DELETE)
	if err != nil {
		le.Errorf("failed creating actions from ingress object: %v", err)
		return
	}
	if len(actions) == 0 {
		le.Info("no actions produced from ingress object, doing nothing")
		return
	}

	controller.SendActionsToHandler(c.eventHandler, actions...)
}

func (c *Controller) ingressResync(ingObj *extv1beta1.Ingress) {
	le := log.WithFields(log.Fields{
		"k8_event":        "resync",
		"k8_uid":          ingObj.UID,
		"k8_ingress_name": ingObj.Name,
		"k8_namespace":    ingObj.Namespace,
	})
	le.Info("resync ingress object")

	actions, err := actionsFromIngress(ingObj, controller.METHOD_CHECK)
	if err != nil {
		le.Errorf("failed creating actions from ingress object: %v", err)
		return
	}

	if len(actions) == 0 {
		le.Info("no actions produced from ingress object, doing nothing")
		return
	}

	controller.SendActionsToHandler(c.eventHandler, actions...)
}

func (c *Controller) ingressUpdate(oObj, nObj interface{}) {
	ingOld, ok := oObj.(*extv1beta1.Ingress)
	if !ok {
		log.Error("unknown ingress object for informer update event")
	}

	ingNew, ok := nObj.(*extv1beta1.Ingress)
	if !ok {
		log.Error("unknown ingress object for informer ipdate event")
	}

	if ingOld.ResourceVersion == ingNew.ResourceVersion {
		c.ingressResync(ingNew)
		return
	}

	le := log.WithFields(log.Fields{
		"k8_event":        "update",
		"k8_uid":          ingNew.UID,
		"k8_ingress_name": ingNew.Name,
		"k8_namespace":    ingNew.Namespace,
	})
	le.Info("update ingress object")

	delActions, err := actionsFromIngress(ingOld, controller.METHOD_DELETE)
	if err != nil {
		le.Errorf("failed creating actions from ingress object: %v", err)
		return
	}

	if len(delActions) != 0 {
		controller.SendActionsToHandler(c.eventHandler, delActions...)
	}

	addActions, err := actionsFromIngress(ingNew, controller.METHOD_ADD)
	if err != nil {
		le.Errorf("failed creating actions from ingress object: %v", err)
		return
	}

	if len(addActions) != 0 {
		controller.SendActionsToHandler(c.eventHandler, addActions...)
	}
}

func NewController(lw *ListWatch, resync time.Duration) *Controller {
	c := &Controller{
		eventHandler: controller.EmptyActionHandler,
	}
	c.store, c.controller = cache.NewInformer(
		lw, &extv1beta1.Ingress{}, resync,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.ingressAdd,
			DeleteFunc: c.ingressDelete,
			UpdateFunc: c.ingressUpdate,
		})
	return c
}

func (c *Controller) Run(s <-chan struct{}) {
	c.controller.Run(s)
}

func (c *Controller) SetEventHandler(h controller.ActionEventHandler) {
	c.eventHandler = h
}
