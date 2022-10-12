resource "hsdp_iam_role" "tdr_contract" {
    name = "TDRCONTRACT"
    managing_organization = var.tdr_org_id
    permissions = [
        "CONTRACT.READ", 
        "CONTRACT.CREATE"
    ]
}

resource "hsdp_iam_group" "tdr_contractadmin" {
    name = "TDR Contract Admin Group"
    managing_organization = var.tdr_org_id
    description           = "Group for TDR Users with Contract roles"
    roles                 = [hsdp_iam_role.tdr_contract.id]
    users                 = []
}
