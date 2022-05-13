package main

import (
	"strings"

	registry "oras.land/oras-go/v2/registry/remote"
)

func setPlainHttp(repo *registry.Repository) {
	if repo.Reference.Host() == "localhost" || strings.HasPrefix(repo.Reference.Host(), "localhost:") {
		repo.PlainHTTP = true
	}
}
