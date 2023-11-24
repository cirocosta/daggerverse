package main

import (
	"context"
	"fmt"

	"github.com/google/shlex"
)

const (
	imageGolang = "golang:1.21"

	toolsModuleDirpath = "./internal/tools"
)

type GolangciLint struct{}

// example usage: "dagger shell container-with-source --src ."
func (m *GolangciLint) ContainerWithSource(src *Directory) *Container {
	return dag.Container().
		From(imageGolang).
		WithMountedDirectory("/mnt", src).
		WithWorkdir("/mnt").
		WithFile("/bin/golangci-lint", executableFile())
}

// example usage: "dagger call lint --dir ."
func (m *GolangciLint) Lint(
	ctx context.Context,
	dir *Directory,
) (string, error) {
	return m.ContainerWithSource(dir).
		WithExec(argv(`golangci-lint run`)).
		Stdout(ctx)
}

func argv(s string) []string {
	a, err := shlex.Split(s)
	if err != nil {
		panic(fmt.Errorf("shlex: %w", err))
	}

	return a
}

func golangContainer() *Container {
	return dag.
		Container().
		From(imageGolang)
}

func toolsDirectory() *Directory {
	return dag.
		Host().
		Directory(toolsModuleDirpath, HostDirectoryOpts{})
}

func executableFile() *File {
	return golangContainer().
		WithDirectory("/tools", toolsDirectory()).
		WithWorkdir("/tools").
		WithEnvVariable("GOBIN", "/bin").
		WithExec(argv(`go install -v github.com/golangci/golangci-lint/cmd/golangci-lint`)).
		Directory("/bin").
		File("golangci-lint")
}
