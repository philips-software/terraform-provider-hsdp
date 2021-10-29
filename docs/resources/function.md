---
subcategory: "Functions"
---

# hsdp_function

Define function-as-a-service using various backends. Currently,
only `siderite` (HSDP Iron) is supported.

## Example usage

```hcl
resource "hsdp_function" "rds_backup" {
  name = "streaming-backup"
  
  # The docker packaged function business logic
  docker_image = var.streaming_backup_image
  docker_credentials = {
    username = var.docker_username
    password = var.docker_password
  }
  
  # Environment variables available in the container
  environment = {
    db_name = "hsdp_pg"
    db_host = "rds.aws.com"
    db_username = var.db_username
    db_password = var.db_password
    
    s3_access_key = "AAA"
    s3_secret_key = "BBB"
    s3_bucket = "cf-s3-xxx"
    s3_prefix = "/backups"
  }

  # Run every day at 4am
  schedule = "0 4 * * *"

  backend {
    credentials = module.siderite_backend.credentials
  }  
}
```

## Argument reference

The following arguments are supported:

* `name` - (Required) The name of the function
* `docker_image` - (Required) The docker image that contains the logic of the function
* `docker_credentials` - (Optional) The docker registry credentials
  * `username` - (Required) The registry username
  * `password` - (Required) The registry password  
* `command` - (Optional) The command to execute in the container. Default is `/app/server`
* `environment` - (Optional, map) The environment variables to set in the docker container before executing the function
* `schedule` - (Optional) set schedule using cron format. This requires a backend with activated gateway. Conflicts with `run_every`
* `run_every` - (Optional) Run the function every `{value}{unit}` period. Supported units are `s`, `m`, `h`, `d` for second, minute, hours, days respectively. Conflicts with `cron`
  Example: a value of `"20m"` would run the function every 20 minutes.
* `start_at` - (Optional) Only valid for `run_every`. This is a hint for when the first run should be.
  This determines the time of day the schedule will run at. Depending on the frequency of the runs and
  the time of day when the Terraform script was run, it can take up to 24 hours for the first run to happen.
  Use `schedule` for more accurate scheduling behaviour.
* `timeout` - (Optional, int) When set, limits the execution time (seconds) to this value. Default: `1800` (30 minutes)
* `backend` - (Required) The backend to use for scheduling your functions.
  * `credentials` - (Required, map) The backend credentials. Must be iron configuration details at this time.

## Attribute reference

The following attributes are exported:

* `endpoint` - The gateway endpoint where you can trigger this function
* `async_endpoint` - The gateway endpoint where you can schedule the function asynchronously  
* `token` - The token to use in case `auth_type` is set to `token`. This token must be pasted in the HTTP `Authorization` header as `Token TOKENHERE`  
* `auth_type` - The authentication type. Possible values [`none`, `token`, `iam`]
