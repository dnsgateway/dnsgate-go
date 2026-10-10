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

// Package gatewayadapter is the Gateway Adapter: it presents the Gateway
// to DNS tools built on the shared Provider interface
// (github.com/dnsgateway/dnsprovider), such as dnscli, so that they manage
// records through the Gateway as they would through a Provider. The
// Gateway itself is never a Provider (dnsgate ADR 0008).
//
// A Provider the Factory builds holds a Gateway address and a Service
// Token, and sends every call as one Gateway request authorized by that
// token, acting as Egress Agent when the Gateway relays through Agent
// Egress. Zones have their name as ID. Each value of an RRset is one
// Record, with the RRset's TTL, its own Record Extras and an ID built from
// its content. Writes go through ApplyChangeset, which writes each RRset
// a Changeset touches whole, with create, update or delete, and names the
// Expected RRset on each update and delete.
//
// The Provider's TTL range, automatic TTL and Record Extras are those of
// the token's Provider type, which the Gateway names on InspectToken.
//
// A Bound Token is not usable: the Gateway Adapter sends no Identity
// Token.
package gatewayadapter
