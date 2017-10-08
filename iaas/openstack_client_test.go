package iaas_test

import (
	. "github.com/EngineerBetter/concourse-up/iaas"

	"errors"
	"git.openstack.org/openstack/golang-client/openstack"
	"github.com/EngineerBetter/concourse-up/iaas/iaasfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var (
	openstackHostname    string
	openstackUsername    string
	openstackPassword    string
	openStackProjectName string
	openStackProjectID   string
)

type FakeAuthRef struct{}

func (a *FakeAuthRef) GetToken() string                           { return "" }
func (a *FakeAuthRef) GetExpiration() time.Time                   { return time.Now() }
func (a *FakeAuthRef) GetEndpoint(string, string) (string, error) { return "", nil }
func (a *FakeAuthRef) GetProject() (string)                       { return "" }

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
				fakeAuthRef *FakeAuthRef
			)
			BeforeEach(func() {
				os.Setenv(OpenStackHostname, "Hostname")
				os.Setenv(OpenStackUserName, "User")
				os.Setenv(OpenStackPassword, "Password")
				fakeAuthRef = &FakeAuthRef{}
				openStackAdapterMock.DoAuthRequestReturns(fakeAuthRef, nil)
			})
			It("is ok when authentication information are set with project name", func() {
				os.Setenv(OpenStackProjectName, "ProjectName")

				client, err := NewOpenStackClient(openStackAdapterMock)

				Expect(openStackAdapterMock.DoAuthRequestArgsForCall(0)).To(Equal(
					openstack.AuthOpts{
						AuthUrl:     "Hostname",
						Username:    "User",
						Password:    "Password",
						ProjectName: "ProjectName",
					}))
				Expect(err).Should(BeNil())
				Expect(client).To(BeAssignableToTypeOf(&OpenStackClient{}))
				Expect(client.AuthRef).To(Equal(fakeAuthRef))
			})
			It("is ok when authentication information are set with project id", func() {
				os.Setenv(OpenStackProjectID, "ProjectId")

				client, err := NewOpenStackClient(openStackAdapterMock)

				Expect(openStackAdapterMock.DoAuthRequestArgsForCall(0)).To(Equal(
					openstack.AuthOpts{
						AuthUrl:   "Hostname",
						Username:  "User",
						Password:  "Password",
						ProjectId: "ProjectId",
					}))
				Expect(err).Should(BeNil())
				Expect(client).To(BeAssignableToTypeOf(&OpenStackClient{}))
				Expect(client.AuthRef).To(Equal(fakeAuthRef))
			})
		})
		Describe("returns openstack api errors", func() {
			var (
				fakeAuthRef *FakeAuthRef
			)
			BeforeEach(func() {
				os.Setenv(OpenStackHostname, "Hostname")
				os.Setenv(OpenStackUserName, "User")
				os.Setenv(OpenStackPassword, "Password")
				fakeAuthRef = &FakeAuthRef{}
				openStackAdapterMock.DoAuthRequestReturns(nil, errors.New("any error"))
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
