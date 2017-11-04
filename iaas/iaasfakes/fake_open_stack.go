// Code generated by counterfeiter. DO NOT EDIT.
package iaasfakes

import (
	"sync"

	"github.com/EngineerBetter/concourse-up/iaas"
	"github.com/rackspace/gophercloud"
)

type FakeOpenStack struct {
	AuthenticatedClientStub        func(authopts gophercloud.AuthOptions) (*gophercloud.ProviderClient, error)
	authenticatedClientMutex       sync.RWMutex
	authenticatedClientArgsForCall []struct {
		authopts gophercloud.AuthOptions
	}
	authenticatedClientReturns struct {
		result1 *gophercloud.ProviderClient
		result2 error
	}
	authenticatedClientReturnsOnCall map[int]struct {
		result1 *gophercloud.ProviderClient
		result2 error
	}
	NewComputeV2Stub        func(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)
	newComputeV2Mutex       sync.RWMutex
	newComputeV2ArgsForCall []struct {
		client *gophercloud.ProviderClient
		eo     gophercloud.EndpointOpts
	}
	newComputeV2Returns struct {
		result1 *gophercloud.ServiceClient
		result2 error
	}
	newComputeV2ReturnsOnCall map[int]struct {
		result1 *gophercloud.ServiceClient
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOpenStack) AuthenticatedClient(authopts gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	fake.authenticatedClientMutex.Lock()
	ret, specificReturn := fake.authenticatedClientReturnsOnCall[len(fake.authenticatedClientArgsForCall)]
	fake.authenticatedClientArgsForCall = append(fake.authenticatedClientArgsForCall, struct {
		authopts gophercloud.AuthOptions
	}{authopts})
	fake.recordInvocation("AuthenticatedClient", []interface{}{authopts})
	fake.authenticatedClientMutex.Unlock()
	if fake.AuthenticatedClientStub != nil {
		return fake.AuthenticatedClientStub(authopts)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.authenticatedClientReturns.result1, fake.authenticatedClientReturns.result2
}

func (fake *FakeOpenStack) AuthenticatedClientCallCount() int {
	fake.authenticatedClientMutex.RLock()
	defer fake.authenticatedClientMutex.RUnlock()
	return len(fake.authenticatedClientArgsForCall)
}

func (fake *FakeOpenStack) AuthenticatedClientArgsForCall(i int) gophercloud.AuthOptions {
	fake.authenticatedClientMutex.RLock()
	defer fake.authenticatedClientMutex.RUnlock()
	return fake.authenticatedClientArgsForCall[i].authopts
}

func (fake *FakeOpenStack) AuthenticatedClientReturns(result1 *gophercloud.ProviderClient, result2 error) {
	fake.AuthenticatedClientStub = nil
	fake.authenticatedClientReturns = struct {
		result1 *gophercloud.ProviderClient
		result2 error
	}{result1, result2}
}

func (fake *FakeOpenStack) AuthenticatedClientReturnsOnCall(i int, result1 *gophercloud.ProviderClient, result2 error) {
	fake.AuthenticatedClientStub = nil
	if fake.authenticatedClientReturnsOnCall == nil {
		fake.authenticatedClientReturnsOnCall = make(map[int]struct {
			result1 *gophercloud.ProviderClient
			result2 error
		})
	}
	fake.authenticatedClientReturnsOnCall[i] = struct {
		result1 *gophercloud.ProviderClient
		result2 error
	}{result1, result2}
}

func (fake *FakeOpenStack) NewComputeV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	fake.newComputeV2Mutex.Lock()
	ret, specificReturn := fake.newComputeV2ReturnsOnCall[len(fake.newComputeV2ArgsForCall)]
	fake.newComputeV2ArgsForCall = append(fake.newComputeV2ArgsForCall, struct {
		client *gophercloud.ProviderClient
		eo     gophercloud.EndpointOpts
	}{client, eo})
	fake.recordInvocation("NewComputeV2", []interface{}{client, eo})
	fake.newComputeV2Mutex.Unlock()
	if fake.NewComputeV2Stub != nil {
		return fake.NewComputeV2Stub(client, eo)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newComputeV2Returns.result1, fake.newComputeV2Returns.result2
}

func (fake *FakeOpenStack) NewComputeV2CallCount() int {
	fake.newComputeV2Mutex.RLock()
	defer fake.newComputeV2Mutex.RUnlock()
	return len(fake.newComputeV2ArgsForCall)
}

func (fake *FakeOpenStack) NewComputeV2ArgsForCall(i int) (*gophercloud.ProviderClient, gophercloud.EndpointOpts) {
	fake.newComputeV2Mutex.RLock()
	defer fake.newComputeV2Mutex.RUnlock()
	return fake.newComputeV2ArgsForCall[i].client, fake.newComputeV2ArgsForCall[i].eo
}

func (fake *FakeOpenStack) NewComputeV2Returns(result1 *gophercloud.ServiceClient, result2 error) {
	fake.NewComputeV2Stub = nil
	fake.newComputeV2Returns = struct {
		result1 *gophercloud.ServiceClient
		result2 error
	}{result1, result2}
}

func (fake *FakeOpenStack) NewComputeV2ReturnsOnCall(i int, result1 *gophercloud.ServiceClient, result2 error) {
	fake.NewComputeV2Stub = nil
	if fake.newComputeV2ReturnsOnCall == nil {
		fake.newComputeV2ReturnsOnCall = make(map[int]struct {
			result1 *gophercloud.ServiceClient
			result2 error
		})
	}
	fake.newComputeV2ReturnsOnCall[i] = struct {
		result1 *gophercloud.ServiceClient
		result2 error
	}{result1, result2}
}

func (fake *FakeOpenStack) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.authenticatedClientMutex.RLock()
	defer fake.authenticatedClientMutex.RUnlock()
	fake.newComputeV2Mutex.RLock()
	defer fake.newComputeV2Mutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOpenStack) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ iaas.OpenStack = new(FakeOpenStack)