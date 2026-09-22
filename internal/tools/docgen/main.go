package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/kubewarden/runtime-enforcer/internal/kubectlplugin"
)

const defaultOut = "docs/kubectl-plugin/cli-docs.adoc"

func main() {
	out := flag.String("out", defaultOut, "output file for the AsciiDoc CLI reference")
	flag.Parse()

	if err := os.MkdirAll(filepath.Dir(*out), 0o750); err != nil {
		log.Fatal(err)
	}

	_ = os.Setenv("KUBECACHEDIR", "$HOME/.kube/cache")

	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}

	err = genASCIIDoc(kubectlplugin.NewRootCmd(), "kubectl ", f)
	if err != nil {
		_ = f.Close()
		log.Fatal(err)
	}

	// Close is not deferred. A write error can surface at Close, and a
	// truncated file must not exit 0.
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
