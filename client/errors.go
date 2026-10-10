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

package client

import (
	"strings"

	"google.golang.org/grpc/status"

	dnsgatev1 "github.com/dnsgateway/dnsgate-go/api/dnsgate/v1"
)

// Error is a failed Gateway request. Its gRPC status stays reachable
// through status.FromError.
type Error struct {
	status *status.Status
	detail *dnsgatev1.Error
}

// newError wraps the error a stream ended with. Errors that carry no gRPC
// status, such as a canceled context, are returned as they are.
func newError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	return &Error{status: st, detail: dnsgatev1.ErrorDetail(st)}
}

// Error returns the Gateway error category and message, such as
// "unauthenticated: wrong or missing bootstrap secret". An error without a
// Gateway category, such as a connection failure, shows its gRPC code.
func (e *Error) Error() string {
	if e.detail == nil {
		return strings.ToLower(e.status.Code().String()) + ": " + e.status.Message()
	}

	return strings.ReplaceAll(e.detail.GetCategory().Name(), "_", " ") + ": " + e.status.Message()
}

// Detail returns the Gateway's Error detail, or nil if there was none.
func (e *Error) Detail() *dnsgatev1.Error { return e.detail }

// GRPCStatus returns the status the stream ended with.
func (e *Error) GRPCStatus() *status.Status { return e.status }
