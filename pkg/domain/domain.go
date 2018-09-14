package domain

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"

	"golang.org/x/net/idna"
)

type Domain string

const (
	empty     Domain = ""
	Root      rune   = '.'
	Separator rune   = '.'
)

func (d Domain) CompareStr(s string) bool {
	return strings.EqualFold(string(d), s)
}

func (d Domain) EndsWithRoot() bool {
	r, _ := utf8.DecodeLastRuneInString(string(d))
	return r == Root
}

func (d Domain) HashValue() (Domain, error) {
	hash, err := idna.Display.ToASCII(string(d))
	if err != nil {
		return empty, errors.Wrap(err, "domain hash value")
	}
	hash = strings.ToLower(hash)

	if !d.EndsWithRoot() {
		hash += string(Root)
	}
	return Domain(hash), nil
}

func (d Domain) IsEmpty() bool {
	return d == empty
}

func (d Domain) IsSubDomain(s Domain) (Domain, bool) {
	if d.IsEmpty() || s.IsEmpty() {
		return empty, false
	}
	if sp := strings.TrimSuffix(string(s), string(Separator)+string(d)); Domain(sp) != s {
		return Domain(sp), true
	}
	return empty, false
}

// Seperates the components of the domain name via it labels. If the domain name
// is empty or ends with a '.' then an empty slice is returned. Not other
// conversion is done on the labels.
func (d Domain) Labels() []string {
	return strings.FieldsFunc(string(d), func(c rune) bool {
		return c == Separator
	})
}
func (d Domain) Prepend(c ...string) Domain {
	return Domain(strings.Join(append(c, string(d)), string(Separator)))
}
