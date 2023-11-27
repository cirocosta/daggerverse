package main

import (
	"context"
)

const (
	defaultGolangImage = "golang:alpine"
	defaultGitUrl      = "https://github.com/golangci/golangci-lint"
	defaultGitRev      = "v1.55.2"

	toolsModuleDirpath = "./internal/tools"
)

type GolangciLint struct{}

// example usage: "dagger shell container-with-source --src ."
func (m *GolangciLint) ContainerWithSource(src *Directory) *Container {
	return m.Container().
		WithMountedDirectory("/mnt", src).
		WithWorkdir("/mnt")
}

// example usage: "dagger call lint --src ."
func (m *GolangciLint) Lint(
	ctx context.Context,
	src *Directory,
) (string, error) {
	runArgv := []string{
		"golangci-lint", "run",
	}

	return m.ContainerWithSource(src).
		WithExec(runArgv).
		Stdout(ctx)
}

func (m *GolangciLint) Container() *Container {
	buildArgs := []BuildArg{
		BuildArg{
			Name:  "GOLANG_IMAGE",
			Value: defaultGolangImage,
		},
		BuildArg{
			Name:  "GOLANGCI_LINT_SRC_URL",
			Value: defaultGitUrl,
		},
		BuildArg{
			Name:  "GOLANGCI_LINT_SRC_REV",
			Value: defaultGitRev,
		},
	}

	dockerfile := dag.Host().File("Dockerfile")

	return dag.Directory().
		WithFile("./Dockerfile", dockerfile).
		DockerBuild(DirectoryDockerBuildOpts{
			BuildArgs: buildArgs,
		})
}
