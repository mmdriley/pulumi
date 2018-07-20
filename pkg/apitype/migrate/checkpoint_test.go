// Copyright 2016-2018, Pulumi Corporation.
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

package migrate

import (
	"testing"

	"github.com/pulumi/pulumi/pkg/apitype"
	"github.com/pulumi/pulumi/pkg/resource/config"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/stretchr/testify/assert"
)

func TestCheckpointV1ToV2(t *testing.T) {
	v1 := apitype.CheckpointV1{
		Stack: tokens.QName("mystack"),
		Config: config.Map{
			config.MustMakeKey("foo", "number"): config.NewValue("42"),
		},
		Latest: &apitype.DeploymentV1{
			Manifest:  apitype.ManifestV1{},
			Resources: []apitype.ResourceV1{},
		},
	}

	v2 := UpToCheckpointV2(v1)
	assert.Equal(t, tokens.QName("mystack"), v2.Stack)
	assert.Equal(t, config.Map{
		config.MustMakeKey("foo", "number"): config.NewValue("42"),
	}, v2.Config)
	assert.Len(t, v2.Latest.Resources, 0)
}

func TestCheckpointV2ToV1(t *testing.T) {
	v2 := apitype.CheckpointV2{
		Stack: tokens.QName("mystack"),
		Config: config.Map{
			config.MustMakeKey("foo", "number"): config.NewValue("42"),
		},
		Latest: &apitype.DeploymentV2{
			Manifest:  apitype.ManifestV1{},
			Resources: []apitype.ResourceV2{},
		},
	}

	v1 := DownToCheckpointV1(v2)
	assert.Equal(t, tokens.QName("mystack"), v1.Stack)
	assert.Equal(t, config.Map{
		config.MustMakeKey("foo", "number"): config.NewValue("42"),
	}, v1.Config)
	assert.Len(t, v1.Latest.Resources, 0)
}
