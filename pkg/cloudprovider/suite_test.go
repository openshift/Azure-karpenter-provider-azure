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

package cloudprovider

// TODO v1beta1 extra refactor into suite_test.go / cloudprovider_test.go
import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	. "github.com/Azure/karpenter-provider-azure/pkg/test/expectations"
	"github.com/awslabs/operatorpkg/object"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	clock "k8s.io/utils/clock/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	karpv1 "sigs.k8s.io/karpenter/pkg/apis/v1"
	karpv1alpha1 "sigs.k8s.io/karpenter/pkg/apis/v1alpha1"
	corecloudprovider "sigs.k8s.io/karpenter/pkg/cloudprovider"
	"sigs.k8s.io/karpenter/pkg/controllers/nodeoverlay"
	"sigs.k8s.io/karpenter/pkg/controllers/provisioning"
	"sigs.k8s.io/karpenter/pkg/controllers/state"
	"sigs.k8s.io/karpenter/pkg/events"
	coreoptions "sigs.k8s.io/karpenter/pkg/operator/options"
	coretest "sigs.k8s.io/karpenter/pkg/test"
	. "sigs.k8s.io/karpenter/pkg/test/expectations"
	"sigs.k8s.io/karpenter/pkg/test/v1alpha1"
	. "sigs.k8s.io/karpenter/pkg/utils/testing"

	"github.com/Azure/karpenter-provider-azure/pkg/apis"
	"github.com/Azure/karpenter-provider-azure/pkg/apis/v1beta1"
	"github.com/Azure/karpenter-provider-azure/pkg/consts"
	"github.com/Azure/karpenter-provider-azure/pkg/controllers/nodeclass/status"
	"github.com/Azure/karpenter-provider-azure/pkg/fake"
	"github.com/Azure/karpenter-provider-azure/pkg/operator/options"
	"github.com/Azure/karpenter-provider-azure/pkg/providers/instance"
	"github.com/Azure/karpenter-provider-azure/pkg/test"
	"github.com/Azure/karpenter-provider-azure/pkg/utils/zones"
)

var ctx context.Context
var testOptions *options.Options
var stop context.CancelFunc
var env *coretest.Environment
var azureEnv *test.Environment
var azureEnvNonZonal *test.Environment
var coreProvisioner *provisioning.Provisioner
var coreProvisionerNonZonal *provisioning.Provisioner
var cloudProvider *CloudProvider
var cloudProviderNonZonal *CloudProvider
var cluster *state.Cluster
var clusterNonZonal *state.Cluster
var fakeClock *clock.FakeClock
var recorder events.Recorder
var statusController *status.Controller

var nodePool *karpv1.NodePool
var nodeClass *v1beta1.AKSNodeClass
var nodeClaim *karpv1.NodeClaim

var fakeZone1 = zones.MakeAKSLabelZoneFromARMZone(fake.Region, "1")
var defaultTestSKU = fake.MakeSKU("Standard_D2_v3")

func TestCloudProvider(t *testing.T) {
	ctx = TestContextWithLogger(t)
	RegisterFailHandler(Fail)
	RunSpecs(t, "cloudProvider/Azure")
}

var _ = BeforeSuite(func() {
	env = coretest.NewEnvironment(coretest.WithCRDs(apis.CRDs...), coretest.WithCRDs(v1alpha1.CRDs...), coretest.WithFieldIndexers(coretest.NodeProviderIDFieldIndexer(ctx)))
	ctx = coreoptions.ToContext(ctx, coretest.Options())
	ctx, stop = context.WithCancel(ctx)
	fakeClock = clock.NewFakeClock(time.Now())
	recorder = events.NewRecorder(&record.FakeRecorder{})
})

var _ = AfterSuite(func() {
	stop()
	Expect(env.Stop()).To(Succeed(), "Failed to stop environment")
})

var _ = BeforeEach(func() {
	nodeClass = test.AKSNodeClass()
	nodePool = coretest.NodePool(karpv1.NodePool{
		Spec: karpv1.NodePoolSpec{
			Template: karpv1.NodeClaimTemplate{
				Spec: karpv1.NodeClaimTemplateSpec{
					NodeClassRef: &karpv1.NodeClassReference{
						Group: object.GVK(nodeClass).Group,
						Kind:  object.GVK(nodeClass).Kind,
						Name:  nodeClass.Name,
					},
				},
			},
		},
	})

	nodeClaim = coretest.NodeClaim(karpv1.NodeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{karpv1.NodePoolLabelKey: nodePool.Name},
		},
		Spec: karpv1.NodeClaimSpec{
			NodeClassRef: &karpv1.NodeClassReference{
				Group: object.GVK(nodeClass).Group,
				Kind:  object.GVK(nodeClass).Kind,
				Name:  nodeClass.Name,
			},
		},
	})
})

