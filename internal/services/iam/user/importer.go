package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/philips-software/terraform-provider-hsdp/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importUserContext(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()

	if strings.HasPrefix(importId, "login/") {
		loginID := strings.TrimPrefix(importId, "login/")
		c := m.(*config.Config)
		client, err := c.IAMClient()
		if err != nil {
			return nil, fmt.Errorf("IAMClient error: %w", err)
		}
		user, _, err := client.Users.LegacyGetUserIDByLoginID(loginID)
		if err != nil {
			return nil, fmt.Errorf("GetUserIDByLoginID error: %w", err)
		}
		importId = user
	}
	d.SetId(importId)
	return []*schema.ResourceData{d}, nil
}
