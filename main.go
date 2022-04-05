package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	registry "oras.land/oras-go/v2/registry/remote"
)

func main() {

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("No reference specified")
		os.Exit(1)
	}

	ref := flag.Arg(0)
	mnft, err := fetchManifest(ref)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(mnft)
}

func fetchManifest(ref string) (string, error) {
	ref = "localhost:5000/hello-world:latest"
	//	ref = "localhost:5000/hello-world:latest@sha256:f54a58bc1aac5ea1a25d796ae155dc228b3f0e11d046ae276b39c4bf2f13d8c4"
	fmt.Println(ref)
	ctx := context.Background()
	repo, err := registry.NewRepository(ref)
	if err != nil {
		panic(err)
	}

	repo.PlainHTTP = true
	desc, rc, err := repo.FetchReference(ctx, ref)
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
