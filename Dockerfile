FROM golang:1.12.7-alpine3.9 as builder
LABEL maintainer="andy.lo-a-foe@philips.com"
RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /terraform-provider-hsdp
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# Build
COPY . .
RUN go build .

FROM hashicorp/terraform:0.12.8
ENV HOME /root
COPY --from=builder /terraform-provider-hsdp/terraform-provider-hsdp $HOME/.terraform.d/plugins/linux_amd64/terraform-provider-hsdp
