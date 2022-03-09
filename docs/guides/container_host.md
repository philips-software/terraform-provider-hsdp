---
page_title: "Working with Container Host"
---
# Working with Container Host

The Cartel `container-host` role provides a shared responsibility and security-enhanced host for running Docker containers within HSDP.

## Cartel vs Container Host

Cartel is a simple REST like JSON API which provides orchestration and provisioning capabilities to the environment.
Cartel is designed to deploy the managed services over generic server instances. One of the key arguments of the Cartel
Create request is the type of role to assign to the Cartel instance (really, an EC2 instance). While there are quite a few roles defined the only
supported role today is `container-host`. This is why Cartel and Container Host are effectively interchangeable in conversations.

-> Takeaway: The Cartel REST API is used to provision Container-Host instances

## The Container Host role

The `container-host` role is not intended to replace Cloud Foundry for application deployment, but as a complement to Cloud Foundry.
It can be used to deploy services not supported by HSDP Service Brokers, custom applications that do not work well on Cloud Foundry, or for administrative purposes.

-> In the remainder of this guide when referring to `Container Host` we mean a Cartel instance provisioned with the `container-host` role

## When to use Container Host

There are some common patterns when Container Host makes sense. We'll attempt to list
a few of them here. Please reach out to your Technical Account Manager in case you have
further questions. Common reasons include:

- Legacy applications that rely on persistent disk files
- Monolithic applications with extreme memory requirements (`> 20GB`)
- Applications which require dedicated CPU (compute) resources
- Off-The-Shelf software which does not fit well on Cloud foundry

-> Rule of thumb: treat Container Host as an escape hatch (last resort) i.e. when your workload absolutely cannot be accommodated on Cloud foundry

## When not to use Container Host

- If your app works fine on Cloud foundry
- If you need to run a `Windows` Legacy application (Container Host is `Linux only!`)
- If your application is not dockerized (Dockerize it first)

### Alternatives to Container Host

