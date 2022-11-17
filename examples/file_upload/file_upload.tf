terraform {
  required_providers {
    hsdp = {
      source = "registry.terraform.io/philips-software/hsdp"
    }
  }
}

provider "hsdp" {
  iam_url            = var.iam_url
  idm_url            = var.iam_url
  org_admin_username = var.org_admin_username
  org_admin_password = var.org_admin_password
  oauth2_client_id   = var.oauth2_client_id
  oauth2_password    = var.oauth2_password
}


data "archive_file" "config_template_folder" {
  type        = "zip"
  source_dir  = "C:/path"
  output_path = "./files/file_or_dir.zip"
}

resource "hsdp_file_upload" "name" {
  url       = "http://localhost:8989/reporting/7d32cc51-f715-49eb-a61f-15d6ecb6b277/templateconfig"
  file_path = "./examples/file_upload/files/file_or_dir.zip"
  checksum  = data.archive_file.config_template_folder.output_sha
}
