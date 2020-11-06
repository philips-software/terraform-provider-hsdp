ARG hsdp_provider_version=0.6.9

FROM golang:1.15.2-alpine3.12 as build_base
RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /terraform-provider-hsdp
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
LABEL builder=true

# Build
FROM build_base AS builder
COPY . .
RUN ./buildscript.sh

FROM hashicorp/terraform:0.13.5
ARG hsdp_provider_version
ENV HSDP_PROVIDER_VERSION ${hsdp_provider_version}
LABEL maintainer="Andy Lo-A-Foe <andy.lo-a-foe@philips.com>"
ENV HOME /root
COPY --from=builder /terraform-provider-hsdp/build/terraform-provider-hsdp $HOME/.terraform.d/plugins/registry.terraform.io/philips-software/hsdp/${HSDP_PROVIDER_VERSION}/linux_amd64/terraform-provider-hsdp_v${HSDP_PROVIDER_VERSION}
