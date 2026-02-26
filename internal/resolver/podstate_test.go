//nolint:testpackage // testing internal methods
package resolver

import (
	"testing"

	"github.com/rancher-sandbox/runtime-enforcer/api/v1alpha1"
	"github.com/stretchr/testify/require"
)

func TestPodState(t *testing.T) {
	// Test the podInfo struct and its methods here
	namespace := "test-namespace"
	name := "test-name"
	policyName := "test-policy"
	labels := Labels{
		"test-label":            "test-value",
		v1alpha1.PolicyLabelKey: policyName,
	}

	podState := &podState{
		info: &podInfo{
			namespace: namespace,
			name:      name,
			labels:    labels,
		},
		containers: map[ContainerID]*containerInfo{
			"test-container": {
				cgID: CgroupID(12),
				name: "test-container-name",
			},
		},
	}

	require.Equal(t, name, podState.podName())
	require.Equal(t, namespace, podState.podNamespace())
	require.True(t, podState.matchPolicy(policyName, namespace))
	// same name but another namespace.
	require.False(t, podState.matchPolicy(policyName, "random-namespace"))
	// same namespace but different name.
	require.False(t, podState.matchPolicy("random-name", namespace))
}
