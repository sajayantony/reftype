package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	artifactspec "github.com/oci-playground/artifact-spec/specs-go/v1"
	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry/remote"
)

const (
	tagStaged = "staged"
)

type pushOptions struct {
	targetRef    string
	fileRefs     []string
	artifactType string
	verbose      bool
}

func pushCmd() *cobra.Command {
	var opts pushOptions
	cmd := &cobra.Command{
		Use:   "push name[:tag|@digest] file[:type] [file...]",
		Short: "Push a reference type artifact",
		Long: `Push reference artifact to remote registry

Example - Push reference artifact "hi.txt" with the "application/vnd.oci.image.layer.v1.tar" media type (default):
  reftype push localhost:5000/hello:latest hi.txt

Example - Push reference artifact "hi.txt" with the custom "application/vnd.me.hi" media type:
  reftype push localhost:5000/hello:latest hi.txt:application/vnd.me.hi

Example - Push multiple files with different media types as a reference artifact:
  reftype push localhost:5000/hello:latest hi.txt:application/vnd.me.hi bye.txt:application/vnd.me.bye
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.targetRef = args[0]
			opts.fileRefs = args[1:]
			return runPush(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.artifactType, "artifact-type", "", "", "artifact type")
	return cmd
}

func runPush(opts pushOptions) error {

	// Prepare client
	repo, err := remote.NewRepository(opts.targetRef)
	//repo.ManifestMediaTypes = append(repo.ManifestMediaTypes, artifactspec.MediaTypeArtifactManifest)
	if err != nil {
		return err
	}
	setPlainHttp(repo)

	// Prepare manifest
	store := file.New("")
	defer store.Close()

	// Pack manifests
	var desc ocispec.Descriptor
	ctx := context.Background()

	desc, err = packArtifact(ctx, repo, store, &opts)

	if err != nil {
		return err
	}

	// ready to push
	tracker := &statusTracker{
		Target:  repo,
		out:     os.Stdout,
		prompt:  "Uploading",
		verbose: opts.verbose,
	}

	err = oras.CopyGraph(ctx, store, tracker, desc)

	if err != nil {
		return err
	}

	fmt.Println("Pushed", opts.targetRef)
	fmt.Println("Digest:", desc.Digest)

	return nil
}

func packArtifact(ctx context.Context, remote content.Resolver, store *file.Store, opts *pushOptions) (ocispec.Descriptor, error) {
	subject, err := remote.Resolve(ctx, opts.targetRef)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	files, err := loadFiles(ctx, store, opts)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	var annotations = map[string]string{}
	annotations["org.opencontainers.artifact.type"] = opts.artifactType
	annotations["org.opencontainers.artifact.created"] = time.Now().Format(time.RFC3339)

	manifest := artifactspec.ArtifactManifest{
		MediaType: artifactspec.MediaTypeArtifactManifest,
		//ArtifactType: opts.artifactType,
		Blobs:       ociToArtifactSlice(files),
		Reference:   ociToArtifact(subject),
		Annotations: annotations,
		//Annotations: annotations[annotationManifest],
	}

	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	manifestDesc := ocispec.Descriptor{
		MediaType: artifactspec.MediaTypeArtifactManifest,
		Digest:    digest.FromBytes(manifestBytes),
		Size:      int64(len(manifestBytes)),
	}

	// store manifest
	if err := store.Push(ctx, manifestDesc, bytes.NewReader(manifestBytes)); err != nil && !errors.Is(err, errdef.ErrAlreadyExists) {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push manifest: %w", err)
	}
	if err := store.Tag(ctx, manifestDesc, tagStaged); err != nil {
		return ocispec.Descriptor{}, err
	}
	return manifestDesc, nil
}

func loadFiles(ctx context.Context, store *file.Store, opts *pushOptions) ([]ocispec.Descriptor, error) {
	var files []ocispec.Descriptor
	for _, fileRef := range opts.fileRefs {
		filename, mediaType := parseFileRef(fileRef, "")
		name := filepath.Clean(filename)
		if !filepath.IsAbs(name) {
			// convert to slash-separated path unless it is absolute path
			name = filepath.ToSlash(name)
		}
		if opts.verbose {
			fmt.Println("Preparing", name)
		}
		file, err := store.Add(ctx, name, mediaType, filename)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}
	if len(files) == 0 {
		fmt.Println("Uploading empty artifact")
	}
	return files, nil
}

func ociToArtifactSlice(descs []ocispec.Descriptor) []ocispec.Descriptor {
	res := make([]ocispec.Descriptor, 0, len(descs))
	for _, desc := range descs {
		res = append(res, ociToArtifact(desc))
	}
	return res
}

func ociToArtifact(desc ocispec.Descriptor) ocispec.Descriptor {
	return ocispec.Descriptor{
		MediaType:   desc.MediaType,
		Digest:      desc.Digest,
		Size:        desc.Size,
		URLs:        desc.URLs,
		Annotations: desc.Annotations,
	}
}
