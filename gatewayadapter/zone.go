package gatewayadapter

import (
	"context"

	"github.com/dnsgateway/dnsprovider"
)

// CreateZone creates a new DNS zone.
func (a *adapter) CreateZone(context.Context, string) (dnsprovider.Zone, error) {
	return dnsprovider.Zone{}, dnsprovider.ErrNotImplemented
}

// UpdateZone updates an existing DNS zone.
func (a *adapter) UpdateZone(context.Context, dnsprovider.Zone, dnsprovider.ZoneUpdate) (dnsprovider.ZoneStatus, error) {
	return dnsprovider.ZoneStatus{}, dnsprovider.ErrNotImplemented
}

// DeleteZone deletes an existing DNS zone.
func (a *adapter) DeleteZone(context.Context, dnsprovider.Zone) error {
	return dnsprovider.ErrNotImplemented
}
