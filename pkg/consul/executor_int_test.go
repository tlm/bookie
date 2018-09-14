// +build integration

// go test -tags=integration

package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/tlmiller/bookie/pkg/domain"
	"github.com/tlmiller/bookie/pkg/k8/controller"
)

func TestUpsertAction(t *testing.T) {
	ex, err := NewExecutor(&api.Config{
		Address:    "localhost:8500",
		Datacenter: "dc1",
		Scheme:     "http",
	})
	if err != nil {
		t.Fatalf("unexpected error making consul executor: %v", err)
	}

	tests := []struct {
		Action      controller.Action
		ServiceName string
	}{
		{
			Action: controller.Action{
				domain.ResourceRecord{
					FQDN:  domain.Domain("testupserta1.service.dc1.consul"),
					ID:    "ingress/tesupserta1",
					Type:  domain.AAAA,
					Value: "fe80:1::1",
				},
				controller.METHOD_ADD,
			},
			ServiceName: "testupserta1",
		},
		{
			Action: controller.Action{
				domain.ResourceRecord{
					FQDN:  domain.Domain("testupserta1.service.dc1.consul"),
					ID:    "ingress/tesupserta1",
					Type:  domain.A,
					Value: "192.168.0.1",
				},
				controller.METHOD_ADD,
			},
			ServiceName: "testupserta1",
		},
	}

	for _, test := range tests {
		err = ex.upsertAction(&test.Action)
		if err != nil {
			t.Fatalf("unexpected error upserting new consul record %s: %v",
				test.Action.FQDN, err)
		}

		qs, _, err := ex.client.Catalog().Service(test.ServiceName, "", &api.QueryOptions{
			Datacenter: ex.Datacenter(),
		})

		if err != nil {
			t.Fatalf("unexpected error querying for upserted consul record %s: %v",
				test.Action.FQDN, err)
		}

		found := false
		for _, s := range qs {
			if s.ServiceID != fmt.Sprintf("%s-%s", test.Action.ID, string(test.Action.Type)) {
				continue
			}
			if found {
				t.Fatalf("found more then one service for upsert %s", test.Action.FQDN)
			}
			found = true

			if s.Node != test.ServiceName {
				t.Errorf("consul service upsert record has incorrect node name, expected %s got %s",
					test.ServiceName, s.Node)
			}
			if s.ServiceName != test.ServiceName {
				t.Errorf("consul service upsert record has incorrect service name, expected %s got %s",
					test.ServiceName, s.ServiceName)
			}
			if s.ServiceAddress != test.Action.Value {
				t.Errorf("consul service upsert record has incorrect service address, expected %s got %s",
					test.Action.Value, s.ServiceAddress)
			}
		}
	}
}
