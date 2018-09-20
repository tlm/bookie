package domain

import (
	"testing"
)

func TestFromlabels(t *testing.T) {
	tests := []struct {
		Labels []string
		Value  Domain
	}{
		{
			[]string{"example", "com"},
			Domain("example.com."),
		},
		{
			[]string{"test", "example", "com"},
			Domain("test.example.com."),
		},
		{
			[]string{"Www", "tESt", "台灣"},
			Domain("Www.tESt.台灣."),
		},
	}

	for _, test := range tests {
		if o := FromLabels(test.Labels...); o != test.Value {
			t.Errorf("expected output of labels %s == %s", o, test.Value)
		}
	}
}

func TestHashValue(t *testing.T) {
	tests := []struct {
		Domain string
		Hash   string
	}{
		{"example.com.", "example.com."},
		{"example.com", "example.com."},
		{"ExamPle.com", "example.com."},
		{"www.test.台灣.", "www.test.xn--kpry57d."},
		{"Www.tESt.台灣.", "www.test.xn--kpry57d."},
		{"Www.tESt.台灣", "www.test.xn--kpry57d."},
		{"www.test.台灣", "www.test.xn--kpry57d."},
	}

	for _, test := range tests {
		d := Domain(test.Domain)
		var (
			err error
			h   Domain
		)
		if h, err = d.HashValue(); err != nil {
			t.Errorf("unexpected hash value error: %v", err)
		} else if string(h) != test.Hash {
			t.Errorf("hash value for '%s' expected '%s', got '%s'",
				d, test.Hash, h)
		}
	}
}

func TestLabelSplitting(t *testing.T) {
	tests := []struct {
		Domain string
		Labels []string
	}{
		{"example.com", []string{"example", "com"}},
		{"example.com.", []string{"example", "com"}},
		{".", []string{}},
		{"..", []string{}},
		{"", []string{}},
		{"Www.tESt.台灣.", []string{"Www", "tESt", "台灣"}},
	}

	for _, test := range tests {
		labels := Domain(test.Domain).Labels()
		if len(labels) != len(test.Labels) {
			t.Fatalf("labels len for %s, expected %d got %d", test.Domain,
				len(test.Labels), len(labels))
		}
		for index, l := range labels {
			if l != test.Labels[index] {
				t.Fatalf("label at index %d, expected %s got %s", index,
					test.Labels[index], l)
			}
		}
	}
}

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
