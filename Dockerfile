FROM golang:alpine as builder
LABEL maintainer="andy.lo-a-foe@philips.com"

RUN apk add --no-cache git openssh
ADD . $GOPATH/src/github.com/philips-software/terraform-provider-hsdp
WORKDIR $GOPATH/src/github.com/philips-software/terraform-provider-hsdp
RUN go get . && go build .

FROM hashicorp/terraform
ENV GOPATH /go
ENV HOME /root
COPY --from=builder $GOPATH/src/github.com/philips-software/terraform-provider-hsdp/terraform-provider-hsdp $HOME/.terraform.d/plugins/linux_amd64/terraform-provider-hsdp
ENTRYPOINT ["/bin/terraform"]