If your workload does not fit on either Cloud foundry or Container Host then HSP also
offers [HSDP Hosting Services](https://www.hsdp.io/documentation/hsdp-hosting-services). This is a more traditional hosting environment where you
will have to specify the resources up front (limited elasticity) and will need to provide a runbook for operations. This is
typically the landing zone for lift-and-shift operations and is in general not self-serviceable. Use this service if your
workload is not Cloud Native and you need time for refactoring.

## Minimum requirements for using Container Host

You'll need to request one-time onboarding onto the Cartel service
via a `ServiceNow` ticket. Please contact your Technical Account Manager for
assistance with this if needed.

### Cartel keys

As part of the onboarding procedure you will receive a set of credentials via SDT (Secure Digital Transfer)
consisting of a `Cartel Key` and a `Cartel Secret`. These will be used to sign the Cartel API calls e.g. to
provision instances or modify running instances.

### LDAP account

An HSP LDAP account is required to interact with Container Host instances. This is provided
as part of your onboarding onto the HSP platform.

### SSH private key

The `SSH public key` which is associate to your `LDAP` account is used to authenticate and gain
access to Container Host instances. The associated `SSH private key` should be
loaded in your local SSH agent and agent forwarding should be enabled (recommended). Alternatively, you
can pass the private key as an argument to the HSP Terraform provider configuration

## Provisioning your first instance

With all prerequisites in place we can now talk about provisioning your
first Container Host instance. We'll briefly illustrate how this can be done using the Cartel CLI
but focus of this guide will be on leveraging Terraform.

### Using the CLI

The Cartel CLI is installed on the [HSP SSH regional jump gates](https://www.hsdp.io/develop/get-started-healthsuite/set-up-ssh-access/access-services-behind-ssh-gateway/connect-to-gateway).

Listing all your instances:

```shell
cartel --get-all-instances -y

[]
```

Get all available roles:

```shell
cartel --get-all-roles -y
```

```shell
- description: Cartel container hosting.
  role: container-host
```

Create a new instance:

```shell
cartel --create \
  -r container-host \
  -n ron.dev \
  --ldap-groups rswanson
```

```shell
message:
- eip_address: null
  instance_id: i-0b4c962a374e47456
  ip_address: 127.0.44.112
  name: ron.dev
  role: container-host
result: Success
```

For more information on other CLI arguments:

```shell
cartel --help
```

## Using Terraform

The HSP Terraform provider supports managing Container Host instances
through the [hsdp_container_host](https://registry.terraform.io/providers/philips-software/hsdp/latest/docs/resources/container_host) and [hsdp_container_exec](https://registry.terraform.io/providers/philips-software/hsdp/latest/docs/resources/container_host_exec) resource types.

Use `hsdp_container_host` to declare an instance. Example:

```hcl
resource "hsdp_container_host" "server" {
  name             = "pawnee-server.dev"

  user             = "cdrummer"
  user_groups      = ["cdrummer"]

  instance_type    = "m5.8xlarge"
  security_groups  = ["analytics"]
  subnet_type      = "public"
  
  volumes          = 1
  volume_size      = 100
  
  agent            = true

  tags = {
    created_by = "rswanson"
  }
}
```

## Configuration options

There are many configuration options and we'll discuss the common arguments used for
Container Host resource definitions in Terraform, starting with the only required attribute, `name`

The name should be a `DNS` friendly string post-fixed with `.dev`, example:

```hcl
resource "hsdp_container_host" "tynan" {
  name = "tynan-server.dev"
}
```

Make sure the name is unique in the region you are deploying in. The best practice is to introduce
some randomness in the name to ensure this using for example a `random_pet` resource:

```hcl
resource "random_pet" "deploy" {
}

resource "hsdp_container_host" "server" {
  name = "server-${random_pet.deploy.id}.dev"
}
```

### User and user groups

You'll want to assign one or more `user groups` to your instance. The assigned user groups determine
which LDAP accounts have access to the instance via SSH. The best practice is to at least add
your own LDAP group which has the name same as your LDAP `login`. So, a user with LDAP `cdrummer` would
add the following arguments:

```hcl
resource "hsdp_container_host" "tynan" {
  name = "tynan-server.dev"

  # Assign user and user_groups
  user        = "cdrummer"
  user_groups = ["cdrummer"]
}
```

### Instance types

You can choose which underlying EC2 instance type to use for your Container Host.
Most EC2 instance types are supported with some exceptions, notably `i3` family instances types or `GPU`
compute instances are **not supported**. The default instance type is `m5.large`

The below argument would allocate an EC2 instance with `32 VCPUs` and `128GB RAM`

```hcl
resource "hsdp_container_host" "tynan" {
  name        = "tynan-server.dev"
  user        = "cdrummer"
  user_groups = ["cdrummer"]

  # Override the default with a more powerful instance type
  instance_type = "m5.8xlarge"
}
```

### Security groups

Security groups determine which ports are opened up on the Container Host instance.
Hosting any services requires assigning one or more relevant `security groups`. To get
a list of available security groups:

```hcl
data "hsdp_container_host_security_groups" "all" {
}

output "all_security_groups" {
  value = data.hsdp_container_host_security_groups.all.names
}
```

You can also request creation of custom security groups via ServiceNOW. Requests will go through
an approval process so might take a while, so it's best to check weather a
pre-approved security group will fit your needs.

Once you find a suitable security group you can assign it:

```hcl
resource "hsdp_container_host" "tynan" {
  name          = "tynan-server.dev"
  user          = "cdrummer"
  user_groups   = ["cdrummer"]
  instance_type = "m5.8xlarge"

  # Assign a few security groups
  security_groups  = ["analytics", "http-from-cloudfoundry"]
}

```

### Subnet types

There are two subnet types available:

- `public` - Your instance will be assigned a public IP address as well as a private one.
- `private` - Your instance will be assigned a private IP only

You can instruct Terraform to pick a specific subnet by just specifying
the type only:

```hcl
resource "hsdp_container_host" "tynan" {
  name             = "tynan-server.dev"
  user             = "cdrummer"
  user_groups      = ["cdrummer"]
  instance_type    = "m5.8xlarge"
  security_groups  = ["analytics", "http-from-cloudfoundry"]
  
  # Provision this instance in a public subnet, we don't care which one
  subnet_type = "public"
}
```

The assigned subnet name will be followed by a letter in the name.
The letter corresponds to the availability zone that the instance is deployed in.

You also have the option to specify the exact subnet and availability zone you want to
deploy your instance to:

```hcl
resource "hsdp_container_host" "tynan" {
  name             = "tynan-server.dev"
  user             = "cdrummer"
  user_groups      = ["cdrummer"]
  instance_type    = "m5.8xlarge"
  security_groups  = ["analytics", "http-from-cloudfoundry"]

  # Provision this instance in a public subnet in availability zone c
  subnet = "sc1-public-c"
}
```

~> Using `public` subnets increases your costs slightly as you will pay extra for the public IP address

### Volumes

Persisting data on Container Host requires attaching (EBS) volumes to your instance. Use the following attributes
to configure this:

- `volumes` - The number of EBS volumes to attach
- `volume_size` - The size (in `GB`) of each volume

```hcl
resource "hsdp_container_host" "tynan" {
  name             = "tynan-server.dev"
  user             = "cdrummer"
  user_groups      = ["cdrummer"]
  instance_type    = "m5.8xlarge"
  security_groups  = ["analytics", "http-from-cloudfoundry"]
  subnet_type      = "public"
  
  # Attach a 2TB EBS volume
  volumes     = 1
  volume_size = 2000
}
```

### SSH access

Once your Container Host instance is up and running you will want to login and start deploying containers. Assuming
your private key is loaded in your local SSH agent you can access your instance using the [regional SSH jump gate](https://www.hsdp.io/develop/get-started-healthsuite/set-up-ssh-access/access-services-behind-ssh-gateway/connect-to-gateway):

```shell
ssh -A -C -J cdrummer@gw-na1.phsdp.com cdrummer@tynan-server.dev
```

## Bootstrapping Container Host instances

To fully automate provisioning and deployment of software on Container Host you can leverage the `hsdp_container_host_exec` resource.
While it is possibly to specify `file` and `commands` in the `hsdp_container_host` resource directly we highly recommend splitting off
bootstrapping steps to a `hsdp_container_host_exec` resource. This decouples your software bootstrapping from the Container Host instance provisioning, which
can take between 7 and 10 minutes on average.

We'll break down the example below:

```hcl
resource "hsdp_container_host_exec" "nomad_server_init" {
  # Use the private_ip from a provisioned hsdp_container_host resource
  # This ensures the exec resource is only executed after the Container Host is up and running
  host  = hsdp_container_host.nomad_server.private_ip
  user  = var.ldap_user
  agent = true

  # Render template and copy it to the Container Host instance
  file {
    content = templatefile("${path.module}/templates/server.hcl", {
      advertise_ip   = hsdp_container_host.nomad_server.private_ip
      region         = "global"
      datacenter     = "dc1"
      name           = "server"
      docker_runtime = var.docker_runtime
    })
    destination = "/home/${var.ldap_user}/server.hcl"
    permissions = "0755"
  }

  # Execute a sequence of commands, using the copied files as input
  commands = [
    "docker stop nomad-server || true",
    "docker rm nomad-server || true",
    "docker rm alpine || true",
    "docker volume rm nomad-server-config || true",
    "docker volume create nomad-server-config",
    "docker create -v nomad-server-config:/config --name alpine alpine",
    "docker cp /home/${var.ldap_user}/server.hcl alpine:/config",
    "docker run -d --restart on-failure --name nomad-server -v nomad-server-config:/config ${var.nomad_image} nomad agent -server -bind=0.0.0.0 -acl-enabled -plugin-dir=/plugins -config=/config/server.hcl -data-dir=/tmp/nomad",
    "sleep 5",
    "docker exec nomad-server nomad acl bootstrap"
  ]
}
```

### host

The `host` parameter should be set to the private IP address of the Container Host instance. The `hsdp_container_exec` will use
an embedded SSH agent which supports tunneling (and proxy traversal!) to establish a connection to your Container Host.

### file

Use `file {}` blocks to prepare and copy content from your Terraform deployment all the way to the Container Host instances

### commands

Finally, you can specify one or more `commands` which will be executed in sequence on the Container Host. You can use this to
create volumes, networks and spin up containers. The output of the last command will exported under the `result` attribute. This allows
you to capture data / state from the Container Host and make it available to other Terraform resources.

## Forwarding Internet traffic to Container Host

Ingress traffic to the platform needs to flow through the Cloud foundry Load balancers. A reverse proxy can be deployed
on CF which then forwards traffic to your Container Host. Make sure the proper `security groups` are in place to allow
traffic from the Cloud foundry VPC (examples: `http-from-cloud-foundry`, `http-8080`)

```shell
Internet --> Proxy on CF --> Container Host instance
```

Below is an example on how to use [Caddy](https://caddyserver.com/) as a reverse proxy for your Container Host:

```hcl
resource "cloudfoundry_app" "superset_proxy" {
  name         = "caddy"
  space        = data.cloudfoundry_space.space.id
  memory       = 128
  disk_quota   = 512

  # Use Caddy from Docker Hub, optionally pull from HSDP Docker Registry to prevent limits 
  docker_image = "caddy:2.4.6"
  docker_credentials = {
    username = var.docker_username
    password = var.docker_password
  }

  # Render the template and inject the Container Host private IP address
  environment = merge({
    CADDYFILE_BASE64 = base64encode(templatefile("${path.module}/templates/Caddyfile", {
      upstream_url = "http://${hsdp_container_host.server.private_ip}:8080"
    }))
  }, {})

  # Override the Docker command to render the config from ENV in the container and let Caddy use it
  command           = "echo $CADDYFILE_BASE64 | base64 -d > /etc/caddy/Caddyfile && cat /etc/caddy/Caddyfile && caddy run -config /etc/caddy/Caddyfile"
  health_check_type = "process"

  routes {
    route = cloudfoundry_route.proxy.id
  }
}

resource "cloudfoundry_route" "proxy" {
  domain   = data.cloudfoundry_domain.domain.id
  space    = data.cloudfoundry_space.space.id
  hostname = "caddy-proxy"
}
```

The Caddy template in `templates/Caddyfile` would look like:

```caddy
:80 {
  @insecure {
    header X-Forwarded-Proto http
  }
  redir @insecure https://{host}{uri} permanent

  reverse_proxy ${upstream_url}
}
```

-> The above resource uses the [Cloud foundry Terraform provider](https://registry.terraform.io/providers/cloudfoundry-community/cloudfoundry/0.15.0)
