ARG GOLANG_IMAGE=golang:alpine
FROM ${GOLANG_IMAGE} AS golang


FROM golang AS runtime

	RUN apk add --no-cache --update git


FROM runtime AS source

	ARG GITLEAKS_SRC_URL=https://github.com/gitleaks/gitleaks
	ARG GITLEAKS_SRC_REV=master

	RUN set -ex && \
		mkdir -p /out && \
		git clone ${GITLEAKS_SRC_URL} /out && \
		cd /out && \
		git checkout ${GITLEAKS_SRC_REV}

	WORKDIR /out


FROM golang AS builder

	COPY 	--from=source /out /in
	WORKDIR /in

	RUN \
		--mount=target=/root/.cache,type=cache \
		--mount=target=/go/pkg/mod,type=cache  \
		set -ex && \
			GOBIN=/bin go install -v .


FROM runtime AS runner

	COPY	--from=builder /bin/gitleaks /bin/gitleaks
