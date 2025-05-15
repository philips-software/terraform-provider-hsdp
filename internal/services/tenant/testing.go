package tenant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestGenerateAPIKeyAndSignature is an exported version of generateAPIKeyAndSignature for testing
func TestGenerateAPIKeyAndSignature(d *schema.ResourceData) (string, string, error) {
	return generateAPIKeyAndSignature(d)
}
