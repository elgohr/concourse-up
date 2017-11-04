package iaas

import (
	"errors"
	"github.com/rackspace/gophercloud"
	"os"
)

const (
	OpenStackHostname    = "OPENSTACK_HOSTNAME"
	OpenStackUserName    = "OPENSTACK_USER"
	OpenStackPassword    = "OPENSTACK_PASSWORD"
	OpenStackProjectName = "OPENSTACK_PROJECT_NAME"
	OpenStackProjectID   = "OPENSTACK_PROJECT_ID"
)

type OpenStackClient struct {
	provider      *gophercloud.ProviderClient
	ComputeClient *gophercloud.ServiceClient
	region        string
}

func NewOpenStackClient(openStackAdapter OpenStack) (*OpenStackClient, []error) {
	var errorList []error
	hostname := os.Getenv(OpenStackHostname)
	username := os.Getenv(OpenStackUserName)
	password := os.Getenv(OpenStackPassword)
	projectName := os.Getenv(OpenStackProjectName)
	projectId := os.Getenv(OpenStackProjectID)

	if len(hostname) == 0 {
		errorList = append(errorList, errors.New("openstack hostname required"))
	}
	if len(username) == 0 {
		errorList = append(errorList, errors.New("openstack username required"))
	}
	if len(password) == 0 {
		errorList = append(errorList, errors.New("openstack password required"))
	}
	if len(projectName) == 0 &&
		len(projectId) == 0 {
		errorList = append(errorList, errors.New("openstack project id or project name required"))
	}
	if len(projectName) != 0 &&
		len(projectId) != 0 {
		errorList = append(errorList, errors.New("openstack project name and project id are set"))
	}
	if len(errorList) > 0 {
		return nil, errorList
	}
	client := &OpenStackClient{region: ""}
	if len(projectName) != 0 {
		authOptions := gophercloud.AuthOptions{
			IdentityEndpoint: hostname,
			Username:         username,
			Password:         password,
			TenantName:       projectName,
		}
		provider, err := openStackAdapter.AuthenticatedClient(authOptions)
		if err != nil {
			errorList = append(errorList, err)
		}
		client.ComputeClient, err = openStackAdapter.NewComputeV2(provider, gophercloud.EndpointOpts{
			Region: "",
		})
		if err != nil {
			errorList = append(errorList, err)
		}
	}
	if len(projectId) != 0 {
		authOptions := gophercloud.AuthOptions{
			IdentityEndpoint: hostname,
			Username:         username,
			Password:         password,
			TenantID:         projectId,
		}
		provider, err := openStackAdapter.AuthenticatedClient(authOptions)
		if err != nil {
			errorList = append(errorList, err)
		}
		client.ComputeClient, err = openStackAdapter.NewComputeV2(provider, gophercloud.EndpointOpts{
			Region: "",
		})
	}
	if len(errorList) > 0 {
		return nil, errorList
	}
	client.region = hostname
	return client, nil
}

func (o *OpenStackClient) IaaS() string {
	return "Openstack"
}

func (o *OpenStackClient) Region() string {
	return o.region
}

func (o *OpenStackClient) DeleteVMsInVPC(vpcID string) error {
	return errors.New("not implemented")
}
