package domain

import (
	"strings"
)

type Domain string

const (
	empty     Domain = ""
	SEPARATOR string = "."
)

func (d Domain) CompareStr(s string) bool {
	return strings.EqualFold(string(d), s)
}

func (d Domain) IsEmpty() bool {
	return d == empty
}

func (d Domain) IsSubDomain(s Domain) (Domain, bool) {
	if d.IsEmpty() || s.IsEmpty() {
		return empty, false
	}
	if sp := strings.TrimSuffix(string(s), SEPARATOR+string(d)); Domain(sp) != s {
		return Domain(sp), true
	}
	return empty, false
}

func (d Domain) Prepend(c ...string) Domain {
	return Domain(strings.Join(append(c, string(d)), SEPARATOR))
}