var _ = AfterEach(func() {
	ExpectCleanedUp(ctx, env.Client)
})

// Helper functions for NodeClaim validation
func validateNodeClaimCommon(nodeClaim *karpv1.NodeClaim, nodePool *karpv1.NodePool) {
	// Basic validation
	Expect(nodeClaim).ToNot(BeNil())
	Expect(nodeClaim.Status.Capacity).ToNot(BeEmpty())

	// Status fields
	Expect(nodeClaim.Status.ProviderID).ToNot(BeEmpty())
	Expect(nodeClaim.Status.ImageID).ToNot(BeEmpty())

	// Common labels validation
	Expect(nodeClaim.Labels).To(HaveKey(karpv1.CapacityTypeLabelKey))
	Expect(nodeClaim.Labels).To(HaveKey(karpv1.NodePoolLabelKey))
	if nodePool != nil {
		Expect(nodeClaim.Labels[karpv1.NodePoolLabelKey]).To(Equal(nodePool.Name))
	}
	Expect(nodeClaim.Labels).To(HaveKey(v1.LabelInstanceTypeStable))
	Expect(nodeClaim.Labels).To(HaveKey(v1.LabelArchStable))
	Expect(nodeClaim.Labels).To(HaveKey(v1beta1.LabelSKUName))
	Expect(nodeClaim.Labels).To(HaveKey(v1beta1.LabelSKUFamily))
	Expect(nodeClaim.Labels).To(HaveKey(v1beta1.LabelSKUCPU))
	Expect(nodeClaim.Labels).To(HaveKey(v1beta1.LabelSKUMemory))

	// Zone validation (conditional)
	zone := nodeClaim.Labels[v1.LabelTopologyZone]
	if zone != "" && zone != zones.Regional {
		Expect(zone).To(MatchRegexp(`^[a-z0-9-]+-[0-9]+$`))
	}

	// Capacity and Allocatable resources
	Expect(nodeClaim.Status.Capacity).To(HaveKey(v1.ResourceCPU))
	Expect(nodeClaim.Status.Capacity).To(HaveKey(v1.ResourceMemory))
	Expect(nodeClaim.Status.Allocatable).To(HaveKey(v1.ResourceCPU))
	Expect(nodeClaim.Status.Allocatable).To(HaveKey(v1.ResourceMemory))

	// Lifecycle validation
	Expect(nodeClaim.DeletionTimestamp).To(BeNil())
}

func validateVMNodeClaim(nodeClaim *karpv1.NodeClaim, nodePool *karpv1.NodePool) {
	// Common validations
	validateNodeClaimCommon(nodeClaim, nodePool)

	// VM-specific validation (should NOT have AKS machine annotation)
	Expect(nodeClaim.Annotations).ToNot(HaveKey(v1beta1.AnnotationAKSMachineResourceID))
}

func reconcileCapacityOverlay(customResource v1.ResourceName, overlayCapacity resource.Quantity) {
	GinkgoHelper()
	nodeOverlay := coretest.NodeOverlay(karpv1alpha1.NodeOverlay{
		Spec: karpv1alpha1.NodeOverlaySpec{
			Requirements: []karpv1alpha1.NodeSelectorRequirement{{
				Key:      karpv1.NodePoolLabelKey,
				Operator: v1.NodeSelectorOpIn,
				Values:   []string{nodePool.Name},
			}},
			Capacity: v1.ResourceList{customResource: overlayCapacity},
		},
	})
	ExpectApplied(ctx, env.Client, nodeOverlay)
	nodeOverlayController := nodeoverlay.NewController(env.Clock, env.Client, cloudProvider, azureEnv.InstanceTypeStore, cluster)
	ExpectReconcileSucceeded(ctx, nodeOverlayController, client.ObjectKeyFromObject(nodeOverlay))
}

type nodeOverlayCapacityTestOptions struct {
	validateNodeClaim func(*karpv1.NodeClaim)
	resetCreateCalls  func()
	expectCreateCalls func()
}

