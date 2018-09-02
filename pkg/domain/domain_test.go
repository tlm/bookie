package domain

import (
	"testing"
)

func TestIsSubDomain(t *testing.T) {
	base, sub := Domain("example.com"), Domain("tst.example.com")
	part, check := base.IsSubDomain(sub)
	if check == false {
		t.Errorf("%s was not a subdomain of %s", sub, base)
	}
	if part != "tst" {
		t.Error("expected subdomain part to be tst")
	}
}
