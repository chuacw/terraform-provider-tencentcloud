package tencentcloud

import (
	"testing"
)

func Test_dataSourceTencentCloudSecurityGroupRuleRead(t *testing.T) {
	RequiredVars := map[string]string{
		"Region": "ap-singapore",
	}
	GetRequiredEnvVars(RequiredVars)
	config := Config{
		SecretId:  RequiredVars["SecretId"],
		SecretKey: RequiredVars["SecretKey"],
		Region:    RequiredVars["Region"],
	}
	resource := dataSourceTencentCloudSecurityGroupRule()
	d := resource.Data(nil)
	m, _ := config.Client()

	// The key used must be defined in the schema, otherwise, a panic will occur when d.Get is called on
	// the key.
	d.Set("sgId", "sg-elhg6l30")

	dataSourceTencentCloudSecurityGroupRuleRead(d, m)
}
