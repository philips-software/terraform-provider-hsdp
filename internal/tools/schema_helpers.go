package tools

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func StringSchema() *schema.Schema {
	return &schema.Schema{Type: schema.TypeString}
}

func SchemaSetStrings(ss []string) *schema.Set {
	s := &schema.Set{F: schema.HashSchema(StringSchema())}
	for _, str := range ss {
		s.Add(str)
	}
	return s
}

func IntSchema() *schema.Schema {
	return &schema.Schema{Type: schema.TypeInt}
}

func SchemaSetInts(l []int) *schema.Set {
	s := &schema.Set{F: schema.HashSchema(IntSchema())}
	for _, i := range l {
		s.Add(i)
	}
	return s
}
