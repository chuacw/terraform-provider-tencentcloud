package tencentcloud

import (
	"testing"
)

func Test_dataSourceTencentCloudAutoScalingGroupsRead(t *testing.T) {
	RequiredVars := map[string]string{
		"Region": "ap-beijing",
	}
	GetRequiredEnvVars(RequiredVars)
	config := Config{
		SecretId:  RequiredVars["SecretId"],
		SecretKey: RequiredVars["SecretKey"],
		Region:    RequiredVars["Region"],
	}
	resource := dataSourceTencentCloudAutoScalingGroups()
	d := resource.Data(nil)
	m, _ := config.Client()
	// d.Set("scalingConfigurationId", "asg-4gpnrf3l")
	dataSourceTencentCloudAutoScalingGroupsRead(d, m)
}
