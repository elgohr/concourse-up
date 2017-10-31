package iaas

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

//go:generate counterfeiter . OpenStack
type OpenStack interface {
	AuthenticatedClient(authopts gophercloud.AuthOptions) (*gophercloud.ProviderClient, error)
}

// This struct only exists to wrap static method-calls from OpenStack-Api
// due to testability. If you write an Api yourself, remember this!
type OpenStackAdapter struct {
}

func (oa *OpenStackAdapter) AuthenticatedClient(authopts gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	return openstack.AuthenticatedClient(authopts)
}
