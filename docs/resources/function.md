# hsdp_function

Define function-as-a-service using various backends. Currently
only `iron` is supported. 

## Example usage

```hcl
resource "hsdp_function" "rds_backup" {
  name = "streaming-backup"
  
  docker_image = var.streaming_backup_image
  docker_credentails = {
    username = var.docker_username
    password = var.docker_password
  }
  
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

  schedule {
    start = "2021-01-01T04:00Z"
    run_every = "1d"
  }

  backend {
    type = "iron"
    credentials = module.iron_backend.credentials
  }  
}
```

## Argument reference
The following arguments are supported:

* `name` - (Required) The name of the function
* `docker_image` - (Required) The docker image that contains the logic for the function
* `docker_credentials` - (Optionals) The docker registry credentials
  * `username` - (Required) The registry username
  * `password` - (Required) The registry password  
* `environment` - (Optional, map) The environment variables to set in the docker container before executing the function
* `schedule` - (Optional) Schedule the function. When not set, the function is just a task.
  * `start` - (Required, RFC3339) When to start the schedule
  * `run_every` - (Required) Run the function every `{value}{unit}` period. Supported units are `s`, `m`, `h`, `d` for second, minute, hours, days respectively.
    Example: a value of `"20m"` would run the function every 20 minutes.
* `backend` - (Required) The backend to use for scheduling your functions.
  * `type` - (Required) The backend type. Only `iron` is supported at this time.
  * `credentials` - (Required) The backend credentials. Must be iron configuration details at this time.
    
## Attribute reference
