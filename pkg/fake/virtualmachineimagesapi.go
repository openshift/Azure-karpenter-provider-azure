/*
Portions Copyright (c) Microsoft Corporation.

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

package fake

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/karpenter-provider-azure/pkg/providers/imagefamily/types"
)

type VirtualMachineImagesAPI struct {
	GetFunc func(context.Context, string, string, string, string, string, *armcompute.VirtualMachineImagesClientGetOptions) (armcompute.VirtualMachineImagesClientGetResponse, error)
}

var _ types.VirtualMachineImagesAPI = &VirtualMachineImagesAPI{}

func (v *VirtualMachineImagesAPI) Get(ctx context.Context, location, publisher, offer, sku, version string, options *armcompute.VirtualMachineImagesClientGetOptions) (armcompute.VirtualMachineImagesClientGetResponse, error) {
	if v.GetFunc != nil {
		return v.GetFunc(ctx, location, publisher, offer, sku, version, options)
	}
	return armcompute.VirtualMachineImagesClientGetResponse{}, fmt.Errorf("no marketplace image response configured for %s:%s:%s:%s in %s", publisher, offer, sku, version, location)
}

func (v *VirtualMachineImagesAPI) Reset() {
	v.GetFunc = nil
}
