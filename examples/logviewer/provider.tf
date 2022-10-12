variable "iam_url" {}
variable "idm_url" {}
variable "oauth2_client_id" {}
variable "oauth2_password" {}
variable "org_id" {}
variable "org_admin_username" {}
variable "org_admin_password" {}
variable "shared_key" {}
variable "secret_key" {}

provider "hsdp" {
  iam_url            = var.iam_url
  idm_url            = var.idm_url
  oauth2_client_id   = var.oauth2_client_id
  oauth2_password    = var.oauth2_password
  org_id             = var.org_id
  org_admin_username = var.org_admin_username
  org_admin_password = var.org_admin_password
  shared_key         = var.shared_key
  secret_key         = var.secret_key
}
