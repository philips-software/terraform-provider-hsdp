---
page_title: "Working with HSDP functions"
---
# Working with HSDP functions

The `hsdp_function` resource is a higher level abstraction of the [HSDP Iron](https://www.hsdp.io/documentation/ironio-service-broker)
service. It uses an Iron service broker instance together with an (optional) function Gateway running on Cloud foundry. This combination
unlocks capabilities beyond the standard Iron services:

- No need to use Iron CLI to schedule tasks or upload code
- Manage Iron codes fully via Terraform
- **CRON** compatible scheduling of Docker workloads using Terraform, leapfrogging Iron.io scheduling capabilities
- Full control over the Docker container **ENVIRONMENT** variables, allowing easy workload configuration
- Automatic encryption of workload payloads
- Synchronously call a Docker workload running on an IronWorker, with **streaming support**
- Asynchronously schedule a Docker workload with HTTP Callback support (POST output of workload)
- Function Gateway can be configured with Token auth (default)
- Optionally integrates with **HSDP IAM** for Organization Role Based Access Control (RBAC) to functions
- Asynchronous jobs are scheduled and can take advantage of Iron autoscaling
- Designed to be **Iron agnostic**

## Configuring the backend

The execution plane is pluggable but at this time we  support  the `siderite` backend type which utilizes the HSDP Iron services.
The `siderite` backend should be provisioned using the [siderite-backend](https://registry.terraform.io/modules/philips-labs/siderite-backend/cloudfoundry/latest) terraform module.
Example:

```hcl
module "siderite-backend" {
  source  = "philips-labs/siderite-backend/cloudfoundry"
  version = "0.8.0"

  cf_region   = "eu-west"
  cf_org_name = "my-cf-org"
  cf_space    = "myspace"
  cf_user     = var.cf_user

  iron_plan   = "large-encrypted-gpu"
}
```

> Iron service broker plan names can differ between CF regions so make sure the `iron_plan` you specify is available in the region

The module will provision an Iron service instance and deploy the function Gateway to the specified
Cloud foundry space. If no space is specified one will be created automatically.

> The (optional) Gateway app is very lean and is set up to use no more than 64MB RAM

## Defining your first function

With the above module in place you can continue on to defining a function:

```hcl
resource "hsdp_function" "cuda_test" {
  name         = "cuda-test"
  docker_image = "philipslabs/hsdp-task-cuda-test:v0.0.4"
  command      = ["/app/cudatest"]
  
  backend {
    credentials = module.siderite_backend.credentials
  }
}
```

When applied, the provider will perform the following actions:

- Create an iron `code` based on the specified docker image
- Create two (2) `schedules` in the Iron backend which use the `code`, one for synchronous calls and one for asynchronous calls

The `hsdp_function` resource will export a number of attributes:

| Name | Description |
|------|-------------|
| `async_endpoint` | The endpoint to trigger your function asychronously |
| `endpoint` | The endpoint to trigger your function synchronously |
| `auth_type` |  The auth type configuration of the API gateway |
| `token` | The security token to use for authenticating against the endpoints |

## Creating your own Docker function image

A `hsdp_function` compatible Docker image needs to adhere to a number of criteria. We use
a helper application called `siderite`. Siderite started as a convenience tool to ease IronWorker usage. It now has a
`function` mode where it will look for an `/app/server` (configurable) and execute it. The `siderite` helper also
provides integration with [HSDP Logging](https://www.hsdp.io/documentation/logging) and supports streaming output in
real-time to your central logging accounts.

## Asynchronous function

In asynchronous mode the Siderite helper will pull the payload from the Gateway and execute the request
(again by spawning `/app/server`). It will `POST` the response back to a URL specified in the original request Header called
`X-Callback-URL` header.

## Example Docker file

```dockerfile
FROM golang:1.17.3-alpine3.14 as builder
RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /src
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# Build
COPY . .
RUN go build -o server .

FROM philipslabs/siderite:v0.12.0 AS siderite

FROM alpine:latest
WORKDIR /app
COPY --from=siderite /app/siderite /app/siderite
COPY --from=builder /src/server /app/server

CMD ["/app/siderite","function"]
```

Notes:

- The above docker image builds a Go binary from source and copies it as `/app/server` in the final image
- You can use ANY programming language (even COBOL), as long as you produce an executable binary or script which spawns
  listens on port `8080` after startup.
- We pull the `siderite` binary from the official `philipslabs/siderite` registry. Use a version tag for stability.
- The `CMD` statement should execute `/app/siderite function` as the main command
- If your function is always scheduled use `/app/siderite task` instead. This will automatically exit after a single run.
- Include any additional tools in your final image

## Example using curl

```text
curl -v \
    -X POST \
    -H "Authorization: Token XXX" \
    -H "X-Callback-URL: https://hookb.in/XYZ" \
    https://hsdp-func-gateway-yyy.eu-west.philips-healthsuite.com/function/zzz
```

This would schedule the function to run. The result of the request will then be posted to `https://hook.bin/XYZ`. Calls
will be queued up and picked up by workers.

## Scheduling a function to run periodically (Task)

Enabling the gateway in the `siderite` backend unlocks full **CRON** compatible scheduling of `hsdp_function` resources.
It provides much finer control over scheduling behaviour compared to the standard Iron.io `run_every`
option. To achieve this the gateway runs an internal CRON scheduler which is driven by the provider managed schedule entries
in the Iron.io backend, syncing the config every few seconds.

```hcl
resource "hsdp_function" "cuda_test" {
  name         = "cuda-test"
  docker_image = "philipslabs/hsdp-task-cuda-test:v0.0.5"
  command      = ["/app/cudatest"]

  schedule = "14 15 * * * *"
  timeout  = 120

  backend {
    credentials = module.siderite_backend.credentials
  }
}
```

The above example would queue your `hsdp_function` every day at exactly 3.14pm.
The following one would queue your function every Sunday morning at 5am:

```hcl
resource "hsdp_function" "cuda_test" {
  name         = "cuda-test"
  docker_image = "philipslabs/hsdp-task-cuda-test:v0.0.4"
  command      = ["/app/cudatest"]

  schedule = "0 5 * * * 0"
  timeout  = 120

  backend {
    credentials = module.siderite_backend.credentials
  }
}
```

-> Even though you can specify an up-to-the-minute accurate schedule, your function is still queued on the
Iron cluster, so the exact start time is always determined by how busy the cluster is at that moment.

Finally, an example of using the Iron.io native scheduler:

```hcl
resource "hsdp_function" "cuda_test" {
  name         = "cuda-test"
  docker_image = "philipslabs/hsdp-task-cuda-test:v0.0.5"
  command      = ["/app/cudatest"]

  run_every = "20m"
  start_at = "2021-01-01T07:00:00.00Z" # Start at 7am UTC
  timeout  = 120

  backend {
    credentials = module.siderite_backend.credentials
  }
}
```

This will run your function every 20 minutes.

-> Always set a timeout value for your scheduled function. This sets a limit on the runtime for each invocation.

### cron field description

```text
1. Entry: Minute when the process will be started [0-60]
2. Entry: Hour when the process will be started [0-23]
3. Entry: Day of the month when the process will be started [1-28/29/30/31]
4. Entry: Month of the year when the process will be started [1-12]
5. Entry: Weekday when the process will be started [0-6] [0 is Sunday]

all x min = */x
```

## Function vs Task

The `hsdp_function` resource supports defining functions which are automatically executed
periodically i.e. `Tasks`. A Docker image which defines a task should use the following `CMD`:

```dockerfile
FROM philipslabs/siderite:debian-v0.12.0 as siderite

FROM minio/mc:latest
COPY --from=siderite /app/siderite /usr/bin/siderite

## Copy other tools or applications as needed
#COPY s3mirror.sh /usr/bin/s3mirror.sh

CMD ["siderite","task"]
```

This ensures that after a single run the container exits gracefully instead of waiting to timeout.

## Naming convention

Please name and publish your `hsdp_function` compatible Docker images using a repository name starting with `hsdp-function-...`.
This will help others identify the primary usage pattern for your image.
If your image represents a task, please use the prefix `hsdp-task-...`

## Gateway authentication

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

## Logging

The siderite helper supports direct logging to HSDP logging. You can either use API signing or an IAM service identity
which has the `LOG.CREATE` scope (recommended). Configuration can be done using below environment variables:

| environment | description | required |
|-------------|-------------|----------|
| `SIDERITE_LOGINGESTOR_PRODUCT_KEY`| The HSDP logging product key | Required |
| `SIDERITE_LOGINGESTOR_KEY` | The HSDP logging shared key | Optional |
| `SIDERITE_LOGINGESTOR_SECRET` | The HSDP logging shared secret | Optional |
| `SIDERITE_LOGINGESTOR_URL` | The HSDP logging base URL | Required when not setting region and environment |
| `SIDERITE_LOGINGESTOR_SERVICE_ID` | The HSDP service identity ID to use | Optional |
| `SIDERITE_LOGINGESTOR_SERVICE_PRIVATE_KEY` | The private key belonging to the service identity | Optional |
| `SIDERITE_LOGINGESTOR_REGION` | The HSDP region | Required for service identity |
| `SIDERITE_LOGINGESTOR_ENVIRONMENT` | The HSDP environment (`client-test`, `prod`) | Required for service identity |

### Logging using Logdrainer URL

If you have access to a Cloud foundry Logdrainer endpoint you can also leverage that for easy logging configuration. In that
case you only need to specify the Logdrainer endpoint URL:

| environment | description | required |
|-------------|-------------|----------|
| `SIDERITE_LOGDRAINER_URL` | The HSDP Logdrainer endpoint | Required |

Below is an example of using logging in a task:

```hcl
resource "hsdp_function" "request" {
  name = "http-request"
  docker_image = "philipslabs/hsdp-function-http-request:v0.6.0"

  environment = {
    REQUEST_METHOD   = "GET"
    REQUEST_URL      = "https://go-hello-world.eu-west.philips-healthsuite.com/dump"
    
    SIDERITE_LOGINGESTOR_REGION              = var.region
    SIDERITE_LOGINGESTOR_ENVIRONMENT         = var.environment
    SIDERITE_LOGINGESTOR_PRODUCT_KEY         = var.logging_product_key
    SIDERITE_LOGINGESTOR_SERVICE_ID          = hsdp_iam_service.logger.id
    SIDERITE_LOGINGESTOR_SERVICE_PRIVATE_KEY = hsdp_iam_service.logger.private_key
  }

  timeout = 30

  backend {
    credentials = module.siderite_backend.credentials
  }
}
```
