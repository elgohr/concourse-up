package iaas

import "git.openstack.org/openstack/golang-client/openstack"

//go:generate counterfeiter . OpenStack
type OpenStack interface {
	DoAuthRequest(authopts openstack.AuthOpts) (openstack.AuthRef, error)
}

// This struct only exists to wrap static method-calls from OpenStack-Api
// due to testability. If you write an Api yourself, remember this!
type OpenStackAdapter struct {
}

func (oa *OpenStackAdapter) DoAuthRequest(authopts openstack.AuthOpts) (openstack.AuthRef, error) {
	return openstack.DoAuthRequest(authopts)
}
