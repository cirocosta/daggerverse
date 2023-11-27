ARG GOLANG_IMAGE=golang:alpine
FROM ${GOLANG_IMAGE} AS golang


FROM golang AS runtime


FROM runtime AS source

	RUN apk add --no-cache --update git

	ARG GOLANGCI_LINT_SRC_URL=https://github.com/golangci/golangci-lint
	ARG GOLANGCI_LINT_SRC_REV=master

	RUN set -ex && \
		mkdir -p /out && \
		git clone ${GOLANGCI_LINT_SRC_URL} /out && \
		cd /out && \
		git checkout ${GOLANGCI_LINT_SRC_REV}

	WORKDIR /out


FROM golang AS builder

	COPY 	--from=source /out /in
	WORKDIR /in

	RUN \
		--mount=target=/root/.cache,type=cache \
		--mount=target=/go/pkg/mod,type=cache  \
		set -ex && \
			GOBIN=/bin go install -v ./cmd/golangci-lint


FROM runtime AS runner

	COPY	--from=builder /bin/golangci-lint /bin/golangci-lint
