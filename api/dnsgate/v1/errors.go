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

package dnsgatev1

import (
	"strings"

	"google.golang.org/grpc/status"
)

// ErrorDetail returns the Error detail of st, or nil if it has none.
func ErrorDetail(st *status.Status) *Error {
	for _, d := range st.Details() {
		if e, ok := d.(*Error); ok {
			return e
		}
	}

	return nil
}

// Name returns the category's enum name in lower case without the
// prefix, such as "unauthenticated" or "permission_denied".
func (c Error_Category) Name() string {
	return strings.ToLower(strings.TrimPrefix(c.String(), "CATEGORY_"))
}
