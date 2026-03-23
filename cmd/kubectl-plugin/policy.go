package main

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func newPolicyCmd(f cmdutil.Factory, streams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage WorkloadPolicy",
	}

	cmd.SetUsageTemplate(groupUsageTemplate)

	cmd.AddCommand(newPolicyModeProtectCmd(f, streams))
	cmd.AddCommand(newPolicyModeMonitorCmd(f, streams))
	cmd.AddCommand(newPolicyShowCmd(f, streams))
	cmd.AddCommand(newPolicyExecAllowCmd(f, streams))
	cmd.AddCommand(newPolicyExecDenyCmd(f, streams))

	return cmd
}
