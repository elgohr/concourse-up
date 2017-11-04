package iaas_test

import (
	. "github.com/EngineerBetter/concourse-up/iaas"

	"errors"
	"github.com/EngineerBetter/concourse-up/iaas/iaasfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rackspace/gophercloud"
	"os"
)

var (
	openstackHostname    string
	openstackUsername    string
	openstackPassword    string
	openStackProjectName string
	openStackProjectID   string
)

var _ = Describe("OpenstackClient", func() {

	var (
		openStackAdapterMock *iaasfakes.FakeOpenStack
	)

	BeforeEach(func() {
		saveCurrentEnvironmentVariables()
		clearEnvironmentVariables()
		openStackAdapterMock = new(iaasfakes.FakeOpenStack)
	})

	AfterEach(func() {
		restoreEnvironmentVariables()
	})

	Describe("methods", func() {
		var (
			client *OpenStackClient
		)
		BeforeEach(func() {
			os.Setenv(OpenStackHostname, "https://host.name")
			os.Setenv(OpenStackUserName, "User")
			os.Setenv(OpenStackPassword, "Password")
			os.Setenv(OpenStackProjectName, "ProjectName")
			openStackAdapterMock.AuthenticatedClient(gophercloud.AuthOptions{})
			client, _ = NewOpenStackClient(openStackAdapterMock)

		})
		Describe("IaaS", func() {
			It("returns it's name", func() {
				Expect(client.IaaS()).To(Equal("Openstack"))
			})
		})
		Describe("Region", func() {
			It("returns the configured endpoint", func() {
				Expect(client.Region()).To(Equal("https://host.name"))
			})
		})
		Describe("DeleteVMsInVPC", func() {
			It("deletes all machines in a region", func() {
				//Expect(client.DeleteVMsInVPC("ANY")).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("NewOpenStackClient", func() {
		It("errors when hostname is not set", func() {
			_, err := NewOpenStackClient(openStackAdapterMock)
			Expect(err).To(ContainElement(errors.New("openstack hostname required")))
		})
		It("errors when username is not set", func() {
			_, err := NewOpenStackClient(openStackAdapterMock)
			Expect(err).To(ContainElement(errors.New("openstack username required")))
		})
		It("errors when password is not set", func() {
			_, err := NewOpenStackClient(openStackAdapterMock)
			Expect(err).To(ContainElement(errors.New("openstack password required")))
		})
		It("errors when projectName is not set", func() {
			_, err := NewOpenStackClient(openStackAdapterMock)
			Expect(err).To(ContainElement(errors.New("openstack project id or project name required")))
		})
		It("errors when projectId is not set", func() {
			_, err := NewOpenStackClient(openStackAdapterMock)
			Expect(err).To(ContainElement(errors.New("openstack project id or project name required")))
		})
		It("errors when projectName and projectId are set", func() {
			os.Setenv(OpenStackProjectName, "ProjectName")
			os.Setenv(OpenStackProjectID, "ProjectId")
			_, err := NewOpenStackClient(openStackAdapterMock)
			Expect(err).To(ContainElement(errors.New("openstack project name and project id are set")))
		})

		Describe("passes checks", func() {
			var (
				providerMock      *gophercloud.ProviderClient
				serviceClientMock *gophercloud.ServiceClient
			)

			BeforeEach(func() {
				os.Setenv(OpenStackHostname, "https://host.name")
				os.Setenv(OpenStackUserName, "User")
				os.Setenv(OpenStackPassword, "Password")
				providerMock = &gophercloud.ProviderClient{}
				openStackAdapterMock.AuthenticatedClientReturns(providerMock, nil)
				serviceClientMock = &gophercloud.ServiceClient{}
				openStackAdapterMock.NewComputeV2Returns(serviceClientMock, nil)
			})
			It("is ok when authentication information are set with project name", func() {
				os.Setenv(OpenStackProjectName, "ProjectName")

				client, err := NewOpenStackClient(openStackAdapterMock)

				Expect(openStackAdapterMock.AuthenticatedClientArgsForCall(0)).To(Equal(
					gophercloud.AuthOptions{
						IdentityEndpoint: "https://host.name",
						Username:         "User",
						Password:         "Password",
						TenantName:       "ProjectName",
					}))
				Expect(err).Should(BeNil())
				Expect(client).To(BeAssignableToTypeOf(&OpenStackClient{}))
				Expect(client.ComputeClient).To(Equal(serviceClientMock))
				Expect(openStackAdapterMock.NewComputeV2CallCount()).To(Equal(1))
				providerClient, eo := openStackAdapterMock.NewComputeV2ArgsForCall(0)
				Expect(providerClient).To(Equal(providerMock))
				Expect(eo).To(Equal(gophercloud.EndpointOpts{
					Region: "",
				}))
			})
			It("is ok when authentication information are set with project id", func() {
				os.Setenv(OpenStackProjectID, "ProjectId")

				client, err := NewOpenStackClient(openStackAdapterMock)

				Expect(openStackAdapterMock.AuthenticatedClientArgsForCall(0)).To(Equal(
					gophercloud.AuthOptions{
						IdentityEndpoint: "https://host.name",
						Username:         "User",
						Password:         "Password",
						TenantID:         "ProjectId",
					}))
				Expect(err).Should(BeNil())
				Expect(client).To(BeAssignableToTypeOf(&OpenStackClient{}))
				Expect(client.ComputeClient).To(Equal(serviceClientMock))
				Expect(openStackAdapterMock.NewComputeV2CallCount()).To(Equal(1))
				providerClient, eo := openStackAdapterMock.NewComputeV2ArgsForCall(0)
				Expect(providerClient).To(Equal(providerMock))
				Expect(eo).To(Equal(gophercloud.EndpointOpts{
					Region: "",
				}))
			})
		})
		Describe("returns openstack api errors", func() {
			BeforeEach(func() {
				os.Setenv(OpenStackHostname, "Hostname")
				os.Setenv(OpenStackUserName, "User")
				os.Setenv(OpenStackPassword, "Password")
				openStackAdapterMock.AuthenticatedClientReturns(nil, errors.New("any error"))
			})
			It("is ok when authentication information are set with project name", func() {
				os.Setenv(OpenStackProjectName, "ProjectName")

				_, err := NewOpenStackClient(openStackAdapterMock)

				Expect(err).To(ContainElement(errors.New("any error")))
			})
			It("is ok when authentication information are set with project id", func() {
				os.Setenv(OpenStackProjectID, "ProjectId")

				_, err := NewOpenStackClient(openStackAdapterMock)

				Expect(err).To(ContainElement(errors.New("any error")))
			})
		})
	})

})

func clearEnvironmentVariables() {
	os.Clearenv()
}

func restoreEnvironmentVariables() {
	os.Setenv(OpenStackHostname, openstackHostname)
	os.Setenv(OpenStackUserName, openstackUsername)
	os.Setenv(OpenStackPassword, openstackPassword)
	os.Setenv(OpenStackProjectName, openStackProjectName)
	os.Setenv(OpenStackProjectID, openStackProjectID)
}

func saveCurrentEnvironmentVariables() {
	openstackHostname = os.Getenv(OpenStackHostname)
	openstackUsername = os.Getenv(OpenStackUserName)
	openstackPassword = os.Getenv(OpenStackPassword)
	openStackProjectName = os.Getenv(OpenStackProjectName)
	openStackProjectID = os.Getenv(OpenStackProjectID)
}
