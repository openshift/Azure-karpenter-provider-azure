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

package fake_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/karpenter-provider-azure/pkg/fake"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
)

func TestVirtualMachineImagesAPI(t *testing.T) {
	g := NewWithT(t)
	api := &fake.VirtualMachineImagesAPI{}
	ctx := context.Background()
	opts := &armcompute.VirtualMachineImagesClientGetOptions{}
	wantErr := errors.New("lookup failed")
	api.GetFunc = func(gotCtx context.Context, location, publisher, offer, sku, version string, options *armcompute.VirtualMachineImagesClientGetOptions) (armcompute.VirtualMachineImagesClientGetResponse, error) {
		g.Expect(gotCtx).To(Equal(ctx))
		g.Expect([]string{location, publisher, offer, sku, version}).To(Equal([]string{"eastus", "publisher", "offer", "sku", "latest"}))
		g.Expect(options).To(BeIdenticalTo(opts))
		return armcompute.VirtualMachineImagesClientGetResponse{VirtualMachineImage: armcompute.VirtualMachineImage{
			Properties: &armcompute.VirtualMachineImageProperties{Architecture: lo.ToPtr(armcompute.ArchitectureTypesArm64)},
		}}, nil
	}
	response, err := api.Get(ctx, "eastus", "publisher", "offer", "sku", "latest", opts)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(response.Properties.Architecture).To(HaveValue(Equal(armcompute.ArchitectureTypesArm64)))
	api.GetFunc = func(context.Context, string, string, string, string, string, *armcompute.VirtualMachineImagesClientGetOptions) (armcompute.VirtualMachineImagesClientGetResponse, error) {
		return armcompute.VirtualMachineImagesClientGetResponse{}, wantErr
	}
	_, err = api.Get(ctx, "eastus", "publisher", "offer", "sku", "latest", nil)
	g.Expect(err).To(MatchError(wantErr))
	api.Reset()
	_, err = api.Get(ctx, "eastus", "publisher", "offer", "sku", "latest", nil)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err).NotTo(MatchError(wantErr))
}
