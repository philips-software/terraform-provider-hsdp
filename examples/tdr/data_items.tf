resource "hsdp_iam_role" "tdr_data_items" {
    name = "TDRDATAITEMS"
    managing_organization = "${var.tdr_org_id}"
    permissions = [
        "DATAITEM.CREATE", 
        "DATAITEM.READ",
        "DATAITEM.DELETE",
        "DATAITEM.PATCH"
    ]
}

resource "hsdp_iam_group" "tdr_data_items" {
    name = "TDR Data Items Admin Group"
    managing_organization = "${var.tdr_org_id}"
    description           = "Group for TDR Users with Dataitem roles"
    roles                 = ["${hsdp_iam_role.tdr_data_items.id}"]
    users                 = []
}
