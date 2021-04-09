package hsdp

import (
	"testing"

	"github.com/google/fhir/go/jsonformat"
	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
	"github.com/stretchr/testify/assert"
)

func TestParameters(t *testing.T) {
	body := `{
    "resourceType": "Parameters",
    "parameter": [
        {
            "name": "status",
            "valueString": "SUCCESS"
        },
        {
            "name": "submissionTime",
            "valueDateTime": "2018-07-13T00:05:29.981681+00:00"
        },
        {
            "name": "lastUpdated",
            "valueDateTime": "2018-07-13T00:08:37.563+00:00"
        },
        {
            "name": "requestor",
            "valueString": "77d6e95d-6f2a-4739-9d9c-bfa52f39a3e9"
        }
    ]
}`
	um, err := jsonformat.NewUnmarshaller("UTC", jsonformat.STU3)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, um) {
		return
	}
	unmarshalled, err := um.Unmarshal([]byte(body))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, unmarshalled) {
		return
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	params := contained.GetParameters()

	assert.Len(t, params.Parameter, 4)
	assert.Equal(t, "status", params.Parameter[0].Name.Value)
	assert.Equal(t, "SUCCESS", params.Parameter[0].Value.GetStringValue().Value)
}
