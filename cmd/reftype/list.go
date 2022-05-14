package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	registry "oras.land/oras-go/v2/registry/remote"
)

type refsOptions struct {
	targetRef string
}

func refsCmd() *cobra.Command {
	var opts refsOptions
	cmd := &cobra.Command{
		Use:   "ls <name:tag|name@digest>",
		Short: "Lists ref-types for a given reference",

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.targetRef = args[0]
			return runRefs(opts)
		},
	}

	return cmd
}

func runRefs(opts refsOptions) error {
	return fetchReferrers(opts.targetRef)
}

func fetchReferrers(ref string) error {

	ctx := context.Background()
	repo, err := registry.NewRepository(ref)

	if err != nil {
		panic(err)
	}

	setPlainHttp(repo)

	desc, err := repo.Manifests().Resolve(ctx, ref)
	if err != nil {
		return err
	}

	if err = repo.Referrers(ctx, desc, func(refs []ocispec.Descriptor) error {
		var index ocispec.Index
		index.Manifests = refs
		d, err := json.Marshal(index)
		if err == nil {
			fmt.Println(string(d))
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
