package rescue

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

var provisionerOnlineClient online.Client

func applyRescueProvisioner(ctx context.Context) error {
	d := ctx.Value(schema.ProvConfigDataKey).(*schema.ResourceData)

	if provisionerOnlineClient == nil {
		token := d.Get("token").(string)
		provisionerOnlineClient = online.NewClient(token)
	}

	enabled := d.Get("enabled").(bool)
	server := d.Get("server").(int)
	imageInterface, imageExists := d.GetOk("image")
	pathInterface, pathExists := d.GetOk("credentials_dir")

	if enabled && !imageExists {
		return errors.New("Need image to enable rescue mode")
	}

	if enabled && !pathExists {
		return errors.New("Need credentials_dir to enable rescue mode")
	}

	if enabled {
		credentials, err := provisionerOnlineClient.BootRescueMode(server, imageInterface.(string))
		if err != nil {
			return err
		}
		err = writeCredentials(credentials, pathInterface.(string))
		if err != nil {
			return err
		}
	} else {
		err := provisionerOnlineClient.BootNormalMode(server)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeCredentials is a hack as provisioners can not have output data.
// This will write the credentials to individual files that can be picked up by ${file()}
func writeCredentials(credentials *online.RescueCredentials, to string) error {

	err := os.Mkdir(to, 0755)
	if err != nil && !strings.Contains(err.Error(), "file exists") { // do not error if the directory already exists
		return err
	}

	err = ioutil.WriteFile(path.Join(to, "username"), []byte(credentials.Login), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(to, "password"), []byte(credentials.Password), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(to, "ip"), []byte(credentials.IP), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(to, "protocol"), []byte(credentials.Protocol), 0644)
	if err != nil {
		return err
	}

	return nil
}
