# hsdp_ai_inference_job

Manages HSDP IA Inference Jobs

## Example usage

```hcl
resource "hsdp_ai_inference_job" "job" {
  endpoint = data.hsdp_ai_inference_service_instance.inference.endpoint
  
  name          = "job1"
  description   = "Long running Inference Job"
  
  timeout       = 60
  
  model {
    reference = hsdp_ai_inference_model.model.reference
  }
  
  compute_target {
    reference = hsdp_ai_inference_compute_target.target.reference
  }
  
  input {
    name = "train"
    url  = "s3://input-sagemaker-64q6eey/data/input"
  }
  
  output {
    name = "train"
    url  = "s3://input-sagemaker-64q6eey/data/prediction"
  }
  
  environment = {
    FOO = "bar"
    BAR = "baz"
  }

  command_args = ["-f", "abc"]
  
  labels = ["BONEAGE", "CNN"]
}
```

## Argument reference

The following arguments are supported:

* `endpoint` - (Required) The AI Inference instance endpoint
* `name` - (Required) The name of Compute Environment
* `compute_target` - (Required) The compute Target to use
  * `reference` - (Required) The reference of the Compute Target
* `model` - (Optional) The model to use
  * `reference` - (Required) The reference of the Inference module
* `description` - (Optional) Description of the Compute Target
* `timeout` - (Optional) How long the job should run max.
* `input` - (Optional) Input data. Can have mulitple
  * `name` - (Required) Name of the input
  * `url` - (Required) URL pointing to the input
* `output` - (Optional) Output data. Can have mulitple
  * `name` - (Required) Name of the output
  * `url` - (Required) URL pointing to the output
* `environment` - (Optional, Map) Environment to set for Job
* `command_args` - (Optional, list(string)) Arguments to use for job

## Attributes reference

In addition to all arguments above, the following attributes are exported:
attributes are exported:

* `id` - The GUID of the job
* `reference` - The reference of this job
* `created` - The date this job was created
* `created_by` - Who created the environment
* `completed` - When the job was completed
* `duration` - How long (seconds) the job ran for
* `status` - The status of the job
* `status_message` - The status message, if available

## Import

An existing Compute Environment can be imported using `terraform import hsdp_ai_inference_compute_target`, e.g.

```bash
terraform import hsdp_ai_inference_compute_target.target a-guid
```
