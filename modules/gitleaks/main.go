package main

const (
	// reportFpath is the absolute path inside the gitleaks container where
	// reports from `gitleaks detect` are written to.
	//
	// See `.ReportFile()` for having a direct reference to the report.
	reportFpath = "/out/report.json"

	// sourceCodeMountPath is the absolute path to the directory where
	// source code will be mounted for credential detection.
	sourceCodeMountPath = "/in/src"
)

// Gitleaks provides the ability of detecting possible credentials leaks in a
// codebase.
type Gitleaks struct{}

// Detect runs credential detection on source code provided via the `src` Directory.
func (m *Gitleaks) Detect(src *Directory) *Container {
	return m.Container().
		WithMountedDirectory(sourceCodeMountPath, src).
		WithWorkdir(sourceCodeMountPath).
		WithExec([]string{
			"gitleaks", "detect",
			"-v",
			"--source", ".",
			"--report-path", reportFpath,
		})
}

// ReportFile provides a reference to the report file that will be generated once the detection runs on `src`.
func (m *Gitleaks) ReportFile(src *Directory) *File {
	return m.Detect(src).
		File(reportFpath)
}

// Container provides a container with `gitleaks` to run detections from.
//
// Note that `gitleaks` itself is compiled in the form of a static binary but
// it relies on `git` being available in $PATH, which the container image here
// has.
func (m *Gitleaks) Container() *Container {
	buildArgs := []BuildArg{
		BuildArg{
			Name:  "GOLANG_IMAGE",
			Value: "golang:alpine",
		},
		BuildArg{
			Name:  "GITLEAKS_SRC_URL",
			Value: "https://github.com/gitleaks/gitleaks",
		},
		BuildArg{
			Name:  "GITLEAKS_SRC_REV",
			Value: "v8.18.1",
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
