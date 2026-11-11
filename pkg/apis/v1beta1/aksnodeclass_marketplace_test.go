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
)

func TestHasMarketplaceImages(t *testing.T) {
	cases := map[string]struct {
		images []string
		want   bool
	}{
		"omitted": {},
		"empty":   {images: []string{}},
		"one":     {images: []string{"publisher:offer:sku:1.0.0"}, want: true},
		"multiple": {images: []string{
			"publisher:offer:amd64:1.0.0", "publisher:offer:arm64:1.0.0",
		}, want: true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)
			nodeClass := &v1beta1.AKSNodeClass{Spec: v1beta1.AKSNodeClassSpec{MarketplaceImages: tc.images}}
			g.Expect(nodeClass.HasMarketplaceImages()).To(Equal(tc.want))
		})
	}
}

func TestParseMarketplaceImageURN(t *testing.T) {
	cases := map[string]struct {
		urn  string
		want *v1beta1.MarketplaceImage
	}{
		"explicit version": {urn: "publisher:offer:sku:1.0.0", want: &v1beta1.MarketplaceImage{
			Publisher: "publisher", Offer: "offer", SKU: "sku", Version: "1.0.0",
		}},
		"nonnumeric version": {urn: "publisher:offer:sku:some-version", want: &v1beta1.MarketplaceImage{
			Publisher: "publisher", Offer: "offer", SKU: "sku", Version: "some-version",
		}},
		"latest":            {urn: "publisher:offer:sku:latest"},
		"newline":           {urn: "publisher:offer:sku:1.0.0\n"},
		"space":             {urn: "pub lisher:offer:sku:1.0.0"},
		"tab":               {urn: "publisher:of\tfer:sku:1.0.0"},
		"empty":             {},
		"missing component": {urn: "publisher:offer:sku"},
		"extra component":   {urn: "publisher:offer:sku:version:extra"},
		"empty publisher":   {urn: ":offer:sku:version"},
		"empty offer":       {urn: "publisher::sku:version"},
		"empty sku":         {urn: "publisher:offer::version"},
		"empty version":     {urn: "publisher:offer:sku:"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)
			image, err := v1beta1.ParseMarketplaceImageURN(tc.urn)
			if tc.want == nil {
				g.Expect(err).To(HaveOccurred())
				g.Expect(image).To(BeNil())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(image).To(Equal(tc.want))
			}
		})
	}
}
