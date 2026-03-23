package main

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	utilcomp "k8s.io/kubectl/pkg/util/completion"
)

var version = "dev"

// Custom usage template: no "kubectl [command]" line.
const (
	rootUsageTemplate = `Usage:
  {{.UseLine}}

Available Commands:
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}  {{rpad .Name .NamePadding}} {{.Short}}
{{end}}{{end}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
`
)

func registerCompletionFuncForGlobalFlags(cmd *cobra.Command, f cmdutil.Factory) {
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"namespace",
		func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return utilcomp.CompGetResource(f, "namespace", toComplete), cobra.ShellCompDirectiveNoFileComp
		}))
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"context",
		func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return utilcomp.ListContextsInConfig(toComplete), cobra.ShellCompDirectiveNoFileComp
		}))
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"cluster",
		func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return utilcomp.ListClustersInConfig(toComplete), cobra.ShellCompDirectiveNoFileComp
		}))
	cmdutil.CheckErr(cmd.RegisterFlagCompletionFunc(
		"user",
		func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return utilcomp.ListUsersInConfig(toComplete), cobra.ShellCompDirectiveNoFileComp
		}))
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "runtime-enforcer",
		Long:    "Kubernetes plugin for SUSE Security Runtime Enforcer",
		Version: version,
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.SetUsageTemplate(rootUsageTemplate)

	// Create a shared iostream.
	streams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	// Add flags to persistent flags so they are inherited by all subcommands
	configFlags := genericclioptions.NewConfigFlags(true).WithWarningPrinter(streams)
	configFlags.AddFlags(cmd.PersistentFlags())

	// Create cmdutil.Factory for use in completion functions
	f := cmdutil.NewFactory(configFlags)
	utilcomp.SetFactoryForCompletion(f)

	// Register completion functions, so we can auto-complete global flags like --namespace, --context, etc.
	registerCompletionFuncForGlobalFlags(cmd, f)

	cmd.AddCommand(newProposalCmd(f, streams))
	cmd.AddCommand(newPolicyCmd(f, streams))

	return cmd
}

func main() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
