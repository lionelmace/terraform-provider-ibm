package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMPIVolume() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPIVolumeCreate,
		Read:     resourceIBMPIVolumeRead,
		Update:   resourceIBMPIVolumeUpdate,
		Delete:   resourceIBMPIVolumeDelete,
		Exists:   resourceIBMPIVolumeExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			helpers.PIVolumeName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Volume Name to create",
			},

			helpers.PIVolumeShareable: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to indicate if the volume can be shared across multiple instances?",
			},
			helpers.PIVolumeSize: {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Size of the volume in GB",
			},
			helpers.PIVolumeType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ssd", "standard", "tier1", "tier3"}),
			},

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},

			// Computed Attributes

			"volume_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMPIVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	name := d.Get(helpers.PIVolumeName).(string)
	volType := d.Get(helpers.PIVolumeType).(string)
	size := float64(d.Get(helpers.PIVolumeSize).(float64))
	shared := d.Get(helpers.PIVolumeShareable).(bool)
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	vol, err := client.Create(name, size, volType, shared, powerinstanceid)

	if err != nil {
		return err
	}

	volumeid := *vol.VolumeID
	d.SetId(fmt.Sprintf("%s/%s", powerinstanceid, volumeid))
	if err != nil {
		log.Printf("[DEBUG]  err %s", isErrorToString(err))
		return err
	}
	_, err = isWaitForIBMPIVolumeAvailable(client, volumeid, powerinstanceid, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	//return nil
	return resourceIBMPIVolumeRead(d, meta)
}

func resourceIBMPIVolumeRead(d *schema.ResourceData, meta interface{}) error {
	sess, _ := meta.(ClientSession).IBMPISession()
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	powerinstanceid := parts[0]
	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	vol, err := client.Get(parts[1], powerinstanceid)
	if err != nil {
		return err
	}

	d.Set(helpers.PIVolumeName, vol.Name)
	d.Set(helpers.PIVolumeSize, vol.Size)
	d.Set(helpers.PIVolumeShareable, vol.Shareable)
	d.Set(helpers.PIVolumeType, vol.DiskType)
	d.Set("volume_status", vol.State)
	d.Set("volume_id", vol.VolumeID)
	d.Set(helpers.PICloudInstanceId, powerinstanceid)

	return nil
}

func resourceIBMPIVolumeUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("Calling the IBM Power Volume update call")
	sess, _ := meta.(ClientSession).IBMPISession()
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	powerinstanceid := parts[0]
	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	//name := ""
	//if d.HasChange(helpers.PIVolumeName) {
	name := d.Get(helpers.PIVolumeName).(string)
	//}

	size := float64(d.Get(helpers.PIVolumeSize).(float64))
	shareable := bool(d.Get(helpers.PIVolumeShareable).(bool))

	volrequest, err := client.Update(parts[1], name, size, shareable, powerinstanceid)
	if err != nil {
		return err
	}

	_, err = isWaitForIBMPIVolumeAvailable(client, *volrequest.VolumeID, powerinstanceid, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return resourceIBMPIVolumeRead(d, meta)
}

func resourceIBMPIVolumeDelete(d *schema.ResourceData, meta interface{}) error {

	sess, _ := meta.(ClientSession).IBMPISession()
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	powerinstanceid := parts[0]

	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	vol, err := client.Get(parts[1], powerinstanceid)
	if err != nil {
		return err
	}

	log.Printf("The volume to be deleted is in the following state .. %s", vol.State)
	_, err = isWaitForIBMPIVolumeAvailable(client, parts[1], powerinstanceid, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	voldelete_err := client.Delete(parts[1], powerinstanceid)
	if voldelete_err != nil {
		return voldelete_err
	}

	d.SetId("")
	return nil
}
func resourceIBMPIVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	powerinstanceid := parts[0]
	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	vol, err := client.Get(parts[1], powerinstanceid)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	log.Printf("Calling the existing function.. %s", *(vol.VolumeID))

	volumeid := *vol.VolumeID
	return volumeid == parts[1], nil
}

/*
func isWaitForIBMPIVolumeDeleted(vol *st.IBMPIVolumeClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVolumeDeleting},
		Target:     []string{},
		Refresh:    isIBMPIVolumeDeleteRefreshFunc(vol, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}*/

/*func isIBMPIVolumeDeleteRefreshFunc(vol *st.IBMPIVolumeClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := vol.Get(id)
		if err == nil {
			return vol, isVolumeDeleting, nil
		}

		iserror, ok := err.(iserrors.Power)
		if ok {
			log.Printf("[DEBUG] %s", iserror.Error())
			if len(iserror.Payload.Errors) == 1 &&
				iserror.Payload.Errors[0].Code == "volume_not_found" {
				return nil, isVolumeDeleted, nil
			}
		}
		return nil, isVolumeDeleting, err
	}
}

*/
func isWaitForIBMPIVolumeAvailable(client *st.IBMPIVolumeClient, id, powerinstanceid string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIVolumeProvisioning},
		Target:     []string{helpers.PIVolumeProvisioningDone},
		Refresh:    isIBMPIVolumeRefreshFunc(client, id, powerinstanceid),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isIBMPIVolumeRefreshFunc(client *st.IBMPIVolumeClient, id, powerinstanceid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id, powerinstanceid)
		if err != nil {
			return nil, "", err
		}

		if vol.State == "available" {
			return vol, helpers.PIVolumeProvisioningDone, nil
		}

		return vol, helpers.PIVolumeProvisioning, nil
	}
}
