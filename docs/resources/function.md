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

## Attribute reference
