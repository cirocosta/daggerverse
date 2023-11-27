package main

const (
	// sourceCodeMountPath is the absolute path to the directory where
	// source code will be mounted when performing direct source code
	// scanning.
	sourceCodeMountPath = "/in/src"

	// srcImageFpath is the absolute path inside the trivy container where
	// the oci archive gets put into to be used by the `trivy image`
	// command to check for dependencies.
	//
	srcImageFpath = "/in/image.tar"

	// outputFpath is the absolute path inside the trivy container where
	// the sbom is written to.
	//
	outputFpath = "/out/result.json"
)

// Trivy provides the ability of detecting possible credentials leaks in a
// codebase.
type Trivy struct {
}

// ScanSource provides a container for running trivy's scanning on the source
// directory provided.
func (m *Trivy) ScanSource(src *Directory) *Container {
	return m.Container().
		WithMountedDirectory(sourceCodeMountPath, src).
		WithWorkdir(sourceCodeMountPath).
		WithExec([]string{
			"trivy", "filesystem",
			"--format", "spdx-json",
			"--output", outputFpath,
			".",
		})
}

// ScanSource provides a container for running trivy's scanning on the image
// reference provided.
func (m *Trivy) ScanImageRef(imageRef string) *Container {
	return m.Container().
		WithExec([]string{
			"trivy", "image",
			"--format", "spdx-json",
			"--output", outputFpath,
			imageRef,
		})
}

// ScanOCITarball provides a container for running trivy's scanning on the oci
// image archive provided.
func (m *Trivy) ScanOCITarball(src *File) *Container {
	return m.Container().
		WithMountedFile(srcImageFpath, src).
		WithExec([]string{
			"trivy", "image",
			"--format", "spdx-json",
			"--input", srcImageFpath,
			"--output", outputFpath,
		})
}

// OutputFile is a helper function for extracting the output of the scan
// directly from a container that represents the scanning of "something"
// (image-ref/source/oci-archive).
func OutputFile(container *Container) *File {
	return container.File(outputFpath)
}

// Container provides a container with `trivy` to run detections from.
func (m *Trivy) Container() *Container {
	buildArgs := []BuildArg{
		BuildArg{
			Name:  "GOLANG_IMAGE",
			Value: "golang:alpine",
		},
		BuildArg{
			Name:  "TRIVY_SRC_URL",
			Value: "https://github.com/aquasecurity/trivy",
		},
		BuildArg{
			Name:  "TRIVY_SRC_REV",
			Value: "v0.47.0",
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