func runNodeOverlayCapacityTests(testOptions nodeOverlayCapacityTestOptions) {
	Context("NodeOverlay", func() {
		It("should launch a NodeClaim that requests capacity added by a NodeOverlay", func() {
			ctx = coreoptions.ToContext(ctx, coretest.Options(coretest.OptionsFields{
				FeatureGates: coretest.FeatureGates{NodeOverlay: lo.ToPtr(true)},
			}))
			customResource := v1.ResourceName("example.com/dongle")
			overlayCapacity := resource.MustParse("100")
			nodeClaim.Spec.Resources.Requests = v1.ResourceList{customResource: resource.MustParse("1")}

			ExpectApplied(ctx, env.Client, nodeClass, nodePool, nodeClaim)
			reconcileCapacityOverlay(customResource, overlayCapacity)

			if testOptions.resetCreateCalls != nil {
				testOptions.resetCreateCalls()
			}
			cloudProviderMachine, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(err).ToNot(HaveOccurred())
			Expect(cloudProviderMachine).ToNot(BeNil())
			if testOptions.validateNodeClaim != nil {
				testOptions.validateNodeClaim(cloudProviderMachine)
			}
			if testOptions.expectCreateCalls != nil {
				testOptions.expectCreateCalls()
			}
			capacity, ok := cloudProviderMachine.Status.Capacity[customResource]
			Expect(ok).To(BeTrue())
			Expect(capacity.Cmp(overlayCapacity)).To(Equal(0))
			allocatable, ok := cloudProviderMachine.Status.Allocatable[customResource]
			Expect(ok).To(BeTrue())
			Expect(allocatable.Cmp(overlayCapacity)).To(Equal(0))
		})

		It("should not use overlaid capacity when NodeOverlay is disabled", func() {
			// Explicitly disable the NodeOverlay feature gate so this test does not
			// depend on ordering with the previous It block that enables it.
			ctx = coreoptions.ToContext(ctx, coretest.Options(coretest.OptionsFields{
				FeatureGates: coretest.FeatureGates{NodeOverlay: lo.ToPtr(false)},
			}))
			customResource := v1.ResourceName("example.com/dongle")
			overlayCapacity := resource.MustParse("100")
			nodeClaim.Spec.Resources.Requests = v1.ResourceList{customResource: resource.MustParse("1")}

			ExpectApplied(ctx, env.Client, nodeClass, nodePool, nodeClaim)
			reconcileCapacityOverlay(customResource, overlayCapacity)

			if testOptions.resetCreateCalls != nil {
				testOptions.resetCreateCalls()
			}
			cloudProviderMachine, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(corecloudprovider.IsInsufficientCapacityError(err)).To(BeTrue())
			Expect(cloudProviderMachine).To(BeNil())
		})
	})
}

func vmNodeOverlayCapacityTestOptions() nodeOverlayCapacityTestOptions {
	return nodeOverlayCapacityTestOptions{
		validateNodeClaim: func(nodeClaim *karpv1.NodeClaim) {
			validateVMNodeClaim(nodeClaim, nodePool)
		},
		resetCreateCalls: func() {
			azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Reset()
			azureEnv.AKSMachinesAPI.AKSMachineCreateOrUpdateBehavior.CalledWithInput.Reset()
		},
		expectCreateCalls: func() {
			Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(1))
			Expect(azureEnv.AKSMachinesAPI.AKSMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(0))
		},
	}
}

