resource "hsdp_iam_group" "log_viewers" {
  managing_organization = "${hsdp_iam_org.testdev.id}"
  name                  = "Log Viewers"
  description           = "Group for Log viewers"
  roles                 = ["${hsdp_iam_role.KIBANALOGVIEWERS.id}"]
  users                 = ["${hsdp_iam_user.developer.id}"]
}
