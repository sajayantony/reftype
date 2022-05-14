module github.com/sajayantony/reftype

go 1.17

require (
	github.com/oci-playground/artifact-spec v0.0.0-20220506233500-8fed0a29d06f
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.2
	github.com/spf13/cobra v1.4.0
	oras.land/oras-go/v2 v2.0.0-00010101000000-000000000000
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/oras-project/artifacts-spec v1.0.0-draft.1.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)

replace oras.land/oras-go/v2 => github.com/sajayantony/oras-go/v2 v2.0.0-20220514032313-2aca3fbf5e9d
