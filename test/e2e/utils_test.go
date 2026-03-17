package e2e_test

import (
	"context"

	"github.com/rancher-sandbox/runtime-enforcer/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	apimachinerywait "k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
)

func getDeploymentPolicyMutateOption(
	namespace string,
	policy string, //nolint:unparam // we want to keep the flexibility to support different policy name.
) decoder.DecodeOption {
	// Support only deployment right now.
	return decoder.MutateOption(func(obj k8s.Object) error {
		deployment := obj.(*appsv1.Deployment)
		deployment.SetNamespace(namespace)
		deployment.Spec.Template.Labels[v1alpha1.PolicyLabelKey] = policy
		return nil
	})
}

func DaemonSetUpToDate(r *resources.Resources, daemonset *appsv1.DaemonSet) apimachinerywait.ConditionWithContextFunc {
	return func(ctx context.Context) (bool, error) {
		if err := r.Get(ctx, daemonset.GetName(), daemonset.GetNamespace(), daemonset); err != nil {
			return false, err
		}
		status := daemonset.Status
		if status.UpdatedNumberScheduled != status.DesiredNumberScheduled {
			return false, nil
		}
		return true, nil
	}
}

func DeploymentUpToDate(
	r *resources.Resources,
	deployment *appsv1.Deployment,
) apimachinerywait.ConditionWithContextFunc {
	return func(ctx context.Context) (bool, error) {
		if err := r.Get(ctx, deployment.GetName(), deployment.GetNamespace(), deployment); err != nil {
			return false, err
		}
		status := deployment.Status
		if status.Replicas != *deployment.Spec.Replicas {
			return false, nil
		}
		if status.UpdatedReplicas != status.Replicas {
			return false, nil
		}
		return true, nil
	}
}
