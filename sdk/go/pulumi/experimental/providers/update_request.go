// Copyright 2016-2023, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "generate"; DO NOT EDIT.

//nolint:lll
package providers

import (
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// A request to update a resource.
type UpdateRequest struct {
	// The ID of the resource to update.
	Id string
	// The Pulumi URN for this resource.
	URN string
	// The Pulumi type for this resource.
	Type string
	// The Pulumi name for this resource.
	Name string
	// The old values of provider inputs for the resource to update.
	Old resource.PropertyMap
	// The new values of provider inputs for the resource to update.
	New resource.PropertyMap
	// The create request timeout.
	Timeout time.Duration
	// A set of property paths that should be treated as unchanged.
	IgnoreChanges []string
	// true if this is a preview and the provider should not actually update the resource.
	Preview bool
}
