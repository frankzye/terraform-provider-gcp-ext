package main

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	zone := d.Get("zone").(string)
	instance := d.Get("instance").(string)
	policy := d.Get("policy").(string)
	uid, _ := uuid.NewUUID()
	// delete resource policy no matter exist or not
	deleteResourcePolicy(project, zone, instance, policy)
	addResourcePolicy(project, zone, instance, policy)
	d.SetId(uid.String())
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChanges("instance", "policy", "zone") {
		project := d.Get("project").(string)
		zone, newZone := d.GetChange("zone")
		instance, newInstance := d.GetChange("instance")
		policy, newPolicy := d.GetChange("policy")

		deleteResourcePolicy(project, zone.(string), instance.(string), policy.(string))
		addResourcePolicy(project, newZone.(string), newInstance.(string), newPolicy.(string))
	}
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	zone := d.Get("zone").(string)
	instance := d.Get("instance").(string)
	policy := d.Get("policy").(string)
	deleteResourcePolicy(project, zone, instance, policy)
	return nil
}

func addResourcePolicy(projectID string, zone string, instance string, policy string) {
	ctx := context.Background()
	uid, _ := uuid.NewUUID()
	requestId := uid.String()
	instanceClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		fmt.Println("NewInstancesRESTClient:", err)
		return
	}
	req := &computepb.AddResourcePoliciesInstanceRequest{
		Project:   projectID,
		Zone:      zone,
		RequestId: &requestId,
		Instance:  instance,
		InstancesAddResourcePoliciesRequestResource: &computepb.InstancesAddResourcePoliciesRequest{ResourcePolicies: []string{policy}},
	}
	op, err := instanceClient.AddResourcePolicies(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	waitErr := op.Wait(ctx)
	if waitErr != nil {
		fmt.Println(waitErr.Error())
	}
	return
}

func deleteResourcePolicy(projectID string, zone string, instance string, policy string) {
	ctx := context.Background()
	uid, _ := uuid.NewUUID()
	requestId := uid.String()
	instanceClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		fmt.Println("NewInstancesRESTClient:", err)
		return
	}
	req := &computepb.RemoveResourcePoliciesInstanceRequest{
		Project:   projectID,
		Zone:      zone,
		RequestId: &requestId,
		Instance:  instance,
		InstancesRemoveResourcePoliciesRequestResource: &computepb.InstancesRemoveResourcePoliciesRequest{ResourcePolicies: []string{policy}},
	}
	op, err := instanceClient.RemoveResourcePolicies(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	waitErr := op.Wait(ctx)
	if waitErr != nil {
		fmt.Println(waitErr.Error())
	}
}
