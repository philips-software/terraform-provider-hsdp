---
subcategory: "Configuration"
---

# hsdp_config

Retrieve configuration details from services based on region and environment

## Example Usage

```hcl
data "hsdp_config" "iam_us_east_prod" { 
  service = "iam"
  region = "us-east"
  environment = "prod"
}
```

```hcl
output "iam_url_us_east_prod" {
   value = data.hsdp_config.iam_us_east_prod.url
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) The HSDP service to lookup

Availability of services varies across regions. The following services are discoverable:

| Service         | Description                                                                                                                                                 |
|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| cartel          | The Cartel API service. Manages [Container Host](https://www.hsdp.io/documentation/container-host) instances                                                |
| cf              | HSDP [Cloud foundry](https://www.hsdp.io/develop/architecture/cloud-foundry) regional PaaS configuration                                                    |
| console         | HSDP [Console](https://www.hsdp.io/documentation/metrics-service-broker/service-details) API endpoints                                                      |
| docker-registry | Regional [Docker Registry](https://www.hsdp.io/documentation/docker-registry) details                                                                       |
| gateway         | Regional [SSH gateway](https://www.hsdp.io/develop/get-started-healthsuite/set-up-ssh-access/access-services-behind-ssh-gateway/connect-to-gateway) details |
| has             | Hosted Application Streaming [HAS](https://www.hsdp.io/documentation/hosted-application-streaming/getting-started-with-hosted-application-streaming#)       |
| iam             | Identity and Access Management [IAM](https://www.hsdp.io/documentation/identity-and-access-management-iam)                                                  |
| idm             | Identity and User Management. Part of [IAM](https://www.hsdp.io/documentation/identity-and-access-management-iam)                                           |
| kibana          | Kibana endpoint. Part of [HSDP Logging](https://www.hsdp.io/documentation/logging)                                                                          |
| logging         | [HSDP Logging](https://www.hsdp.io/documentation/logging) API details                                                                                       |
| logquery        | Log query endpoint details. Part of [HSDP Logging](https://www.hsdp.io/documentation/logging)                                                               |
| mdm             | Master Data Management [MDM](https://www.hsdp.io/documentation/master-data-management)                                                                      |
| notification    | HSDP [Notification service](https://www.hsdp.io/documentation/notification)                                                                                 |
| pki             | Public Key Infrastructure [PKI](https://www.hsdp.io/documentation/public-key-infrastructure/getting-started) services                                       |
| s3creds         | [S3 Credentials](https://www.hsdp.io/documentation/s3-credentials) API details                                                                              |
| edge            | Edge / STL API details                                                                                                                                      |
| uaa             | User Account and Authentication [UAA](https://docs.cloudfoundry.org/concepts/architecture/uaa.html). Part of Cloud foundry                                  |
| vault-proxy     | Vault proxy details. Part of [Vault Service Broker](https://www.hsdp.io/documentation/vault-service-broker/service-details)                                 |

* `region` - (Optional) The HSDP region. If not set, defaults to provider level config

The following regions are recognized:

| Region  | Description                                                              |
|---------|--------------------------------------------------------------------------|
| apac2   | [Japan](https://en.wikipedia.org/wiki/Japan) (Tokyo)                     |
| apac3   | [Asia-Pacific](https://en.wikipedia.org/wiki/Asia-Pacific) (Sydney)      |
| ca1     | [Canada](https://en.wikipedia.org/wiki/Canada) (Central Canada)          |
| cn1     | [China](https://en.wikipedia.org/wiki/China) (Beijing)                   |
| dev     | Development (US)                                                         |
| eu-west | [European Union](https://en.wikipedia.org/wiki/European_Union) (Ireland) |
| sa1     | [South America](https://en.wikipedia.org/wiki/South_America) (Sao Paulo) |
| us-east | [United States](https://en.wikipedia.org/wiki/United_States) (Virginia)  |

* `environment` - (Optional) The HSDP environment. If not set, defaults to provider level config

Environments vary across regions. The following environemnts are valid

| Environment | Description                           |
|-------------|---------------------------------------|
| dev         | Development. Only in region `us-east` |
| client-test | Client Test / Development environment |
| prod        | Production                            |

## Attributes Reference

The following attributes are exported:

* `url` - (string) The (base / API) URL of the service
* `host` - (string) The host of the service
* `domain` - (string) The domain associated with the service
* `regions` - (string) The list of known regions
* `services` - (list(string)) The list of available services in the region/environment
* `service_id` - (string) The IAM service ID used for authenticating against IAM
* `org_admin_username` - (string) The IAM OrgAdmin used for authenticating against IAM
* `sliding_expires_on` - (string) A sliding expires on RFC3339 timestamp which can be used to rotate e.g. credentials.
  The value is the first day of the next quarter calculated from the current wall clock time.
