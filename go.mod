module github.com/sajayantony/oras-manifest

go 1.17

require (
	oras.land/oras-go/pkg v0.0.0
	oras.land/oras-go/v2 v2.0.0
)

require (
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/oras-project/artifacts-spec v1.0.0-draft.1.1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)

replace (
	oras.land/oras-go/pkg => /home/sajay/code/src/github.com/sajayantony/oras-go/pkg
	oras.land/oras-go/v2 => /home/sajay/code/src/github.com/sajayantony/oras-go
)
