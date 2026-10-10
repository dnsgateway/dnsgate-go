package gatewayadapter

import (
	"testing"

	dnsgatev1 "github.com/dnsgateway/dnsgate-go/api/dnsgate/v1"
)

func TestMayListZonesFollowsTheGatewaysRules(t *testing.T) {
	for name, tc := range map[string]struct {
		perm   *dnsgatev1.Permission
		egress dnsgatev1.Egress
		want   bool
	}{
		"list on zones":                {&dnsgatev1.Permission{Verbs: []string{"list"}, Resources: []string{"zones"}}, dnsgatev1.Egress_EGRESS_UNSPECIFIED, true},
		"resources left out":           {&dnsgatev1.Permission{Verbs: []string{"list"}}, dnsgatev1.Egress_EGRESS_UNSPECIFIED, false},
		"no list verb":                 {&dnsgatev1.Permission{Verbs: []string{"get"}, Resources: []string{"zones"}}, dnsgatev1.Egress_EGRESS_UNSPECIFIED, false},
		"zones filter does not count":  {&dnsgatev1.Permission{Verbs: []string{"list"}, Resources: []string{"zones"}, Zones: []string{"example.com"}}, dnsgatev1.Egress_EGRESS_UNSPECIFIED, true},
		"egress filter left out":       {&dnsgatev1.Permission{Verbs: []string{"list"}, Resources: []string{"zones"}}, dnsgatev1.Egress_EGRESS_AGENT, true},
		"egress filter holds egress":   {&dnsgatev1.Permission{Verbs: []string{"list"}, Resources: []string{"zones"}, Egress: []string{"agent"}}, dnsgatev1.Egress_EGRESS_AGENT, true},
		"egress filter lacks egress":   {&dnsgatev1.Permission{Verbs: []string{"list"}, Resources: []string{"zones"}, Egress: []string{"agent"}}, dnsgatev1.Egress_EGRESS_GATEWAY, false},
		"auto egress with egress rule": {&dnsgatev1.Permission{Verbs: []string{"list"}, Resources: []string{"zones"}, Egress: []string{"agent"}}, dnsgatev1.Egress_EGRESS_UNSPECIFIED, true},
	} {
		t.Run(name, func(t *testing.T) {
			a := &adapter{egress: tc.egress}
			if got := a.mayListZones([]*dnsgatev1.Permission{tc.perm}); got != tc.want {
				t.Errorf("mayListZones(%v) through %v = %v, want %v", tc.perm, tc.egress, got, tc.want)
			}
		})
	}
}
