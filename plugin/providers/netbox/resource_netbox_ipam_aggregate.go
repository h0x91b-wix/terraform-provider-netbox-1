package netbox

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/h0x91b-wix/go-netbox/netbox/client/ipam"
	"github.com/h0x91b-wix/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// resourceNetboxIpamAggregate is the core Terraform resource structure for the netbox_ipam_aggregate resource.
func resourceNetboxIpamAggregate() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamAggregateCreate,
		Read:   resourceNetboxIpamAggregateRead,
		Update: resourceNetboxIpamAggregateUpdate,
		Delete: resourceNetboxIpamAggregateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"prefix": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network prefix in slash notation for this aggregate. Example: 192.168.10.0/24.",
			},
			"rir_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Netbox ID of the regional internet registry (RIR) that manages this prefix.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of this aggregate.",
			},
			/*
				"date_added": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Date this aggregate was added.",
				},
			*/
		},
	}
}

// resourceNetboxIpamAggregateCreate creates a new aggregate in Netbox.
func resourceNetboxIpamAggregateCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	prefix := d.Get("prefix").(string)
	rirID := int64(d.Get("rir_id").(int))
	description := d.Get("description").(string)
	// TODO dateAdded
	tags := []string{}

	var parm = ipam.NewIpamAggregatesCreateParams().WithData(
		&models.WritableAggregate{
			Prefix:      &prefix,
			Description: description,
			Rir:         &rirID,
			Tags:        tags,
			// TODO DateAdded
		},
	)

	log.Debugf("Executing IpamAggregatesCreate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamAggregatesCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamAggregatesCreate: %v", err)

		return err
	}

	// TODO Probably a better way to parse this ID
	d.SetId(strconv.Itoa(int(out.Payload.ID)))

	log.Debugf("Done Executing IpamAggregatesCreate: %v", out)

	return nil
}

// resourceNetboxIpamAggregateUpdate applies updates to an aggregate by ID when deltas are detected by Terraform.
func resourceNetboxIpamAggregateUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Debugf("Error parsing Aggregate ID %v = %v", d.Id(), err)
		return err
	}

	prefix := d.Get("prefix").(string)
	rirID := int64(d.Get("rir_id").(int))
	description := d.Get("description").(string)
	// TODO dateAdded
	tags := []string{}

	var parm = ipam.NewIpamAggregatesUpdateParams().WithID(int64(id)).WithData(
		&models.WritableAggregate{
			Prefix:      &prefix,
			Description: description,
			Rir:         &rirID,
			Tags:        tags,
			// TODO DateAdded
		},
	)

	log.Debugf("Executing IpamAggregatesUpdate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamAggregatesUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamAggregatesUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing IpamAggregatesUpdate: %v", out)

	return nil
}

// resourceNetboxIpamAggregateRead reads an existing aggregate by ID.
func resourceNetboxIpamAggregateRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Debugf("Error parsing aggregate ID %v = %v", d.Id(), err)
		return err
	}

	var readParams = ipam.NewIpamAggregatesReadParams().WithID(int64(id))

	readResult, err := netboxClient.Ipam.IpamAggregatesRead(readParams, nil)

	if err != nil {
		log.Debugf("Error fetching aggregate ID # %d from Netbox = %v", id, err)
		return err
	}

	d.Set("prefix", readResult.Payload.Prefix)
	d.Set("rir_id", readResult.Payload.Rir.ID)
	d.Set("description", readResult.Payload.Description)
	// TODO date_created

	log.Debugf("Read Aggregate %d from Netbox = %v", readResult.Payload.Rir.ID, readResult.Payload)

	return nil
}

// resourceNetboxIpamAggregateDelete deletes an existing aggregate by ID.
func resourceNetboxIpamAggregateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting Aggregate: %v\n", d)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Debugf("Error parsing Aggregate ID %v = %v", d.Id(), err)
		return err
	}

	var deleteParameters = ipam.NewIpamAggregatesDeleteParams().WithID(int64(id))

	c := meta.(*ProviderNetboxClient).client

	out, err := c.Ipam.IpamAggregatesDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamAggregatesDelete: %v", err)
	}

	log.Debugf("Done Executing IpamAggregatesDelete: %v", out)

	return nil
}
