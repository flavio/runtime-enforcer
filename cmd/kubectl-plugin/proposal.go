package main

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func newProposalCmd(f cmdutil.Factory, streams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal",
		Short: "Manage WorkloadPolicyProposal",
	}

	cmd.SetUsageTemplate(groupUsageTemplate)

	cmd.AddCommand(newProposalPromoteCmd(f, streams))

	return cmd
}
