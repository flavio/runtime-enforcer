package main

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	apiv1alpha1 "github.com/rancher-sandbox/runtime-enforcer/api/v1alpha1"
	"github.com/rancher-sandbox/runtime-enforcer/internal/types/policymode"
	fakeclient "github.com/rancher-sandbox/runtime-enforcer/pkg/generated/clientset/versioned/fake"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRunPolicyModeSet(t *testing.T) {
	t.Parallel()

	const (
		namespace  = "test"
		policyName = "test-policy"
	)

	createTestPolicyWithMode := func(mode string) *apiv1alpha1.WorkloadPolicy {
		return &apiv1alpha1.WorkloadPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      policyName,
				Namespace: namespace,
			},
			Spec: apiv1alpha1.WorkloadPolicySpec{
				Mode: mode,
			},
		}
	}

	tests := []struct {
		name           string
		policy         *apiv1alpha1.WorkloadPolicy
		expectedMode   string
		expectedOutput string
		expectedError  string
	}{
		{
			name:         "monitor to protect",
			policy:       createTestPolicyWithMode(policymode.MonitorString),
			expectedMode: policymode.ProtectString,
		},
		{
			name:         "protect to monitor",
			policy:       createTestPolicyWithMode(policymode.ProtectString),
			expectedMode: policymode.MonitorString,
		},
		{
			name:           "already in target mode",
			policy:         createTestPolicyWithMode(policymode.MonitorString),
			expectedMode:   policymode.MonitorString,
			expectedOutput: fmt.Sprintf("is already in %q mode.", policymode.MonitorString),
		},
		{
			name:          "missing policy",
			policy:        &apiv1alpha1.WorkloadPolicy{},
			expectedMode:  policymode.MonitorString,
			expectedError: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clientset := fakeclient.NewClientset(tt.policy)
			securityClient := clientset.SecurityV1alpha1()

			var out bytes.Buffer
			opts := &policyModeOptions{
				commonOptions: commonOptions{
					Namespace: namespace,
					DryRun:    false,
				},
				PolicyName: policyName,
				Mode:       tt.expectedMode,
			}

			ctx, cancel := context.WithTimeout(context.Background(), defaultOperationTimeout)
			defer cancel()

			err := runPolicyModeSet(ctx, securityClient, opts, &out)
			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			updatedPolicy, getErr := securityClient.WorkloadPolicies(namespace).
				Get(ctx, policyName, metav1.GetOptions{})
			require.NoError(t, getErr)
			require.Equal(t, tt.expectedMode, updatedPolicy.Spec.Mode)
		})
	}
}
