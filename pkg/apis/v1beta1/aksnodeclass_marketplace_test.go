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

package v1beta1_test

import (
	"testing"

	"github.com/Azure/karpenter-provider-azure/pkg/apis/v1beta1"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
)

func TestHasMarketplaceImage(t *testing.T) {
	cases := []struct {
		name     string
		image    *v1beta1.MarketplaceImage
		expected bool
	}{
		{
			name:     "nil marketplace image",
			image:    nil,
			expected: false,
		},
		{
			name: "complete marketplace image",
			image: &v1beta1.MarketplaceImage{
				Publisher: lo.ToPtr("azureopenshift"),
				Offer:     lo.ToPtr("aro4"),
				SKU:       lo.ToPtr("aro_422-v2"),
				Version:   lo.ToPtr("9.8.20260428"),
			},
			expected: true,
		},
		{
			name: "missing version",
			image: &v1beta1.MarketplaceImage{
				Publisher: lo.ToPtr("azureopenshift"),
				Offer:     lo.ToPtr("aro4"),
				SKU:       lo.ToPtr("aro_422-v2"),
			},
			expected: false,
		},
		{
			name: "empty publisher",
			image: &v1beta1.MarketplaceImage{
				Publisher: lo.ToPtr(""),
				Offer:     lo.ToPtr("aro4"),
				SKU:       lo.ToPtr("aro_422-v2"),
				Version:   lo.ToPtr("9.8.20260428"),
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			nodeClass := &v1beta1.AKSNodeClass{
				Spec: v1beta1.AKSNodeClassSpec{
					MarketplaceImage: tc.image,
				},
			}
			g.Expect(nodeClass.HasMarketplaceImage()).To(Equal(tc.expected))
		})
	}
}

func TestMarketplaceImageURN(t *testing.T) {
	cases := []struct {
		name     string
		image    *v1beta1.MarketplaceImage
		expected string
	}{
		{
			name:     "nil marketplace image",
			image:    nil,
			expected: "",
		},
		{
			name: "complete marketplace image",
			image: &v1beta1.MarketplaceImage{
				Publisher: lo.ToPtr("azureopenshift"),
				Offer:     lo.ToPtr("aro4"),
				SKU:       lo.ToPtr("aro_422-v2"),
				Version:   lo.ToPtr("9.8.20260428"),
			},
			expected: "azureopenshift:aro4:aro_422-v2:9.8.20260428",
		},
		{
			name: "incomplete marketplace image",
			image: &v1beta1.MarketplaceImage{
				Publisher: lo.ToPtr("azureopenshift"),
				Offer:     lo.ToPtr("aro4"),
				SKU:       lo.ToPtr("aro_422-v2"),
			},
			expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			nodeClass := &v1beta1.AKSNodeClass{
				Spec: v1beta1.AKSNodeClassSpec{
					MarketplaceImage: tc.image,
				},
			}
			g.Expect(nodeClass.MarketplaceImageURN()).To(Equal(tc.expected))
		})
	}
}
