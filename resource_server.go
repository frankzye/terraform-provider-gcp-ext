package main

import (
	compute "cloud.google.com/go/compute/apiv1"
	"context"
	"fmt"
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
	err := addResourcePolicy(project, zone, instance, policy)
	if err != nil {
		return err
	}
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

		err := deleteResourcePolicy(project, zone.(string), instance.(string), policy.(string))
		if err != nil {
			return err
		}

		err = addResourcePolicy(project, newZone.(string), newInstance.(string), newPolicy.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	zone := d.Get("zone").(string)
	instance := d.Get("instance").(string)
	policy := d.Get("policy").(string)
	err := deleteResourcePolicy(project, zone, instance, policy)
	if err != nil {
		return err
	}
	return nil
}

func addResourcePolicy(projectID string, zone string, instance string, policy string) error {
	ctx := context.Background()
	uid, _ := uuid.NewUUID()
	requestId := uid.String()
	instanceClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
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
		return fmt.Errorf(err.Error())
	}
	waitErr := op.Wait(ctx)
	if waitErr != nil {
		return fmt.Errorf(waitErr.Error())
	}
	return nil
}

func deleteResourcePolicy(projectID string, zone string, instance string, policy string) error {
	ctx := context.Background()
	uid, _ := uuid.NewUUID()
	requestId := uid.String()
	instanceClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
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
		return fmt.Errorf(err.Error())
	}
	waitErr := op.Wait(ctx)
	if waitErr != nil {
		return fmt.Errorf(waitErr.Error())
	}
	return nil
}
