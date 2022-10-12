resource "hsdp_iam_role" "KIBANALOGVIEWERS" {
  name                  = "KIBANALOGVIEWERS"
  description           = "Role to view HSDP Logs from Kibana"
  permissions           = ["LOG.READ"]
  managing_organization = hsdp_iam_org.testdev.id
}