var _ = Describe("CloudProvider", func() {
	// Attention: tests under "ProvisionMode = AKSScriptless" are not applicable to ProvisionMode = AKSMachineAPI option.
	// Due to different assumptions, not all tests can be shared. Add tests for AKS machine instances in a different Context/file.
	// If ProvisionMode = AKSScriptless is no longer supported, their code/tests will be replaced with ProvisionMode = AKSMachineAPI.
	Context("ProvisionMode = AKSScriptless, ManageExistingAKSMachines = false", func() {
		BeforeEach(func() {
			testOptions = test.Options(test.OptionsFields{
				ProvisionMode:             lo.ToPtr(consts.ProvisionModeAKSScriptless),
				ManageExistingAKSMachines: lo.ToPtr(false),
			})
			ctx = coreoptions.ToContext(ctx, coretest.Options())
			ctx = options.ToContext(ctx, testOptions)

			azureEnv = test.NewEnvironment(ctx, env)
			test.ApplyDefaultStatus(nodeClass, env, testOptions.UseSIG)
			cloudProvider = New(azureEnv.InstanceTypesProvider, azureEnv.VMInstanceProvider, azureEnv.AKSMachineProvider, recorder, env.Client, azureEnv.ImageProvider, azureEnv.InstanceTypeStore)

			cluster = state.NewCluster(fakeClock, env.Client, cloudProvider)
			coreProvisioner = provisioning.NewProvisioner(env.Client, recorder, cloudProvider, cluster, fakeClock)
		})

		AfterEach(func() {
			// Wait for any async polling goroutines to complete before resetting
			cloudProvider.WaitForInstancePromises()
			cluster.Reset()
			azureEnv.Reset(ctx)
		})

		It("should list nodeclaim created by the CloudProvider", func() {
			ExpectApplied(ctx, env.Client, nodeClass, nodePool)
			pod := coretest.UnschedulablePod()
			ExpectProvisionedAndWaitForPromises(ctx, env.Client, cluster, cloudProvider, coreProvisioner, azureEnv, pod)
			ExpectScheduled(ctx, env.Client, pod)
			Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(1))

			nodeClaims, _ := cloudProvider.List(ctx)
			Expect(azureEnv.AKSMachinesAPI.AKSMachineNewListPagerBehavior.CalledWithInput.Len()).To(Equal(0))
			Expect(azureEnv.AzureResourceGraphAPI.AzureResourceGraphResourcesBehavior.CalledWithInput.Len()).To(Equal(1))
			queryRequest := azureEnv.AzureResourceGraphAPI.AzureResourceGraphResourcesBehavior.CalledWithInput.Pop().Query
			Expect(*queryRequest.Query).To(Equal(instance.GetVMListQueryBuilder(azureEnv.AzureResourceGraphAPI.ResourceGroup).String()))
			Expect(nodeClaims).To(HaveLen(1))
			validateVMNodeClaim(nodeClaims[0], nodePool)
			resp, _ := azureEnv.VirtualMachinesAPI.Get(ctx, azureEnv.AzureResourceGraphAPI.ResourceGroup, nodeClaims[0].Name, nil)
			Expect(resp.VirtualMachine).ToNot(BeNil())
		})
		It("should list nodeclaim with correct instance type even after capacity error marks offerings unavailable", func() {
			ExpectApplied(ctx, env.Client, nodeClass, nodePool)
			pod := coretest.UnschedulablePod()
			ExpectProvisionedAndWaitForPromises(ctx, env.Client, cluster, cloudProvider, coreProvisioner, azureEnv, pod)
			ExpectScheduled(ctx, env.Client, pod)
			Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(1))

			// Get the instance type from the created VM
			vm := azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Pop().VM
			vmSize := string(lo.FromPtr(vm.Properties.HardwareProfile.VMSize))
			Expect(vmSize).ToNot(BeEmpty())

			// Simulate a capacity error by marking all offerings for this instance type as unavailable
			for _, zone := range azureEnv.Zones() {
				azureEnv.UnavailableOfferingsCache.MarkUnavailable(ctx, "ZonalAllocationFailure", fake.MakeSKU(vmSize), zone, karpv1.CapacityTypeOnDemand)
				azureEnv.UnavailableOfferingsCache.MarkUnavailable(ctx, "ZonalAllocationFailure", fake.MakeSKU(vmSize), zone, karpv1.CapacityTypeSpot)
			}

			// List should still return the nodeclaim with the correct instance type
			nodeClaims, err := cloudProvider.List(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(nodeClaims).To(HaveLen(1))
			validateVMNodeClaim(nodeClaims[0], nodePool)
			Expect(nodeClaims[0].Labels[v1.LabelInstanceTypeStable]).To(Equal(vmSize))
		})
		It("should return an ICE error when there are no instance types to launch", func() {
			// Specify no instance types and expect to receive a capacity error
			nodeClaim.Spec.Requirements = []karpv1.NodeSelectorRequirementWithMinValues{
				{
					Key:      v1.LabelInstanceTypeStable,
					Operator: v1.NodeSelectorOpIn,
					Values:   []string{"doesnotexist"}, // will not match any instance types
				},
			}

			ExpectApplied(ctx, env.Client, nodePool, nodeClass, nodeClaim)
			cloudProviderMachine, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(corecloudprovider.IsInsufficientCapacityError(err)).To(BeTrue())
			Expect(cloudProviderMachine).To(BeNil())
		})

		runNodeOverlayCapacityTests(vmNodeOverlayCapacityTestOptions())

		Context("AKS Machine API integration", func() {
			It("should not call writes to AKS Machine API", func() {
				ExpectApplied(ctx, env.Client, nodeClass, nodePool)
				pod := coretest.UnschedulablePod()
				ExpectProvisionedAndWaitForPromises(ctx, env.Client, cluster, cloudProvider, coreProvisioner, azureEnv, pod)
				ExpectScheduled(ctx, env.Client, pod)

				Expect(azureEnv.AKSMachinesAPI.AKSMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(0))
			})

			Context("AKS Machines Pool Management", func() {
				It("should handle AKS machines pool not found on each CloudProvider operation", func() {
					// First create a successful VM
					ExpectApplied(ctx, env.Client, nodeClass, nodePool)
					pod := coretest.UnschedulablePod()
					ExpectProvisionedAndWaitForPromises(ctx, env.Client, cluster, cloudProvider, coreProvisioner, azureEnv, pod)
					ExpectScheduled(ctx, env.Client, pod)

					// cloudprovider.List should return vm nodeclaim
					nodeClaims, err := cloudProvider.List(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(nodeClaims).To(HaveLen(1))
					validateVMNodeClaim(nodeClaims[0], nodePool)

					// cloudprovider.Delete should be fine also
					err = cloudProvider.Delete(ctx, nodeClaims[0])
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})

	Context("MarketplaceImages", func() {
		BeforeEach(func() {
			testOptions = test.Options(test.OptionsFields{
				ProvisionMode:             lo.ToPtr(consts.ProvisionModeOpenShift),
				ManageExistingAKSMachines: lo.ToPtr(false),
			})
			ctx = coreoptions.ToContext(ctx, coretest.Options())
			ctx = options.ToContext(ctx, testOptions)
			azureEnv = test.NewEnvironment(ctx, env)
			test.ApplyDefaultStatus(nodeClass, env, testOptions.UseSIG)
			nodeClass.Spec.MarketplaceImages = []string{"publisher:offer:amd64:1.0.0", "publisher:offer:arm64:1.0.0"}
			nodeClass.Spec.UserData = lo.ToPtr(`{"ignition":{"version":"3.2.0"}}`)
			nodeClass.Status.Images = []v1beta1.NodeImage{
				{ID: "publisher:offer:amd64:1.0.0", Requirements: []v1.NodeSelectorRequirement{
					{Key: v1.LabelArchStable, Operator: v1.NodeSelectorOpIn, Values: []string{"amd64"}},
					{Key: v1beta1.LabelSKUHyperVGeneration, Operator: v1.NodeSelectorOpIn, Values: []string{"2"}},
				}},
				{ID: "publisher:offer:arm64:1.0.0", Requirements: []v1.NodeSelectorRequirement{
					{Key: v1.LabelArchStable, Operator: v1.NodeSelectorOpIn, Values: []string{"arm64"}},
					{Key: v1beta1.LabelSKUHyperVGeneration, Operator: v1.NodeSelectorOpIn, Values: []string{"2"}},
				}},
			}
			nodeClass.StatusConditions().SetTrue(v1beta1.ConditionTypeImagesReady)
			cloudProvider = New(azureEnv.InstanceTypesProvider, azureEnv.VMInstanceProvider, azureEnv.AKSMachineProvider, recorder, env.Client, azureEnv.ImageProvider, azureEnv.InstanceTypeStore)
			cluster = state.NewCluster(fakeClock, env.Client, cloudProvider)
			coreProvisioner = provisioning.NewProvisioner(env.Client, recorder, cloudProvider, cluster, fakeClock)
		})
		AfterEach(func() {
			cloudProvider.WaitForInstancePromises()
			cluster.Reset()
			azureEnv.Reset(ctx)
		})
		It("excludes Gen1-only sizes from scheduling with a Gen2 marketplace image", func() {
			nodeClass.Status.Images = nodeClass.Status.Images[:1]
			ExpectApplied(ctx, env.Client, nodeClass, nodePool)
			types, err := cloudProvider.GetInstanceTypes(ctx, nodePool)
			Expect(err).NotTo(HaveOccurred(), "listing instance types for a Gen2 marketplace image")
			names := lo.Map(types, func(instanceType *corecloudprovider.InstanceType, _ int) string { return instanceType.Name })
			Expect(names).NotTo(ContainElement("Standard_D2_v2"), "Gen1-only size must be excluded")
			Expect(names).To(ContainElements("Standard_D2as_v6", "Standard_D2s_v3"), "Gen2-only and dual-generation sizes remain eligible")
		})
		DescribeTable("checks marketplace image generation before creating a VM", func(size string, compatible bool) {
			nodeClaim.Spec.Requirements = []karpv1.NodeSelectorRequirementWithMinValues{{Key: v1.LabelInstanceTypeStable, Operator: v1.NodeSelectorOpIn, Values: []string{size}}}
			ExpectApplied(ctx, env.Client, nodeClass, nodePool, nodeClaim)
			created, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			if !compatible {
				Expect(corecloudprovider.IsInsufficientCapacityError(err)).To(BeTrue(), "Gen1-only size cannot launch a Gen2 image")
				Expect(created).To(BeNil())
				Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(BeZero())
				return
			}
			Expect(err).NotTo(HaveOccurred(), "creating a VM with a compatible marketplace image")
			Expect(created.Labels[v1.LabelInstanceTypeStable]).To(Equal(size))
			Expect(created.Status.ImageID).To(Equal("publisher:offer:amd64:1.0.0"))
			Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(1))
		},
			Entry("rejects a Gen1-only size", "Standard_D2_v2", false),
			Entry("accepts a Gen2-only size", "Standard_D2as_v6", true),
			Entry("accepts a dual-generation size", "Standard_D2s_v3", true),
		)
		It("filters scheduling and creation to ARM64 when only ARM64 has an image", func() {
			nodeClass.Status.Images = nodeClass.Status.Images[1:]
			ExpectApplied(ctx, env.Client, nodeClass, nodePool)
			types, err := cloudProvider.GetInstanceTypes(ctx, nodePool)
			Expect(err).NotTo(HaveOccurred())
			Expect(types).NotTo(BeEmpty())
			for _, instanceType := range types {
				Expect(instanceType.Requirements.Get(v1.LabelArchStable).Values()).To(Equal([]string{"arm64"}))
			}
			// The normal scheduler must not choose the smaller, cheaper AMD64 candidate.
			pod := coretest.UnschedulablePod()
			ExpectProvisionedAndWaitForPromises(ctx, env.Client, cluster, cloudProvider, coreProvisioner, azureEnv, pod)
			ExpectScheduled(ctx, env.Client, pod)
			Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(1))
			vm := azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Pop().VM
			Expect(vm.Properties.StorageProfile.ImageReference.SKU).To(HaveValue(Equal("arm64")))
		})
		DescribeTable("launches the image matching the selected architecture", func(arch, sku, version string, reversed bool) {
			if reversed {
				nodeClass.Status.Images[0], nodeClass.Status.Images[1] = nodeClass.Status.Images[1], nodeClass.Status.Images[0]
			}
			nodeClaim.Spec.Requirements = []karpv1.NodeSelectorRequirementWithMinValues{{Key: v1.LabelArchStable, Operator: v1.NodeSelectorOpIn, Values: []string{arch}}}
			ExpectApplied(ctx, env.Client, nodeClass, nodePool, nodeClaim)
			created, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(err).NotTo(HaveOccurred())
			Expect(created.Labels[v1.LabelArchStable]).To(Equal(arch))
			Expect(created.Status.ImageID).To(Equal("publisher:offer:" + sku + ":" + version))
			vm := azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Pop().VM
			Expect(vm.Properties.StorageProfile.ImageReference).To(Equal(&armcompute.ImageReference{
				Publisher: lo.ToPtr("publisher"), Offer: lo.ToPtr("offer"), SKU: lo.ToPtr(sku), Version: lo.ToPtr(version),
			}))
		},
			Entry("AMD64", "amd64", "amd64", "1.0.0", false),
			Entry("ARM64", "arm64", "arm64", "1.0.0", false),
			Entry("AMD64 reversed", "amd64", "amd64", "1.0.0", true),
			Entry("ARM64 reversed", "arm64", "arm64", "1.0.0", true),
		)
		It("reads the selected URN from status rather than same-generation spec", func() {
			nodeClass.Status.Images[0].ID = "resolved:status:amd64:2.0.0"
			nodeClaim.Spec.Requirements = []karpv1.NodeSelectorRequirementWithMinValues{{Key: v1.LabelArchStable, Operator: v1.NodeSelectorOpIn, Values: []string{"amd64"}}}
			ExpectApplied(ctx, env.Client, nodeClass, nodePool, nodeClaim)
			created, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(err).NotTo(HaveOccurred())
			Expect(created.Status.ImageID).To(Equal("resolved:status:amd64:2.0.0"))
			vm := azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Pop().VM
			Expect(vm.Properties.StorageProfile.ImageReference.Publisher).To(HaveValue(Equal("resolved")))
			Expect(vm.Properties.StorageProfile.ImageReference.Offer).To(HaveValue(Equal("status")))
			Expect(vm.Properties.StorageProfile.ImageReference.Version).To(HaveValue(Equal("2.0.0")))
		})
		It("rejects creation when no compatible candidate has an image", func() {
			nodeClass.Status.Images = nodeClass.Status.Images[1:]
			nodeClaim.Spec.Requirements = []karpv1.NodeSelectorRequirementWithMinValues{{Key: v1.LabelArchStable, Operator: v1.NodeSelectorOpIn, Values: []string{"amd64"}}}
			ExpectApplied(ctx, env.Client, nodeClass, nodePool, nodeClaim)
			created, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(corecloudprovider.IsInsufficientCapacityError(err)).To(BeTrue())
			Expect(created).To(BeNil())
			Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(BeZero())
		})
		DescribeTable("rejects invalid image status without creating a VM", func(mutate func(*v1beta1.AKSNodeClass)) {
			mutate(nodeClass)
			ExpectApplied(ctx, env.Client, nodeClass, nodePool, nodeClaim)
			created, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(err).To(HaveOccurred())
			Expect(created).To(BeNil())
			Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(BeZero())
		},
			Entry("empty", func(n *v1beta1.AKSNodeClass) { n.Status.Images = nil }),
			Entry("unready", func(n *v1beta1.AKSNodeClass) {
				n.StatusConditions().SetFalse(v1beta1.ConditionTypeImagesReady, "Failed", "lookup failed")
			}),
			Entry("unconstrained", func(n *v1beta1.AKSNodeClass) {
				n.Status.Images = []v1beta1.NodeImage{{ID: "publisher:offer:amd64:1.0.0", Requirements: []v1.NodeSelectorRequirement{}}}
			}),
		)
		It("reconstructs an existing VM after its architecture loses image coverage", func() {
			nodeClaim.Spec.Requirements = []karpv1.NodeSelectorRequirementWithMinValues{{Key: v1.LabelArchStable, Operator: v1.NodeSelectorOpIn, Values: []string{"amd64"}}}
			ExpectApplied(ctx, env.Client, nodeClass, nodePool, nodeClaim)
			created, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(err).NotTo(HaveOccurred())
			nodeClass = ExpectExists(ctx, env.Client, nodeClass)
			nodeClass.Status.Images = nodeClass.Status.Images[1:]
			ExpectApplied(ctx, env.Client, nodeClass)
			reconstructed, err := cloudProvider.Get(ctx, created.Status.ProviderID)
			Expect(err).NotTo(HaveOccurred())
			Expect(reconstructed.Labels[v1.LabelInstanceTypeStable]).To(Equal(created.Labels[v1.LabelInstanceTypeStable]))
			Expect(reconstructed.Labels[v1.LabelArchStable]).To(Equal("amd64"))
			Expect(reconstructed.Status.Capacity).To(Equal(created.Status.Capacity))
			listed, err := cloudProvider.List(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(listed).To(HaveLen(1))
			Expect(listed[0].Labels[v1.LabelInstanceTypeStable]).To(Equal(created.Labels[v1.LabelInstanceTypeStable]))
		})
		runNodeOverlayCapacityTests(vmNodeOverlayCapacityTestOptions())
	})

	Context("ProvisionMode = AKSScriptless, ManageExistingAKSMachines = true", func() {
		BeforeEach(func() {
			testOptions = test.Options(test.OptionsFields{
				ProvisionMode:             lo.ToPtr(consts.ProvisionModeAKSScriptless),
				ManageExistingAKSMachines: lo.ToPtr(true),
			})
			ctx = coreoptions.ToContext(ctx, coretest.Options())
			ctx = options.ToContext(ctx, testOptions)

			azureEnv = test.NewEnvironment(ctx, env)
			test.ApplyDefaultStatus(nodeClass, env, testOptions.UseSIG)
			cloudProvider = New(azureEnv.InstanceTypesProvider, azureEnv.VMInstanceProvider, azureEnv.AKSMachineProvider, recorder, env.Client, azureEnv.ImageProvider, azureEnv.InstanceTypeStore)

			cluster = state.NewCluster(fakeClock, env.Client, cloudProvider)
			coreProvisioner = provisioning.NewProvisioner(env.Client, recorder, cloudProvider, cluster, fakeClock)
		})

		AfterEach(func() {
			// Wait for any async polling goroutines to complete before resetting
			cloudProvider.WaitForInstancePromises()
			cluster.Reset()
			azureEnv.Reset(ctx)
		})

		It("should list nodeclaim created by the CloudProvider", func() {
			ExpectApplied(ctx, env.Client, nodeClass, nodePool)
			pod := coretest.UnschedulablePod()
			ExpectProvisionedAndWaitForPromises(ctx, env.Client, cluster, cloudProvider, coreProvisioner, azureEnv, pod)
			ExpectScheduled(ctx, env.Client, pod)
			Expect(azureEnv.VirtualMachinesAPI.VirtualMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(1))

			nodeClaims, _ := cloudProvider.List(ctx)
			Expect(azureEnv.AKSMachinesAPI.AKSMachineNewListPagerBehavior.CalledWithInput.Len()).To(Equal(1)) // Expect to be called in case of existing AKS machines
			Expect(azureEnv.AzureResourceGraphAPI.AzureResourceGraphResourcesBehavior.CalledWithInput.Len()).To(Equal(1))
			queryRequest := azureEnv.AzureResourceGraphAPI.AzureResourceGraphResourcesBehavior.CalledWithInput.Pop().Query
			Expect(*queryRequest.Query).To(Equal(instance.GetVMListQueryBuilder(azureEnv.AzureResourceGraphAPI.ResourceGroup).String()))
			Expect(nodeClaims).To(HaveLen(1))
			validateVMNodeClaim(nodeClaims[0], nodePool)
			resp, _ := azureEnv.VirtualMachinesAPI.Get(ctx, azureEnv.AzureResourceGraphAPI.ResourceGroup, nodeClaims[0].Name, nil)
			Expect(resp.VirtualMachine).ToNot(BeNil())
		})
		It("should return an ICE error when there are no instance types to launch", func() {
			// Specify no instance types and expect to receive a capacity error
			nodeClaim.Spec.Requirements = []karpv1.NodeSelectorRequirementWithMinValues{
				{
					Key:      v1.LabelInstanceTypeStable,
					Operator: v1.NodeSelectorOpIn,
					Values:   []string{"doesnotexist"}, // will not match any instance types
				},
			}

			ExpectApplied(ctx, env.Client, nodePool, nodeClass, nodeClaim)
			cloudProviderMachine, err := CreateAndWaitForPromises(ctx, cloudProvider, azureEnv, nodeClaim)
			Expect(corecloudprovider.IsInsufficientCapacityError(err)).To(BeTrue())
			Expect(cloudProviderMachine).To(BeNil())
		})

		runNodeOverlayCapacityTests(vmNodeOverlayCapacityTestOptions())

		Context("AKS Machine API integration", func() {
			It("should not call writes to AKS Machine API", func() {
				ExpectApplied(ctx, env.Client, nodeClass, nodePool)
				pod := coretest.UnschedulablePod()
				ExpectProvisionedAndWaitForPromises(ctx, env.Client, cluster, cloudProvider, coreProvisioner, azureEnv, pod)
				ExpectScheduled(ctx, env.Client, pod)

				Expect(azureEnv.AKSMachinesAPI.AKSMachineCreateOrUpdateBehavior.CalledWithInput.Len()).To(Equal(0))
			})

			Context("AKS Machines Pool Management", func() {
				It("should handle AKS machines pool not found on each CloudProvider operation", func() {
					// First create a successful VM
					ExpectApplied(ctx, env.Client, nodeClass, nodePool)
					pod := coretest.UnschedulablePod()
					ExpectProvisionedAndWaitForPromises(ctx, env.Client, cluster, cloudProvider, coreProvisioner, azureEnv, pod)
					ExpectScheduled(ctx, env.Client, pod)

					// cloudprovider.List should return vm nodeclaim
					nodeClaims, err := cloudProvider.List(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(nodeClaims).To(HaveLen(1))
					validateVMNodeClaim(nodeClaims[0], nodePool)

					// cloudprovider.Delete should be fine also
					err = cloudProvider.Delete(ctx, nodeClaims[0])
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
