package main

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func newPolicyShowCmd(f cmdutil.Factory, streams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show WorkloadPolicy information",
	}

	cmd.SetUsageTemplate(groupUsageTemplate)

	cmd.AddCommand(newPolicyShowProtectionCmd(f, streams))

	return cmd
}
