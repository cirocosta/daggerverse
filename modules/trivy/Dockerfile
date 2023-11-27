ARG GOLANG_IMAGE=golang:alpine
FROM ${GOLANG_IMAGE} AS golang


FROM golang AS runtime


FROM runtime AS source

	RUN apk add --no-cache --update git

	ARG TRIVY_SRC_URL=https://github.com/aquasecurity/trivy
	ARG TRIVY_SRC_REV=main

	RUN set -ex && \
		mkdir -p /out && \
		git clone ${TRIVY_SRC_URL} /out && \
		cd /out && \
		git checkout ${TRIVY_SRC_REV}

	WORKDIR /out


FROM golang AS builder

	COPY 	--from=source /out /in
	WORKDIR /in

	RUN \
		--mount=target=/root/.cache,type=cache \
		--mount=target=/go/pkg/mod,type=cache  \
		set -ex && \
			GOBIN=/bin go install -v ./cmd/trivy


FROM runtime AS runner

	COPY	--from=builder /bin/trivy /bin/trivy
