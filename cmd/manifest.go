package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	registry "oras.land/oras-go/v2/registry/remote"
)

type manifestOptions struct {
	targetRef string
}

func manifestCmd() *cobra.Command {
	var opts refsOptions
	cmd := &cobra.Command{
		Use:   "manifest <name:tag|name@digest>",
		Short: "Show manifest of the given reference",

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.targetRef = args[0]
			return runManifest(opts)
		},
	}

	return cmd
}

func runManifest(opts refsOptions) error {

	mnft, err := fetchManifest(opts.targetRef)
	if err == nil {
		fmt.Print(mnft)
	}
	return err
}

func fetchManifest(ref string) (string, error) {
	//	ref = "localhost:5000/hello-world:latest@sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4"
	// fmt.Println(ref)
	ctx := context.Background()
	repo, err := registry.NewRepository(ref)
	if err != nil {
		panic(err)
	}
	setPlainHttp(repo)

	_, rc, err := repo.FetchReference(ctx, ref)
	if err != nil {
		panic(err)
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, rc)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
