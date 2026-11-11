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
	"github.com/Azure/karpenter-provider-azure/pkg/apis/v1beta1"
	"github.com/Azure/karpenter-provider-azure/pkg/consts"
	"github.com/Azure/karpenter-provider-azure/pkg/utils"
	"github.com/samber/lo"
)

// openShiftDefaultMaxPods matches the kubelet maxPods that the OpenShift machine config renders.
const openShiftDefaultMaxPods = 250

// openShiftMaxPods resolves the pod capacity Karpenter sizes NodeClaims with. OpenShift nodes
// always run OVN-Kubernetes, so the AKS CNI based pod defaults don't apply to them.
func openShiftMaxPods(nodeClass *v1beta1.AKSNodeClass, provisionMode, networkPlugin, networkPluginMode string) int32 {
	if provisionMode != consts.ProvisionModeOpenShift {
		// OpenShift always runs PROVISION_MODE=openshift; other provision modes only exist for the
		// upstream AKS code paths this fork still carries, so they keep the upstream defaults.
		return utils.GetMaxPods(nodeClass, networkPlugin, networkPluginMode)
	}
	if nodeClass.Spec.MaxPods != nil {
		return lo.FromPtr(nodeClass.Spec.MaxPods)
	}
	return openShiftDefaultMaxPods
}
