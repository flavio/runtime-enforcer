package main

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// byName sorts commands by name. This keeps sibling commands in
// alphabetical order in the generated output.
type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

// availableChildren returns the subcommands of cmd, in alphabetical
// order. It skips hidden, deprecated, and help commands.
func availableChildren(cmd *cobra.Command) []*cobra.Command {
	children := cmd.Commands()
	sort.Sort(byName(children))

	out := make([]*cobra.Command, 0, len(children))

	for _, c := range children {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		out = append(out, c)
	}

	return out
}

// hasSeeAlso reports whether cmd needs a "SEE ALSO" section. A command
// needs one when it has a parent or at least one subcommand.
func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}

	return len(availableChildren(cmd)) > 0
}

// anchorID returns the AsciiDoc anchor ID for cmd. It takes the full
// command path and replaces each space with a hyphen. For example, the
// path "runtime-enforcer policy allow" becomes
// "runtime-enforcer-policy-allow". The prefix does not take part in the
// anchor ID, so anchor IDs stay the same when the prefix changes.
func anchorID(cmd *cobra.Command) string {
	return strings.ReplaceAll(cmd.CommandPath(), " ", "-")
}

// displayName returns the name to show to the reader for cmd: prefix
// followed by cmd's full command path. Use this for a title, a
// heading, or link text. Do not use it for the anchor ID.
func displayName(prefix string, cmd *cobra.Command) string {
	return prefix + cmd.CommandPath()
}

// displayUsage returns the usage line to show to the reader for cmd:
// prefix followed by cmd's usage line.
func displayUsage(prefix string, cmd *cobra.Command) string {
	return prefix + cmd.UseLine()
}

// writeOverviewEntry writes one bullet for cmd in the "Command Overview"
// list. It then writes a bullet for each subcommand, in preorder.
func writeOverviewEntry(buf *bytes.Buffer, prefix string, cmd *cobra.Command) {
	fmt.Fprintf(buf, "* <<%s,`%s`>>\n", anchorID(cmd), displayName(prefix, cmd))

	for _, child := range availableChildren(cmd) {
		writeOverviewEntry(buf, prefix, child)
	}
}

// writeOptions writes the "Options" section for cmd's own flags. It
// also writes an "Options inherited from parent commands" section for
// flags cmd inherits. It skips a section when the flag set is empty.
func writeOptions(buf *bytes.Buffer, cmd *cobra.Command) {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)

	if flags.HasAvailableFlags() {
		buf.WriteString("=== Options\n\n....\n")
		flags.PrintDefaults()
		buf.WriteString("....\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)

	if parentFlags.HasAvailableFlags() {
		buf.WriteString("=== Options inherited from parent commands\n\n....\n")
		parentFlags.PrintDefaults()
		buf.WriteString("....\n\n")
	}
}

// writeSeeAlso writes the "SEE ALSO" section for cmd. The section
// links to cmd's parent, if it has one, and to each of its subcommands.
func writeSeeAlso(buf *bytes.Buffer, prefix string, cmd *cobra.Command) {
	if !hasSeeAlso(cmd) {
		return
	}

	buf.WriteString("=== SEE ALSO\n\n")

	if cmd.HasParent() {
		parent := cmd.Parent()
		fmt.Fprintf(buf, "* <<%s,`%s`>> - %s\n", anchorID(parent), displayName(prefix, parent), parent.Short)
	}

	for _, child := range availableChildren(cmd) {
		fmt.Fprintf(buf, "* <<%s,`%s`>> - %s\n", anchorID(child), displayName(prefix, child), child.Short)
	}

	buf.WriteString("\n")
}

// writeCommand writes the AsciiDoc section for one command: an anchor,
// its description, its usage line, its examples, its options, and a
// "SEE ALSO" section. It then writes a section for each subcommand, in
// preorder. This order matches the "Command Overview" list.
//
// prefix is text that goes in front of the command path everywhere the
// reader sees it, for example "kubectl " when the command runs as a
// kubectl plugin. The anchor ID does not use the prefix.
func writeCommand(buf *bytes.Buffer, prefix string, cmd *cobra.Command) {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()
	cmd.InitDefaultVersionFlag()

	fmt.Fprintf(buf, "[#%s]\n", anchorID(cmd))
	fmt.Fprintf(buf, "== `%s`\n\n", displayName(prefix, cmd))

	switch {
	case len(cmd.Long) > 0:
		buf.WriteString(cmd.Long + "\n\n")
	case len(cmd.Short) > 0:
		buf.WriteString(cmd.Short + "\n\n")
	}

	if cmd.Runnable() {
		fmt.Fprintf(buf, "*Usage:* `%s`\n\n", displayUsage(prefix, cmd))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("=== Examples\n\n....\n")
		buf.WriteString(cmd.Example)
		buf.WriteString("\n....\n\n")
	}

	writeOptions(buf, cmd)
	writeSeeAlso(buf, prefix, cmd)

	for _, child := range availableChildren(cmd) {
		writeCommand(buf, prefix, child)
	}
}

// genASCIIDoc writes one AsciiDoc document to w. The document holds
// the help content for root and every one of its subcommands. The
// document starts with a doctitle and a blank line (lines 1 and 2). A
// page in the docs site can include the document with
// `[lines=3..-1]` to skip the doctitle.
//
// prefix is text that goes in front of the command path everywhere the
// reader sees it. Pass "kubectl " for a plugin that always runs through
// kubectl, or "" when the reader runs the binary on its own.
func genASCIIDoc(root *cobra.Command, prefix string, w io.Writer) error {
	buf := new(bytes.Buffer)
	name := displayName(prefix, root)

	fmt.Fprintf(buf, "= Command-Line Help for `%s`\n\n", name)
	fmt.Fprintf(buf, "This document contains the help content for the `%s` command-line program.\n\n", name)

	buf.WriteString("*Command Overview:*\n\n")
	writeOverviewEntry(buf, prefix, root)
	buf.WriteString("\n")

	writeCommand(buf, prefix, root)

	_, err := buf.WriteTo(w)

	return err
}
