package iaas

import (
	"errors"
	"git.openstack.org/openstack/golang-client/openstack"
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
	AuthRef openstack.AuthRef
	region  string
}

func NewOpenStackClient(openStackAdapter OpenStack) (*OpenStackClient, []error) {
	errorList := []error{}
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
		auth, err := openStackAdapter.DoAuthRequest(openstack.AuthOpts{
			AuthUrl:     hostname,
			Username:    username,
			Password:    password,
			ProjectName: projectName,
		})
		if err != nil {
			errorList = append(errorList, err)
		}
		client.AuthRef = auth
	}
	if len(projectId) != 0 {
		auth, err := openStackAdapter.DoAuthRequest(openstack.AuthOpts{
			AuthUrl:   hostname,
			Username:  username,
			Password:  password,
			ProjectId: projectId,
		})
		if err != nil {
			errorList = append(errorList, err)
		}
		client.AuthRef = auth
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
