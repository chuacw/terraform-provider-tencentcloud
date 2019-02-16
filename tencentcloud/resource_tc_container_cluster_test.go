package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	ccs "github.com/zqfan/tencentcloud-sdk-go/services/ccs/unversioned"
)

func TestAccTencentCloudContainerCluster_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudContainerClusterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_container_cluster.foo"),
					checkContainerClusterInstancesAllNormal("tencentcloud_container_cluster.foo"),
				),
			},
		},
	})
}

// For ordinary usage, it doesn't require all nodes in a cluster to be in normal state.
// But for acceptance test, it only has a single node and should be in normal state otherwise
// will cause resource leak such as vpc, subnet and vm resources, these leakage will block
// subsequential acceptance test, hence here we need to do such check to ensure cluster node
// is in an expected state.
func checkContainerClusterInstancesAllNormal(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Container cluster ID is not set")
		}

		conn := testAccProvider.Meta().(*TencentCloudClient).ccsConn
		req := ccs.NewDescribeClusterRequest()
		req.ClusterIds = []*string{&rs.Primary.ID}
		// For now, cluster instance will be reinstalled, hence it needs to wait more time
		err := resource.Retry(20*time.Minute, func() *resource.RetryError {
			resp, err := conn.DescribeCluster(req)
			if err != nil {
				return resource.RetryableError(err)
			}
			if _, ok := err.(*common.APIError); ok {
				return resource.NonRetryableError(err)
			}
			if *resp.Data.Clusters[0].NodeStatus == "AllNormal" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Cluster node status is %s", *resp.Data.Clusters[0].NodeStatus))
		})
		return err
	}
}

const testAccTencentCloudContainerClusterConfig_basic = `
resource "tencentcloud_vpc" "my_vpc" {
  cidr_block = "10.6.0.0/16"
  name       = "terraform_vpc_test"
}

resource "tencentcloud_subnet" "my_subnet" {
  vpc_id = "${tencentcloud_vpc.my_vpc.id}"
  availability_zone = "ap-guangzhou-3"
  name              = "terraform_test_subnet"
  cidr_block        = "10.6.0.0/24"
}

resource "tencentcloud_container_cluster" "foo" {
 cluster_name = "terraform-acc-test"
 cpu    = 1
 mem    = 1
 os_name   = "ubuntu16.04.1 LTSx86_64"
 bandwidth  = 1
 bandwidth_type = "PayByHour"
 require_wan_ip   = 1
 subnet_id  = "${tencentcloud_subnet.my_subnet.id}"
 is_vpc_gateway = 0
 storage_size = 0
 root_size = 100
 root_type = "CLOUD_SSD"
 goods_num  = 1
 password  = "Admin12345678"
 vpc_id   = "${tencentcloud_vpc.my_vpc.id}"
 cluster_cidr = "10.0.0.0/19"
 cvm_type  = "PayByHour"
 cluster_desc = "foofoofoo"
 period   = 1
 zone_id   = 100003
 instance_type = "S2.SMALL1"
 mount_target = ""
 docker_graph_path = ""
 instance_name = "terraform-container-acc-test-vm"
 cluster_version = "1.7.8"
}
`
