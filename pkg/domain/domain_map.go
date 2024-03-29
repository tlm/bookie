package domain

import (
	"github.com/pkg/errors"
)

type DomainMap map[Domain]interface{}

func NewDomainMap() DomainMap {
	return make(DomainMap, 0)
}

func (d DomainMap) Get(key Domain) (interface{}, bool) {
	h, err := key.HashValue()
	if err != nil {
		return nil, false
	}
	val, has := d[h]
	return val, has
}

func (d DomainMap) GetNearestRoot(key Domain) (interface{}, bool) {
	h, err := key.HashValue()
	if err != nil {
		return nil, false
	}

	labels := h.Labels()
	for i, _ := range labels {
		r := FromLabels(labels[i:]...)
		val, has := d[r]
		if has {
			return val, has
		}
	}
	return nil, false
}

func (d DomainMap) Set(key Domain, data interface{}) error {
	h, err := key.HashValue()
	if err != nil {
		return errors.Wrap(err, "domain map getting domain hash value")
	}
	d[h] = data
	return nil
}
