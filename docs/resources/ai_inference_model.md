# hsdp_ai_inference_model

Manages HSDP AI Inference models.

## Example usage
```hcl
data "hsdp_config" "inference" {
  service = "inference"
}

data "hsdp_ai_inference_instance" "inference" {
  base_url        = data.hsdp_config.inference.url
  organization_id = var.inference_tenant_org_id
}

resource "hsdp_ai_inference_model" "model" {
  endpoint = data.hsdp_ai_inference_instance.inference.endpoint
  
  name          = "model1"
  version       = "v1"
  description   = "Test model"
 
  compute_environment  {
    reference  = "foo"
    identifier = "bar"
  }
  
  source_code  {
    url       = "git@github.com:testuser/source.git"
    branch    = "main"
    commit_id = "e1f9366"
    ssh_key   = "..."
  }

  artifact_path = "git@github.com:testuser/example.git"

  entry_commands = ["python main/train.py -s 134786"]
  
  environment = {
    FOO = "bar"
    BAR = "baz"
  }
  
  labels = ["CNN"]

  additional_configuration =  "{\"Tags\": [ { \"Key\": \"name\",\"Value\": \"hsp\"}]}"
}
```

## Argument reference
The following arguments are supported:

* `endpoint` - (Required) The AI Inference instance endpoint
* `name` - (Required) The name of the Model
* `compute_environment` - (Required)
* `source_code` - (Required)
* `entry_commands` - (Required, list(string)) Commands to execute
* `environment` - (Optional, Map) List of environment variables to set
* `labels` - (Optional)
* `artifact_path` - (Optional)
* `description` - (Optional) Description of the Compute Target
* `additional_configuration` - (Optional)

## Attribute reference
In addition to all arguments above, the following attributes are exported:
attributes are exported:

* `id` - The GUID of the Model
* `reference` - The reference of this Model 
* `created` - The date this Model  was created
* `created_by` - Who created the Model
