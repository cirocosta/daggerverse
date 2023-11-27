ARG GOLANG_IMAGE=golang:alpine
FROM ${GOLANG_IMAGE} AS golang


FROM golang AS source

	ARG SKOPEO_SRC_URL=https://github.com/containers/skopeo
	ARG SKOPEO_SRC_REV=main

	RUN apk add --no-cache --update git

	RUN set -ex && \
		mkdir -p /out && \
		git clone ${SKOPEO_SRC_URL} /out && \
		cd /out && \
		git checkout ${SKOPEO_SRC_REV}

	WORKDIR /out


FROM golang AS runtime

	RUN set -ex && \
		apk add --update \
			btrfs-progs-dev \
			glib-dev \
			gpgme-dev \
			libselinux-dev \
			linux-headers \
			lvm2-dev \
			ostree-dev


FROM runtime AS builder

	RUN set -ex && \
		apk add --update \
			gcc \
			go-md2man \
			musl-dev

	COPY 	--from=source /out /in
	WORKDIR /in

	RUN \
		--mount=target=/root/.cache,type=cache \
		--mount=target=/go/pkg/mod,type=cache  \
		set -ex && \
			GOBIN=/bin go install -v ./cmd/skopeo


FROM runtime AS runner

	COPY	--from=builder /bin/skopeo /bin/skopeo
