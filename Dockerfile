FROM golang:1.11.0-alpine3.8 as builder
LABEL maintainer="andy.lo-a-foe@philips.com"

RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /terraform-provider-hsdp
COPY . /terraform-provider-hsdp
RUN cd /terraform-provider-hsdp && go build -o terraform-provider-hsdp .

FROM hashicorp/terraform
ENV HOME /root
COPY --from=builder /terraform-provider-hsdp/terraform-provider-hsdp $HOME/.terraform.d/plugins/linux_amd64/terraform-provider-hsdp
ENTRYPOINT ["/bin/terraform"]
