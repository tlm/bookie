package domain

import (
	"testing"
)

type simpleVal int

func TestDomainMapSetExists(t *testing.T) {
	tests := []struct {
		Key   string
		Value simpleVal
	}{
		{"example.com.", 1},
		{"example.com", 2},
		{"ExamPle.com", 3},
		{"www.test.台灣.", 4},
		{"Www.tESt.台灣.", 5},
		{"Www.tESt.台灣", 6},
		{"www.test.台灣", 7},
	}

	for _, test := range tests {
		dm := NewDomainMap()
		err := dm.Set(Domain(test.Key), test.Value)
		if err != nil {
			t.Fatalf("unexpected error set domain map key: %v", err)
		}

		val, exists := dm.Get(Domain(test.Key))
		if exists == false {
			t.Fatal("set key in domain map does not exist")
		}

		if rval, isSV := val.(simpleVal); !isSV {
			t.Fatalf("wrong value typed recieved for key")
		} else if rval != test.Value {
			t.Fatalf("domain map get for key does not mat set value: %v != %v",
				rval, test.Value)
		}
	}
}

func TestDomainMapGet(t *testing.T) {
	d := Domain("example.com")
	dm := NewDomainMap()
	_, exists := dm.Get(d)
	if exists {
		t.Fatalf("got failse positive for domain map get with no keys")
	}

	err := dm.Set(d, simpleVal(1))
	if err != nil {
		t.Fatalf("unexpected error fot domain map set: %v", err)
	}

	_, exists = dm.Get(d)
	if !exists {
		t.Fatalf("failed to get value for key in domain map, does not exist")
	}
}
