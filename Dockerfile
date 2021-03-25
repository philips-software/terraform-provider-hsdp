# syntax = docker/dockerfile:1-experimental

ARG hsdp_provider_version=0.12.10
FROM --platform=${BUILDPLATFORM} golang:1.16.2-alpine3.13 AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/terraform-provider-hsdp -ldflags "-X main.GitCommit=${GIT_COMMIT}" .

FROM hashicorp/terraform:0.14.9
RUN apk add --no-cache tzdata
ARG TARGETOS
ARG TARGETARCH
ARG hsdp_provider_version
ENV HSDP_PROVIDER_VERSION ${hsdp_provider_version}
ENV HOME /root
COPY --from=build /out/terraform-provider-hsdp $HOME/.terraform.d/plugins/registry.terraform.io/philips-software/hsdp/${HSDP_PROVIDER_VERSION}/${TARGETOS}_${TARGETARCH}/terraform-provider-hsdp_v${HSDP_PROVIDER_VERSION}
