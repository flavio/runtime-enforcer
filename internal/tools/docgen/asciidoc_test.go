package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/kubewarden/runtime-enforcer/internal/kubectlplugin"
)

// newFixtureRootCmd builds a small, fixed command tree. Tests use this
// tree to check the generated AsciiDoc output. It does not depend on
// the real kubectl plugin command tree.
func newFixtureRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "root",
		Short: "Root command",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	root.PersistentFlags().Bool("verbose", false, "Increase verbosity")

	sub := &cobra.Command{
		Use:     "sub [flags] <arg>",
		Short:   "Sub command",
		Long:    "Sub command long description.",
		Example: "root sub --flag value",
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
	sub.Flags().String("flag", "", "A local flag")
	root.AddCommand(sub)

	nested := &cobra.Command{
		Use:   "nested",
		Short: "Nested command",
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
	sub.AddCommand(nested)

	return root
}

func TestGenAsciiDoc_StartsWithDoctitleAndBlankLine(t *testing.T) {
	var buf bytes.Buffer

	require.NoError(t, genASCIIDoc(newFixtureRootCmd(), "", &buf))

	lines := bytes.SplitN(buf.Bytes(), []byte("\n"), 3)
	require.GreaterOrEqual(t, len(lines), 3)
	require.Equal(t, "= Command-Line Help for `root`", string(lines[0]))
	require.Empty(t, string(lines[1]))
}

// TestGenAsciiDoc_Prefix checks that the prefix shows up in the title,
// the usage line, and link text, but not in the anchor ID. A reader
// must see the full command line to run, including the prefix.
func TestGenAsciiDoc_Prefix(t *testing.T) {
	var buf bytes.Buffer

	require.NoError(t, genASCIIDoc(newFixtureRootCmd(), "kubectl ", &buf))
	out := buf.String()

	require.Contains(t, out, "= Command-Line Help for `kubectl root`")
	require.Contains(t, out, "* <<root,`kubectl root`>>")
	require.Contains(t, out, "[#root-sub]\n== `kubectl root sub`")
	require.Contains(t, out, "*Usage:* `kubectl root sub [flags] <arg>`")
	require.Contains(t, out, "* <<root,`kubectl root`>> - Root command")
}

// TestGenAsciiDoc_RealRootCmd checks that the generator works with the
// real kubectl plugin command tree.
func TestGenAsciiDoc_RealRootCmd(t *testing.T) {
	var buf bytes.Buffer

	require.NoError(t, genASCIIDoc(kubectlplugin.NewRootCmd(), "kubectl ", &buf))
	require.NotEmpty(t, buf.String())
	require.Contains(t, buf.String(), "[#runtime-enforcer]")
	require.Contains(t, buf.String(), "== `kubectl runtime-enforcer`")
}
