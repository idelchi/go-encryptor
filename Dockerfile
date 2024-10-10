#[=======================================================================[
# Description : Docker image containing the gocry binary
#]=======================================================================]

# Docker image repository to use. Use `docker.io` for public images.
ARG IMAGE_BASE_REGISTRY

ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.3

#### ---- Build ---- ####
FROM ${IMAGE_BASE_REGISTRY}golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build

LABEL maintainer=arash.idelchi

# (can use root throughout the image since it's a staged build)
# hadolint ignore=DL3002
USER root

# Basic good practices
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# timezone
RUN apk add --no-cache \
    tzdata

WORKDIR /work

ARG GOMODCACHE=/home/${USER}/.cache/.go-mod
ARG GOCACHE=/home/${USER}/.cache/.go

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    go mod download

COPY . .
ARG GOCRY_VERSION="unofficial & built by unknown"
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    CGO_ENABLED=0 go install -ldflags="-s -w -X 'main.version=${GOCRY_VERSION}'" ./...

# Create User (Alpine)
ARG USER=user
RUN addgroup -S -g 1001 ${USER} && \
    adduser -S -u 1001 -G ${USER} -h /home/${USER} -s /bin/ash ${USER}

# Timezone
ENV TZ=Europe/Zurich

#### ---- Standalone ---- ####
FROM scratch AS standalone

LABEL maintainer=arash.idelchi

ARG USER=user

# Copy artifacts from the build stage
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /go/bin/gocry /gocry

USER ${USER}
WORKDIR /home/${USER}

# Clear the base image entrypoint
ENTRYPOINT ["/gocry"]
CMD [""]

# Timezone
ENV TZ=Europe/Zurich

#### ---- App ---- ####
FROM ${IMAGE_BASE_REGISTRY}alpine:${ALPINE_VERSION}

LABEL maintainer=arash.idelchi

USER root

# timezone
RUN apk add --no-cache \
    git \
    tzdata \
    coreutils \
    jq \
    openssl

# Copy artifacts from the build stage
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/bin/go-encryptor /usr/local/bin/go-encryptor

USER ${USER}
WORKDIR /home/${USER}

# Clear the base image entrypoint
ENTRYPOINT [""]
CMD ["/bin/ash"]

# Timezone
ENV TZ=Europe/Zurich
