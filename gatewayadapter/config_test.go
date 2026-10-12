/*
Copyright 2026 DNS Gateway Team.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gatewayadapter_test

import (
	"slices"
	"testing"

	"github.com/dnsgateway/dnsprovider"

	"github.com/dnsgateway/dnsgate-go/gatewayadapter"
)

func TestFactoryDeclaresGatewayTokenAndEgress(t *testing.T) {
	f := gatewayadapter.Factory()

	want := []dnsprovider.Field{
		{Name: "gateway", Label: "Gateway address"},
		{Name: "token", Label: "Service Token", Secret: true},
		{Name: "egress", Label: "Egress", Default: "auto", Choices: []string{"auto", "gateway", "agent"}},
	}

	if !slices.EqualFunc(f.Fields, want, func(a, b dnsprovider.Field) bool {
		return a.Name == b.Name && a.Label == b.Label && a.Secret == b.Secret && a.Tenant == b.Tenant &&
			a.Default == b.Default && slices.Equal(a.Choices, b.Choices)
	}) {
		t.Errorf("Fields = %+v, want %+v", f.Fields, want)
	}

	if want := []dnsprovider.Capability{dnsprovider.CapBatcher}; !slices.Equal(f.Capabilities, want) {
		t.Errorf("Capabilities = %v, want %v: the Gateway batches and has no Zone writes", f.Capabilities, want)
	}
}

func TestNewMakesNoNetworkCall(t *testing.T) {
	// Nothing listens on port 1: building succeeds all the same.
	p, err := gatewayadapter.Factory().New(map[string]string{"gateway": "127.0.0.1:1", "token": "tok"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if p == nil {
		t.Fatal("New returned no Provider")
	}

	if !dnsprovider.Can(p, dnsprovider.CapBatcher) {
		t.Error("the Gateway Adapter is no Batcher")
	}
}

func TestNewRefusesBadFields(t *testing.T) {
	for name, fields := range map[string]map[string]string{
		"unknown field":   {"gateway": "127.0.0.1:1", "token": "tok", "tokn": "x"},
		"missing token":   {"gateway": "127.0.0.1:1"},
		"missing gateway": {"token": "tok"},
		"unknown egress":  {"gateway": "127.0.0.1:1", "token": "tok", "egress": "direct"},
	} {
		t.Run(name, func(t *testing.T) {
			p, err := gatewayadapter.Factory().New(fields)
			if err == nil {
				t.Fatalf("New(%v) = %v, want an error", fields, p)
			}

			if p != nil {
				t.Errorf("New returned a Provider with its error: %v", p)
			}
		})
	}
}

func TestNewTakesEachEgress(t *testing.T) {
	for _, egress := range []string{"", "auto", "gateway", "agent"} {
		if _, err := gatewayadapter.Factory().New(map[string]string{"gateway": "127.0.0.1:1", "token": "tok", "egress": egress}); err != nil {
			t.Errorf("New with egress %q: %v", egress, err)
		}
	}
}
