package main

// Golang provides a base container runtime for Go-based applications with
// pre-configured caching.
type Golang struct {
}

func (m *Golang) Container() *Container {
	buildArgs := []BuildArg{
		BuildArg{
			Name:  "GOLANG_IMAGE",
			Value: "golang:alpine",
		},
	}

	dockerfile := dag.Host().File("Dockerfile")

	return dag.Directory().
		WithFile("./Dockerfile", dockerfile).
		DockerBuild(DirectoryDockerBuildOpts{
			BuildArgs: buildArgs,
		}).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build"))
}
