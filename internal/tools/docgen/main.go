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
	defer f.Close()

	if err := genAsciiDoc(kubectlplugin.NewRootCmd(), "kubectl ", f); err != nil {
		log.Fatal(err)
	}
}
