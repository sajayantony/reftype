package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	memorytarget "oras.land/oras-go/v2/content/memory"
	registry "oras.land/oras-go/v2/registry/remote"
	orasauth "oras.land/oras-go/v2/registry/remote/auth"
)

func main() {

}

func getManifest() (string, error) {
	plainHTTP := false
	insecure := false
	username := "test"
	password := "testpw"
	sourceImage := "test.azurecr.io/repo/image:tag"
	// Create new memory store
	memory_store := memorytarget.New()
	// Create Repository Target
	repository, err := registry.NewRepository(sourceImage)
	if err != nil {
		return "", err
	}

	credentialProvider := func(ctx context.Context, registry string) (orasauth.Credential, error) {
		if username != "" || username != "" {
			return orasauth.Credential{
				Username: username,
				Password: password,
			}, nil
		}
		//TODO: handle dockerconfig
		return orasauth.EmptyCredential, nil
	}

	// Set the Repository Client Credentials
	repoClient := &orasauth.Client{
		Header: http.Header{
			"User-Agent": {"oras-go"},
		},
		Cache:      orasauth.DefaultCache,
		Credential: credentialProvider,
	}
	// Set the TSLClientConfig for HTTP client if insecure set to true
	if insecure {
		repoClient.Client = http.DefaultClient
		repoClient.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	// Set the PlainHTTP to true if specified
	repository.PlainHTTP = plainHTTP
	repository.Client = repoClient
	// Copy the remote source image to local memory target
	retDesc, err := oras.Copy(context.Background(), repository, sourceImage, memory_store, "")
	if err != nil {
		return "", err
	}
	manifestReader, err := memory_store.Fetch(context.Background(), retDesc)
	if err != nil {
		return "", err
	}
	var returnedManifest ocispec.Manifest
	returnedManifestBytes, err := ioutil.ReadAll(manifestReader)
	if err != nil {
		return "", err
	}
	json.Unmarshal(returnedManifestBytes, &returnedManifest)
	fmt.Println(returnedManifest)

	return returnedManifest, err
}
