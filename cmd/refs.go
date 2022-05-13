package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	registry "oras.land/oras-go/v2/registry/remote"
)

func main() {

	mnftcmd := flag.NewFlagSet("manifest", flag.ExitOnError)
	refscmd := flag.NewFlagSet("refs", flag.ExitOnError)

	if len(os.Args) < 2 {
		// switch to better CLI parsing
		fmt.Print("No subcommand specified")
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "manifest":
		err := mnftcmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if mnftcmd.NArg() == 0 {
			fmt.Println("No reference specified")
			os.Exit(1)
		}

		ref := mnftcmd.Arg(0)
		mnft, err := fetchManifest(ref)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println(mnft)

	case "ls":
		err := refscmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if refscmd.NArg() == 0 {
			fmt.Println("No reference specified")
			os.Exit(1)
		}

		ref := refscmd.Arg(0)
		if err := fetchReferrers(ref); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Invalid subcommand %s\n", cmd)
		os.Exit(1)
	}
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
		for _, r := range refs {
			d, err := json.Marshal(r)
			if err == nil {
				fmt.Println(string(d))
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func setPlainHttp(repo *registry.Repository) {
	if repo.Reference.Host() == "localhost" || strings.HasPrefix(repo.Reference.Host(), "localhost:") {
		repo.PlainHTTP = true
	}
}
