---
page_title: "Working with hsdp_function"
---
# Getting started with hsdp_function
The `hsdp_function` resource is a higher level abstraction of the [HSDP Iron](https://www.hsdp.io/documentation/ironio-service-broker) 
service. It uses an Iron service broker instance together with an (optional) function Gateway running in Cloud foundry. This combination
unlocks capabilities beyond the standard Iron services:

- No need to use Iron CLI to schedule tasks or upload code
- Manage Iron codes fully via Terraform
- Periodically schedule Docker workloads using Terraform (**CRONJOB** like functionality)
- Full control over the Docker container **ENVIRONMENT** variables, allowing easy workload configuration
- Automatic encryption of workload payloads 
- Synchronously call a Docker workload running on an Iron Worker, with **streaming support**
- Asynchronously schedule a Docker workload with HTTP Callback support (POST output of workload)
- Function Gateway can be configured with Token auth (default)
- Optionally integrates with **HSDP IAM** for Organization RBAC access to functions
- Asynchronous jobs are scheduled and can take advantage of Iron autoscaling
- Designed to be **Iron agnostic**

# Configuring the backend

The execution plane is pluggable but at this time we only support  the `siderite` backend type which utilizes the HSDP Iron services.
The `siderite` backend should be provisioned using the [siderite-backend](https://registry.terraform.io/modules/philips-labs/siderite-backend/cloudfoundry/latest) terraform module.
Example:

```hcl
module "siderite-backend" {
  source  = "philips-labs/siderite-backend/cloudfoundry"
  version = "0.2.0"

  cf_region   = "eu-west"
  cf_org_name = "my-cf-org"
  cf_space    = "myspacw"
  cf_user     = var.cf_user

  gateway_enabled = true
  auth_type       = "token"
  
  iron_plan   = "large-encrypted-gpu"
}
```
The module will provision an Iron service instance and deploy the function Gateway to the specified
Cloud foundry space. If no space is specified one will be created automatically.

> The (optional) Gateway app is very lean and is set up to use no more than 64MB RAM
 
> Iron service broker plan names can differ between CF regions so make sure the `iron_plan` you specify is available in the region

# Defining your first function

With the above module in place you can continue defining a function:

```hcl
resource "hsdp_function" "hello_world" {
  name         = "hello-world"
  docker_image = "philipslabs/hsdp-function-hello-world:v0.9.2"

  backend {
    type        = "siderite"
    credentials = module.siderite_backend.credentials
  }
}
```

When applied the provider will perform the following actions:
- Create an iron `code` based on the specified docker image
- Create two (2) `schedules` in the Iron backend which use the `code`, one for sychronous calls and one for asychnronous calls

The `hsdp_function` resource will export a number of attributes:

| Name | Description |
|------|-------------|
| `async_endpoint` | The endpoint to trigger your function asychronously |
| `endpoint` | The endpoint to trigger your function synchronously |
| `auth_type` |  The auth type conifguration of the API gateway |
| `token` | The security token to use for authenticating against the endpoints |

# Creating your own Docker function image
A `hsdp_function` compatible Docker image needs to adhere to a number of criteria. We use
a helper application called `siderite`. Siderite started as a convenience tool to ease Iron Worker usage. It now has a 
`function` mode where it will look for an `/app/server` (configurable) and execute it. The server should start up and 
listen on port `8080` for regular HTTP requests. The siderite binary will establish a connection to the gateway and wait
for synchronous requests to come in. In asychronous mode the Siderite helper will pull the payload from the Gateway and
execute the request (again by spawining `/app/server`). Optionally it will `POST` the response back to an URL if one specified in an `X-Callback-URL` header in the HTTP call to the asychronous endpoint.
## Example Docker file
```dockerfile
FROM golang:1.16.0-alpine3.13 as builder
RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /src
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# Build
COPY . .
RUN go build -o server .

FROM philipslabs/siderite:v0.5.1 AS siderite

FROM alpine:latest
RUN apk add --no-cache git openssh openssl bash postgresql-client
WORKDIR /app
COPY --from=siderite /app/siderite /app/siderite
COPY --from=builder /src/server /app

CMD ["/app/siderite","function"]
```
Notes:
- The above docker image builds a Go binary from source and copies it as `/app/server` in the final image
- You can use ANY programming language (even COBOL), as long as you produce an executable binary or script which spawns
  listens on port `8080` after startup.
- We pull the `siderite` binary from the official `philipslabs/siderite` registry. Use a version tag for stability.
- The `CMD` statement should always execute `/app/siderite function` as the main command
- Include any additional tools in your final image

## Naming convention
Please name and publish your `hsdp_function` compatible Docker images as `hsdp-function-xxx` so others can easily identify them.

# Gateway authentication
The gateway supports a number of authentication methods which you can configure via the `auth_type` argument.

| Name | Description |
|------|-------------|
| `none` | Authentication disabled. Only recommended for testing |
| `token` | The default. Token based authentication |
| `iam`  | [HSDP IAM](https://www.hsdp.io/documentation/identity-and-access-management-iam) based authentication |


## Token based authentication
The default authentication method is token based. The endpoint check the following HTTP header for the token

```http
Authorization: Token TOKENHERE
```

if the token matches up the request is allowed.

## IAM integration
The gateway also supports Role Based Access Control (RBAC) using HSDP IAM. The following values should be added to the
siderite backend module block:

```hcl
environment = {
  AUTH_IAM_CLIENT_ID     = "client_id_here"
  AUTH_IAM_CLIENT_SECRET = "Secr3tH3rE"
  AUTH_IAM_REGION        = "eu-west"
  AUTH_IAM_ENVIRONMENT   = "prod"
  AUTH_IAM_ORGS          = "org-uuid1,org-uuid2"
  AUTH_IAM_ROLES         = "HSDP_FUNCTION"
}
```

With the above configuration the gateway will do an introspect call on the Bearer token and if the user/service has the
`HSDP_FUNCTION` role in either of the ORGs specified will be allowed to execute the function.

# Periodically scheduling a function aka CRONJOB
You may schedule a function to run periodically by specifing `schedule` by defining a `schedule` block in your function resource:

```hcl
schedule {
    start = "2021-01-01T04:00:00Z"
    run_every = "1d"
  }
```
The above `schedule` will run your function every day at around `4am`.

> If you only define functions with a `schedule` you can disable the Gateway completely as it is not needed for any operations
