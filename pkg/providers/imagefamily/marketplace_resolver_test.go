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

package imagefamily_test

import (
	"testing"

	"github.com/Azure/karpenter-provider-azure/pkg/apis/v1beta1"
	"github.com/Azure/karpenter-provider-azure/pkg/providers/imagefamily"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/karpenter/pkg/cloudprovider"
	"sigs.k8s.io/karpenter/pkg/scheduling"
)

func TestResolveMarketplaceImageFromNodeClass(t *testing.T) {
	amd64 := v1beta1.NodeImage{ID: "publisher:offer:amd64:1.0.0", Requirements: []corev1.NodeSelectorRequirement{
		{Key: corev1.LabelArchStable, Operator: corev1.NodeSelectorOpIn, Values: []string{"amd64"}},
		{Key: v1beta1.LabelSKUHyperVGeneration, Operator: corev1.NodeSelectorOpIn, Values: []string{"2"}},
	}}
	arm64 := v1beta1.NodeImage{ID: "publisher:offer:arm64:1.0.0", Requirements: []corev1.NodeSelectorRequirement{
		{Key: corev1.LabelArchStable, Operator: corev1.NodeSelectorOpIn, Values: []string{"arm64"}},
		{Key: v1beta1.LabelSKUHyperVGeneration, Operator: corev1.NodeSelectorOpIn, Values: []string{"2"}},
	}}
	gen1 := *amd64.DeepCopy()
	gen1.Requirements[1].Values = []string{"1"}
	cases := map[string]struct {
		arch              string
		generations       []string
		missingGeneration bool
		images            []v1beta1.NodeImage
		mutate            func(*v1beta1.AKSNodeClass)
		want              string
	}{
		"Gen2 image rejects Gen1-only size":       {arch: "amd64", generations: []string{"1"}, images: []v1beta1.NodeImage{amd64}},
		"Gen2 image accepts Gen2-only size":       {arch: "amd64", generations: []string{"2"}, images: []v1beta1.NodeImage{amd64}, want: "publisher:offer:amd64:1.0.0"},
		"Gen2 image accepts dual-generation size": {arch: "amd64", generations: []string{"1", "2"}, images: []v1beta1.NodeImage{amd64}, want: "publisher:offer:amd64:1.0.0"},
		"Gen1 image accepts Gen1-only size":       {arch: "amd64", generations: []string{"1"}, images: []v1beta1.NodeImage{gen1}, want: "publisher:offer:amd64:1.0.0"},
		"Gen1 image rejects Gen2-only size":       {arch: "amd64", generations: []string{"2"}, images: []v1beta1.NodeImage{gen1}},
		"Gen1 image accepts dual-generation size": {arch: "amd64", generations: []string{"1", "2"}, images: []v1beta1.NodeImage{gen1}, want: "publisher:offer:amd64:1.0.0"},
		"missing size generation capability":      {arch: "amd64", missingGeneration: true, images: []v1beta1.NodeImage{amd64}},
		"amd64":                                   {arch: "amd64", images: []v1beta1.NodeImage{arm64, amd64}, want: "publisher:offer:amd64:1.0.0"},
		"arm64":                                   {arch: "arm64", images: []v1beta1.NodeImage{amd64, arm64}, want: "publisher:offer:arm64:1.0.0"},
		"reversed":                                {arch: "arm64", images: []v1beta1.NodeImage{arm64, amd64}, want: "publisher:offer:arm64:1.0.0"},
		"no match":                                {arch: "arm64", images: []v1beta1.NodeImage{amd64}},
		"empty":                                   {arch: "arm64"},
		"unconstrained legacy status":             {arch: "arm64", images: []v1beta1.NodeImage{{ID: "publisher:offer:sku:1.0.0"}}},
		"stale":                                   {arch: "amd64", images: []v1beta1.NodeImage{amd64}, mutate: func(n *v1beta1.AKSNodeClass) { n.Generation++ }},
		"unready": {arch: "amd64", images: []v1beta1.NodeImage{amd64}, mutate: func(n *v1beta1.AKSNodeClass) {
			n.StatusConditions().SetFalse(v1beta1.ConditionTypeImagesReady, "Failure", "lookup failed")
		}},
		"unknown":                  {arch: "amd64", images: []v1beta1.NodeImage{amd64}, mutate: func(n *v1beta1.AKSNodeClass) { n.Status.Conditions = nil }},
		"malformed URN":            {arch: "amd64", images: []v1beta1.NodeImage{{ID: "publisher:offer:sku", Requirements: amd64.Requirements}}},
		"unsupported architecture": {arch: "unknown", images: []v1beta1.NodeImage{{ID: "publisher:offer:sku:1", Requirements: []corev1.NodeSelectorRequirement{{Key: corev1.LabelArchStable, Operator: corev1.NodeSelectorOpIn, Values: []string{"unknown"}}}}}},
		"status differs from spec": {arch: "amd64", images: []v1beta1.NodeImage{amd64}, want: "publisher:offer:amd64:1.0.0"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)
			nodeClass := &v1beta1.AKSNodeClass{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Spec:       v1beta1.AKSNodeClassSpec{MarketplaceImages: []string{"different:spec:sku:2"}},
				Status:     v1beta1.AKSNodeClassStatus{Images: tc.images},
			}
			nodeClass.StatusConditions().SetTrue(v1beta1.ConditionTypeImagesReady)
			if tc.mutate != nil {
				tc.mutate(nodeClass)
			}
			instanceType := &cloudprovider.InstanceType{Name: "test", Requirements: scheduling.NewRequirements(scheduling.NewRequirement(corev1.LabelArchStable, corev1.NodeSelectorOpIn, tc.arch))}
			if tc.missingGeneration {
				instanceType.Requirements.Add(scheduling.NewRequirement(v1beta1.LabelSKUHyperVGeneration, corev1.NodeSelectorOpDoesNotExist))
			} else {
				generations := tc.generations
				if generations == nil {
					generations = []string{"1", "2"}
				}
				instanceType.Requirements.Add(scheduling.NewRequirement(v1beta1.LabelSKUHyperVGeneration, corev1.NodeSelectorOpIn, generations...))
			}
			id, err := imagefamily.ResolveMarketplaceImageFromNodeClass(nodeClass, instanceType)
			if tc.want == "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(id).To(BeEmpty())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(id).To(Equal(tc.want))
			}
		})
	}
}
