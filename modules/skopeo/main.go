package main

const (
	// srcImageFpath is the absolute path inside the skopeo container where
	// the oci archive gets put into to be used by commands that relocate
	// from a tarball.
	//
	srcImageFpath = "/in/image.tar"

	// digestFpath is the absolute path inside the skopeo container where
	// the digest of relocated images gets written to.
	//
	digestFpath = "/out/digest.txt"
)

// Skopeo provides the ability of relocating images:
// - from an oci archive to a remote registry (see `RelocateOCITarballToRef`)
// - from a remote registry to another (see `RelocateRefToRef`).
type Skopeo struct {
}

// RelocateRefToRef relocates images from a registry to another.
func (s *Skopeo) RelocateRefToRef(
	src string,
	dst string,
) *Container {
	return s.Container().
		WithExec(skopeoCopyArgv(
			registryRef(src),
			registryRef(dst),
		))
}

// RelocateOCITarballToRef relocates an image from an oci archive to a
// registry.
func (s *Skopeo) RelocateOCITarballToRef(
	src *File,
	dst string,
) *Container {
	return s.Container().
		WithFile(srcImageFpath, src).
		WithExec(skopeoCopyArgv(
			ociTarballRef(srcImageFpath),
			registryRef(dst),
		))
}

// Container provides a container with `skopeo`.
//
// Note that `skopeo` is not a single static binary that can be copied from
// this container and expect it to work on any other container - it requires a
// few runtime dependencies included in the image (see `./Dockerfile`).
func (s *Skopeo) Container() *Container {
	buildArgs := []BuildArg{
		BuildArg{
			Name:  "GOLANG_IMAGE",
			Value: "golang:alpine",
		},
		BuildArg{
			Name:  "SKOPEO_SRC_URL",
			Value: "https://github.com/containers/skopeo",
		},
		BuildArg{
			Name:  "SKOPEO_SRC_REV",
			Value: "v1.14.0",
		},
	}

	dockerfile := dag.Host().File("Dockerfile")

	return dag.Directory().
		WithFile("./Dockerfile", dockerfile).
		DockerBuild(DirectoryDockerBuildOpts{
			BuildArgs: buildArgs,
		}).
		WithExec([]string{"mkdir", "-p", "/out"}).
		WithExec([]string{"mkdir", "-p", "/in"})
}

func skopeoCopyArgv(src, dst string) []string {
	return []string{
		"skopeo", "copy",
		"--insecure-policy",
		"--digestfile", digestFpath,
		src, dst,
	}
}

func registryRef(v string) string {
	return "docker://" + v
}

func ociTarballRef(v string) string {
	return "oci-archive://" + v
}
