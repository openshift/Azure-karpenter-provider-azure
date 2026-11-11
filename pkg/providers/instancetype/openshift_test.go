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

package instancetype

import (
	"testing"

	"github.com/Azure/karpenter-provider-azure/pkg/apis/v1beta1"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
)

// OpenShift nodes always run OVN-Kubernetes and must get the OpenShift default of 250 pods per
// node, independent of the AKS CNI options that Karpenter is configured with. Karpenter sizes
// NodeClaims from this value.
func TestOpenShiftMaxPods(t *testing.T) {
	cases := map[string]struct {
		nodeClass *v1beta1.AKSNodeClass
		expected  int32
	}{
		"defaults to 250 pods": {
			nodeClass: &v1beta1.AKSNodeClass{},
			expected:  250,
		},
		"honors the node class maxPods override": {
			nodeClass: &v1beta1.AKSNodeClass{Spec: v1beta1.AKSNodeClassSpec{MaxPods: lo.ToPtr(int32(100))}},
			expected:  100,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(openShiftMaxPods(c.nodeClass, "openshift", "", "")).To(Equal(c.expected))
		})
	}
}
