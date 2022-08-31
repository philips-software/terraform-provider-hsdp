package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importRepositoryContext(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("expecting config_url,organization_id,repository_id as import string")
	}
	configURL := parts[0]
	orgID := parts[1]
	id := parts[2]

	d.SetId(id)
	_ = d.Set("config_url", configURL)
	_ = d.Set("organization_id", orgID)
	return []*schema.ResourceData{d}, nil
}
