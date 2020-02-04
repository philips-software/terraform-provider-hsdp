FROM golang:1.13.7-alpine3.10 as build_base
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

FROM hashicorp/terraform:0.12.20
LABEL maintainer="Andy Lo-A-Foe <andy.lo-a-foe@philips.com>"
ENV HOME /root
COPY --from=builder /terraform-provider-hsdp/build/terraform-provider-hsdp $HOME/.terraform.d/plugins/linux_amd64/terraform-provider-hsdp
