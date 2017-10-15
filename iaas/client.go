package iaas

import (
	"fmt"
)

// IClient represents actions taken against AWS
//go:generate counterfeiter . IClient
type IClient interface {
	Platform
	ObjectStore
}

type Platform interface {
	// IaaS returns the name of the provider to operate against
	IaaS() string
	// Region returns the region to operate against
	Region() string
	// DeleteVMsInVPC deletes all the VMs in the given VPC
	DeleteVMsInVPC(vpcID string) error
	// FindLongestMatchingHostedZone finds the longest hosted zone that matches the given subdomain
	FindLongestMatchingHostedZone(subDomain string) (string, string, error)
	// MockProvider enables to use mock objects in testing
	MockProvider(interface{})
}

type ObjectStore interface {
	// DeleteVersionedBucket deletes and empties a versioned bucket
	DeleteVersionedBucket(name string) error
	// DeleteFile deletes a file
	DeleteFile(bucket, path string) error
	// EnsureBucketExists checks if the named bucket exists and creates it if it doesn't
	EnsureBucketExists(name string) error
	// EnsureFileExists checks for the named file and creates it if it doesn't
	// Second argument is true if new file was created
	EnsureFileExists(bucket, path string, defaultContents []byte) ([]byte, bool, error)
	// HasFile returns true if the specified object exists
	HasFile(bucket, path string) (bool, error)
	// LoadFile loads a file
	LoadFile(bucket, path string) ([]byte, error)
	// WriteFile writes the specified object
	WriteFile(bucket, path string, contents []byte) error
}

// New returns a new IAAS client for a particular IAAS and region
func New(iaas string, region string) (IClient, error) {
	if iaas == "AWS" {
		return newAWS(region)
	}

	return nil, fmt.Errorf("IAAS not supported: %s", iaas)
}
