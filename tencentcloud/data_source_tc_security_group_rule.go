package tencentcloud

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"sgId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceTencentCloudSecurityGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":    "DescribeSecurityGroupPolicys",
		"projectId": "0",
		"sgId":      d.Get("sgId").(string),
	}

	log.Printf("[DEBUG] data_source_tc_security_group_rule read params: %v", params)

	response, err := client.SendRequest("dfw", params)
	if err != nil {
		log.Printf("[ERROR] data_source_tc_security_group_rule read client.SendRequest error: %v", err)
		return err
	}

	var jsonresp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TotalNum int `json:"totalNum"`
			Detail   []struct {
				SgId             string `json:"sgId"`
				SgName           string `json:"sgName"`
				SgRemark         string `json:"sgRemark"`
				BeAssociateCount int    `json:"beAssociateCount"`
				CreateTime       string `json:"createTime"`
			}
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group read json.Unmarshal error:%v", err)
		return err
	}
	if jsonresp.Code != 0 {
		log.Printf("[ERROR] resource_tc_security_group read error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		return errors.New(jsonresp.Message)
	} else if jsonresp.Data.TotalNum <= 0 || len(jsonresp.Data.Detail) <= 0 {
		return errors.New("Security group not found")
	}

	sg := jsonresp.Data.Detail[0]

	d.SetId(sg.SgId)
	d.Set("security_group_id", sg.SgId)
	d.Set("name", sg.SgName)
	d.Set("description", sg.SgRemark)
	d.Set("create_time", sg.CreateTime)
	d.Set("be_associate_count", sg.BeAssociateCount)

	return nil
}
