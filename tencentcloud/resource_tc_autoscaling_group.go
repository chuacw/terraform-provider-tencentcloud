package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

var (
	DEBUG bool
)

func init() {
	DEBUG = true
}

type (
	scalingGroupId struct {
		scalingGroupId   string
		scalingGroupName string
	}

	securityGroup struct {
	}

	scalingConfiguration struct {
		scalingConfigurationId   string
		scalingConfigurationName string
		scalingGroupIdSet        []scalingGroupId
		cpu                      int
		mem                      int
		imageType                int
		imageId                  string
		storageType              int
		storageSize              int
		rootSize                 int
		bandwidthType            string
		bandwidth                int
		wanIp                    int
		keyId                    string
		password                 string
		sgSet                    []securityGroup
		needMonitorAgent         int
		needSecurityAgent        int
		createTime               string
		projectId                int
	}
)

func validateAutoScalingGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("Group Name should not be empty, current length: %d", len(value)))
	}
	return
}

func validateAutoScalingGroupPassword(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 8 || len(value) > 16 {
		errors = append(errors, fmt.Errorf("Password length should be between 8 and 16, current length: %d", len(value)))
	}
	pattern := `[A-Za-z]+`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("Password did not contain %s %s", "alphabets", "A-Za-z"))
	}
	pattern = `[0-9]+`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("Password did not contain %s", "numbers"))
	}
	pattern = `[\!\@\#\$\%\^\\\(\)]+`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("Password did not contain %s %s", "special characters", `!@#$%^\()`))
	}
	return
}

func resourceTencentCloudAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAutoScalingGroupCreate,
		Delete: resourceTencentCloudAutoScalingGroupDelete,
		Read:   resourceTencentCloudAutoScalingGroupRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scalingConfigurationName": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"imageId": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: nil,
			},
			"cpu": &schema.Schema{
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: nil,
			},
			"mem": &schema.Schema{
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"storageType": &schema.Schema{
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"bandwidthType": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"bandwidth": &schema.Schema{
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"imageType": &schema.Schema{
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"rootSize": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupPassword,
			},
			"needMonitorAgent": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"needSecurityAgent": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"wanIp": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"sgId": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"projectId": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"dataSnapshotId": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
			"cvmType": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAutoScalingGroupName,
			},
		},
	}
}

func resourceTencentCloudAutoScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Version":                  "2017-03-12",
		"ProjectId":                "0",
		"scalingConfigurationName": d.Get("scalingConfigurationName").(string),
		"imageId":                  d.Get("imageId").(string),
		"cpu":                      IntToStr(d.Get("cpu").(int)),
		"mem":                      IntToStr(d.Get("mem").(int)),
		"storageType":              d.Get("storageType").(string),
		"storageSize":              IntToStr(d.Get("storageSize").(int)),
		"bandwidthType":            d.Get("bandwidthType").(string),
		"bandwidth":                IntToStr(d.Get("bandwidth").(int)),
		"imageType":                IntToStr(d.Get("imageType").(int)),
	}
	if rootSizeInfo, ok := d.GetOk("rootSize"); ok {
		params["rootSize"] = IntToStr(rootSizeInfo.(int))
	}
	if keyId, ok := d.GetOk("keyId"); ok {
		params["keyId"] = keyId.(string)
	}
	if password, ok := d.GetOk("password"); ok {
		params["password"] = password.(string)
	}
	if needMonitorAgent, ok := d.GetOk("needMonitorAgent"); ok {
		params["needMonitorAgent"] = IntToStr(needMonitorAgent.(int))
	}
	if wanIp, ok := d.GetOk("wanIp"); ok {
		params["wanIp"] = IntToStr(wanIp.(int))
	}
	if sgId, ok := d.GetOk("sgId"); ok {
		params["sgId"] = sgId.(string)
	}
	if projectId, ok := d.GetOk("projectId"); ok {
		params["projectId"] = projectId.(string)
	}
	if dataSnapshotId, ok := d.GetOk("dataSnapshotId"); ok {
		params["dataSnapshotId"] = dataSnapshotId.(string)
	}
	if cvmType, ok := d.GetOk("cvmType"); ok {
		params["cvmType"] = cvmType.(string)
	}
	response, err := client.SendRequest("scaling", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		CodeDesc string `json:"codeDesc"`
		Data     struct {
			totalCount              int
			scalingConfigurationSet []scalingConfiguration
		} `json:"data"`
	}
	if DEBUG {
		log.Printf("tencentcloud_autoscaling_group response: %v\n", response)
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("tencentcloud_autoscaling_group got error, code: %v, message: %v",
			jsonresp.Code, jsonresp.CodeDesc)
	}
	d.SetId(jsonresp.Data.scalingConfigurationSet[0].scalingConfigurationId)
	return resourceTencentCloudAutoScalingGroupRead(d, meta)
}

func resourceTencentCloudAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	scalingConfigurationId := d.Id()
	params := map[string]string{
		"Version":                   "2017-03-12",
		"Action":                    "DescribeScalingConfiguration",
		"scalingConfigurationIds.0": scalingConfigurationId,
	}
	client := meta.(*TencentCloudClient).commonConn
	response, err := client.SendRequest("scaling", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		code     int    `json:"code"`
		message  string `json:"message"`
		codeDesc string `json:"codeDesc"`
		data     struct {
			totalCount              int
			scalingConfigurationSet []scalingConfiguration
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.code != 0 {
		return fmt.Errorf("tencentcloud_autoscaling_group got error, code: %v, message: %v",
			jsonresp.code, jsonresp.codeDesc)
	}
	if jsonresp.data.totalCount == 0 {
		d.SetId("")
		return nil
	}
	d.SetId(jsonresp.data.scalingConfigurationSet[0].scalingConfigurationId)
	scalingConfigurationInfo := jsonresp.data.scalingConfigurationSet[0]
	if len(scalingConfigurationInfo.scalingConfigurationName) > 0 {
		d.Set("scalingConfigurationName", scalingConfigurationInfo.scalingConfigurationName)
	}
	if scalingConfigurationInfo.cpu > 0 {
		d.Set("cpu", scalingConfigurationInfo.cpu)
	}
	if scalingConfigurationInfo.mem > 0 {
		d.Set("mem", scalingConfigurationInfo.mem)
	}
	if scalingConfigurationInfo.rootSize > 0 {
		d.Set("rootSize", scalingConfigurationInfo.mem)
	}
	if len(scalingConfigurationInfo.keyId) > 0 {
		d.Set("keyId", scalingConfigurationInfo.keyId)
	}
	if scalingConfigurationInfo.bandwidth > 0 {
		d.Set("keyId", scalingConfigurationInfo.bandwidth)
	}
	if scalingConfigurationInfo.storageType > 0 {
		d.Set("storageType", scalingConfigurationInfo.storageType)
	}
	d.Set("createTime", scalingConfigurationInfo.createTime)
	d.Set("imageId", scalingConfigurationInfo.imageId)
	d.Set("imageType", scalingConfigurationInfo.imageType)
	d.Set("bandwidth", scalingConfigurationInfo.bandwidth)
	d.Set("bandwidthType", scalingConfigurationInfo.bandwidthType)
	d.Set("needMonitorAgent", scalingConfigurationInfo.needMonitorAgent)
	d.Set("needSecurityAgent", scalingConfigurationInfo.needSecurityAgent)
	d.Set("projectId", scalingConfigurationInfo.projectId)
	d.Set("storageSize", scalingConfigurationInfo.storageSize)
	return nil
}

func resourceTencentCloudAutoScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	scalingConfigurationId := d.Id()
	params := map[string]string{
		"Version":                "2017-03-12",
		"Action":                 "DeleteScalingGroup",
		"scalingConfigurationId": scalingConfigurationId,
	}
	client := meta.(*TencentCloudClient).commonConn
	response, err := client.SendRequest("scaling", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		code     int    `json:"code"`
		message  string `json:"message"`
		codeDesc string `json:"codeDesc"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	return nil
}
