package workloadpolicyhandler_test

import (
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rancher-sandbox/runtime-enforcer/api/v1alpha1"
	"github.com/rancher-sandbox/runtime-enforcer/internal/workloadpolicyhandler"
	agentv1 "github.com/rancher-sandbox/runtime-enforcer/proto/agent/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("WorkloadPolicy Controller", func() {
	Context("When reconciling a resource", func() {
		const policyName = "test-policy"
		const testNamespace = "default"

		policy := &v1alpha1.WorkloadPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      policyName,
				Namespace: testNamespace,
			},
			Spec: v1alpha1.WorkloadPolicySpec{
				Mode: "monitor",
				RulesByContainer: map[string]*v1alpha1.WorkloadPolicyRules{
					"main": {
						Executables: v1alpha1.WorkloadPolicyExecutables{
							Allowed: []string{"/usr/bin/sleep"},
						},
					},
				},
			},
		}
		typeNamespacedName := types.NamespacedName{
			Name:      policyName,
			Namespace: testNamespace,
		}

		AfterEach(func() {
			// delete the workload policy to prevent resource leak
			resource := &v1alpha1.WorkloadPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      typeNamespacedName.Name,
					Namespace: typeNamespacedName.Namespace,
				},
			}

			_ = k8sClient.Delete(ctx, resource)
			Expect(k8sClient.Get(ctx, typeNamespacedName, resource)).Should(HaveOccurred())
		})

		It("Should reconcile a WorkloadPolicy correctly", func() {
			By("Creating a new WorkloadPolicy")
			Expect(k8sClient.Create(ctx, policy)).To(Succeed())

			By("Maintaining WorkloadPolicy in cache")
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			resolver := newMockResolver(GinkgoTB())

			By("Reconciling the created resource")
			controllerReconciler := workloadpolicyhandler.NewWorkloadPolicyHandler(
				k8sClient,
				logger,
				resolver,
			)

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Updating internal status correctly")
			policyStatus := resolver.GetPolicyStatuses()
			Expect(policyStatus).NotTo(BeNil())

			status, exists := policyStatus[policy.NamespacedName()]
			Expect(exists).To(BeTrue())
			Expect(status.State).To(Equal(agentv1.PolicyState_POLICY_STATE_READY))
			Expect(status.Mode).To(Equal(agentv1.PolicyMode_POLICY_MODE_MONITOR))

			By("Cleaning up the specific WorkloadPolicy resource")

			// Delete the workload policy
			resource := &v1alpha1.WorkloadPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      typeNamespacedName.Name,
					Namespace: typeNamespacedName.Namespace,
				},
			}

			// Ensure it's deleted correctly.
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			Expect(k8sClient.Get(ctx, typeNamespacedName, resource)).Should(HaveOccurred())

			// Reconcile again to trigger the deletion handling logic.
			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Updating internal status correctly")
			policyStatus = resolver.GetPolicyStatuses()
			Expect(policyStatus).NotTo(BeNil())

			status, exists = policyStatus[policy.NamespacedName()]
			Expect(exists).To(BeFalse())
		})
	})
})
