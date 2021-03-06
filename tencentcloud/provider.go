package tencentcloud

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	PROVIDER_SECRET_ID  = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_REGION     = "TENCENTCLOUD_REGION"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_ID, nil),
				Description: "Secret ID of Tencent Cloud",
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_KEY, nil),
				Description: "Secret key of Tencent Cloud",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_REGION, nil),
				Description: "Region of Tencent Cloud",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"tencentcloud_autoscaling_groups":          dataSourceTencentCloudAutoScalingGroups(),
			"tencentcloud_availability_zones":          dataSourceTencentCloudAvailabilityZones(),
			"tencentcloud_container_cluster_instances": dataSourceTencentCloudContainerClusterInstances(),
			"tencentcloud_container_clusters":          dataSourceTencentCloudContainerClusters(),
			"tencentcloud_eip":                         dataSourceTencentCloudEip(),
			"tencentcloud_image":                       dataSourceTencentCloudSourceImages(),
			"tencentcloud_instance_types":              dataSourceInstanceTypes(),
			"tencentcloud_nats":                        dataSourceTencentCloudNats(),
			"tencentcloud_route_table":                 dataSourceTencentCloudRouteTable(),
			"tencentcloud_security_group":              dataSourceTencentCloudSecurityGroup(),
			"tencentcloud_subnet":                      dataSourceTencentCloudSubnet(),
			"tencentcloud_vpc":                         dataSourceTencentCloudVpc(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"tencentcloud_autoscaling_group":          resourceTencentCloudAutoScalingGroup(),
			"tencentcloud_alb_server_attachment":      resourceTencentCloudAlbServerAttachment(),
			"tencentcloud_cbs_snapshot":               resourceTencentCloudCbsSnapshot(),
			"tencentcloud_cbs_storage":                resourceTencentCloudCbsStorage(),
			"tencentcloud_cbs_storage_attachment":     resourceTencentCloudCbsStorageAttachment(),
			"tencentcloud_container_cluster":          resourceTencentCloudContainerCluster(),
			"tencentcloud_container_cluster_instance": resourceTencentCloudContainerClusterInstance(),
			"tencentcloud_dnat":                       resourceTencentCloudDnat(),
			"tencentcloud_eip":                        resourceTencentCloudEip(),
			"tencentcloud_eip_association":            resourceTencentCloudEipAssociation(),
			"tencentcloud_instance":                   resourceTencentCloudInstance(),
			"tencentcloud_key_pair":                   resourceTencentCloudKeyPair(),
			"tencentcloud_nat_gateway":                resourceTencentCloudNatGateway(),
			"tencentcloud_route_entry":                resourceTencentCloudRouteEntry(),
			"tencentcloud_route_table":                resourceTencentCloudRouteTable(),
			"tencentcloud_security_group":             resourceTencentCloudSecurityGroup(),
			"tencentcloud_security_group_rule":        resourceTencentCloudSecurityGroupRule(),
			"tencentcloud_subnet":                     resourceTencentCloudSubnet(),
			"tencentcloud_vpc":                        resourceTencentCloudVpc(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	secretId, ok := d.GetOk("secret_id")
	if !ok {
		secretId = os.Getenv(PROVIDER_SECRET_ID)
	}
	secretKey, ok := d.GetOk("secret_key")
	if !ok {
		secretKey = os.Getenv(PROVIDER_SECRET_KEY)
	}
	region, ok := d.GetOk("region")
	if !ok {
		region = os.Getenv(PROVIDER_REGION)
		if region == "" {
			region = "ap-guangzhou"
		}
	}
	config := Config{
		SecretId:  secretId.(string),
		SecretKey: secretKey.(string),
		Region:    region.(string),
	}
	return config.Client()
}
