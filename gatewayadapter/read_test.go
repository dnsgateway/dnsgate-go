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

package gatewayadapter

import (
	"context"
	"net"
	"testing"

	"github.com/dnsgateway/dnsprovider"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	dnsgatev1 "github.com/dnsgateway/dnsgate-go/api/dnsgate/v1"
	"github.com/dnsgateway/dnsgate-go/client"
)

// listingGateway answers every Call with the list response it holds for
// the request's operation.
type listingGateway struct {
	dnsgatev1.UnimplementedGatewayServiceServer

	zones   *dnsgatev1.ListZonesResponse
	records *dnsgatev1.ListRecordsResponse
}

func (g *listingGateway) Call(stream grpc.BidiStreamingServer[dnsgatev1.CallRequest, dnsgatev1.CallResponse]) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	resp := &dnsgatev1.Response{}

	switch msg.GetRequest().GetOperation().(type) {
	case *dnsgatev1.Request_ListZones:
		resp.Result = &dnsgatev1.Response_ListZones{ListZones: g.zones}
	case *dnsgatev1.Request_ListRecords:
		resp.Result = &dnsgatev1.Response_ListRecords{ListRecords: g.records}
	}

	return stream.Send(&dnsgatev1.CallResponse{Message: &dnsgatev1.CallResponse_Response{Response: resp}})
}

// adapterFor returns a Gateway Adapter that calls gw.
func adapterFor(t *testing.T, gw dnsgatev1.GatewayServiceServer) *adapter {
	t.Helper()

	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer()
	dnsgatev1.RegisterGatewayServiceServer(srv, gw)

	go func() { _ = srv.Serve(lis) }()

	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = conn.Close() })

	return &adapter{client: client.New(conn), token: "tok"}
}

func TestRecordsCarryThePartialViewFlag(t *testing.T) {
	for name, partial := range map[string]bool{
		"partial view": true,
		"full view":    false,
	} {
		t.Run(name, func(t *testing.T) {
			gw := &listingGateway{records: &dnsgatev1.ListRecordsResponse{
				RecordSets: []*dnsgatev1.RecordSet{{Name: "www", Type: "A", Ttl: 300, Records: []*dnsgatev1.RecordValue{{Value: "192.0.2.1"}}}},
				Partial:    partial,
			}}

			got, err := adapterFor(t, gw).Records(t.Context(), dnsprovider.Zone{ID: "example.com", Name: "example.com"}, dnsprovider.RecordFilter{})
			if err != nil {
				t.Fatalf("Records: %v", err)
			}

			if got.Partial != partial {
				t.Errorf("Records of a response with partial %v: Partial = %v", partial, got.Partial)
			}

			if len(got.Records) != 1 || got.Records[0].Name != "www.example.com" {
				t.Errorf("Records = %+v, want the one www.example.com record", got.Records)
			}
		})
	}
}

func TestZonesCarryThePartialViewFlag(t *testing.T) {
	for name, partial := range map[string]bool{
		"partial view": true,
		"full view":    false,
	} {
		t.Run(name, func(t *testing.T) {
			gw := &listingGateway{zones: &dnsgatev1.ListZonesResponse{
				Zones:   []*dnsgatev1.Zone{{Name: "example.com", State: dnsgatev1.Zone_STATE_ACTIVE}},
				Partial: partial,
			}}

			got, err := adapterFor(t, gw).Zones(t.Context(), dnsprovider.ZoneFilter{})
			if err != nil {
				t.Fatalf("Zones: %v", err)
			}

			if got.Partial != partial {
				t.Errorf("Zones of a response with partial %v: Partial = %v", partial, got.Partial)
			}

			if len(got.Zones) != 1 || got.Zones[0].Name != "example.com" {
				t.Errorf("Zones = %+v, want example.com alone", got.Zones)
			}
		})
	}
}
