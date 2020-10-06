ARG hsdp_provider_version=0.6.3

FROM alpine:latest AS cf
ARG cf_provider_version
ENV CF_PROVIDER_VERSION ${cf_provider_version}
WORKDIR /build
RUN apk update \
 && apk add curl

# Verify the signature file is untampered.
RUN curl -L -Os https://github.com/cloudfoundry-community/terraform-provider-cf/releases/download/v${CF_PROVIDER_VERSION}/terraform-provider-cloudfoundry_linux_amd64
RUN curl -L -Os https://github.com/cloudfoundry-community/terraform-provider-cf/releases/download/v${CF_PROVIDER_VERSION}/checksums.txt

RUN CHECKSUM=$(cat checksums.txt |grep linux_amd64|grep -v zip|cut -f 1 -d ' ') && \
    echo ${CHECKSUM}"  "terraform-provider-cloudfoundry_linux_amd64 |sha1sum -c
RUN chmod +x terraform-provider-cloudfoundry_linux_amd64

FROM golang:1.15-alpine3.12 as build_base
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

FROM hashicorp/terraform:0.13.2
ENV CF_PROVIDER_VERSION ${cf_provider_version}
ENV HSDP_PROVIDER_VERSION ${hsdp_provider_version}
LABEL maintainer="Andy Lo-A-Foe <andy.lo-a-foe@philips.com>"
ENV HOME /root
COPY --from=builder /terraform-provider-hsdp/build/terraform-provider-hsdp $HOME/.terraform.d/plugins/philips.com/hsdp/hsdp/${HSDP_PROVIDER_VERSION}/linux_amd64/terraform-provider-hsdp_v${HSDP_PROVIDER_VERSION}
