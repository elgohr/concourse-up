package config_test

import (
	. "github.com/EngineerBetter/concourse-up/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/EngineerBetter/concourse-up/iaas/iaasfakes"
)

var _ = Describe("Client", func() {
	var (
		iaasClient *iaasfakes.FakeIClient
		client     *Client
		deployArgs *DeployArgs
	)

	BeforeEach(func() {
		iaasClient = new(iaasfakes.FakeIClient)
		iaasClient.RegionReturns("eu-west-1")
		iaasClient.EnsureBucketExistsReturns(nil)
		iaasClient.EnsureFileExistsStub = func(bucket, path string, defaultContents []byte) ([]byte, bool, error) {
			return defaultContents, true, nil
		}

		client = New(iaasClient, "test")

		deployArgs = &DeployArgs{
			IAAS:        "AWS",
			AWSRegion:   "eu-west-1",
			WorkerCount: 1,
			WorkerSize:  "xlarge",
			DBSize:      "medium",
			DBSizeIsSet: false,
		}
	})

	Describe("LoadOrCreate", func() {
		Context("When the there is no existing config", func() {
			var conf *Config
			var createdANewFile bool

			BeforeEach(func() {
				var err error
				conf, createdANewFile, err = client.LoadOrCreate(deployArgs)
				Expect(err).To(Succeed())
			})

			It("creates a new file", func() {
				Expect(createdANewFile).To(BeTrue())
			})

			Describe("the default config file", func() {
				It("Sets the default value for the AvailabilityZone", func() {
					Expect(conf.AvailabilityZone).To(Equal("eu-west-1a"))
				})

				It("Sets the default value for the ConcourseDBName", func() {
					Expect(conf.ConcourseDBName).To(Equal("concourse_atc"))
				})

				It("Generates a secure random string for ConcoursePassword", func() {
					Expect(conf.ConcoursePassword).To(beARandomPassword())
				})

				It("Sets the default value for the ConcourseUsername", func() {
					Expect(conf.ConcourseUsername).To(Equal("admin"))
				})

				It("Sets the default value for the ConcourseWorkerCount", func() {
					Expect(conf.ConcourseWorkerCount).To(Equal(1))
				})

				It("Sets the default value for the ConcourseWorkerSize", func() {
					Expect(conf.ConcourseWorkerSize).To(Equal("xlarge"))
				})

				It("Sets the default value for the ConfigBucket", func() {
					Expect(conf.ConfigBucket).To(Equal("concourse-up-test-eu-west-1-config"))
				})

				It("Sets the default value for the Deployment", func() {
					Expect(conf.Deployment).To(Equal("concourse-up-test"))
				})

				It("Generates a secure random string for DirectorHMUserPassword", func() {
					Expect(conf.DirectorHMUserPassword).To(beARandomPassword())
				})

				It("Generates a secure random string for DirectorMbusPassword", func() {
					Expect(conf.DirectorMbusPassword).To(beARandomPassword())
				})

				It("Generates a secure random string for DirectorNATSPassword", func() {
					Expect(conf.DirectorNATSPassword).To(beARandomPassword())
				})

				It("Generates a secure random string for DirectorPassword", func() {
					Expect(conf.DirectorPassword).To(beARandomPassword())
				})

				It("Generates a secure random string for DirectorRegistryPassword", func() {
					Expect(conf.DirectorRegistryPassword).To(beARandomPassword())
				})

				It("Sets the default value for the DirectorUsername", func() {
					Expect(conf.DirectorUsername).To(Equal("admin"))
				})

				It("Generates a secure random string for EncryptionKey", func() {
					Expect(conf.EncryptionKey).To(MatchRegexp("^[a-z0-9]{32}$"))
				})

				It("Sets the GrafanaPassword to the ConcoursePassword", func() {
					Expect(conf.GrafanaPassword).To(Equal(conf.ConcoursePassword))
				})

				It("Sets the default value for the GrafanaUsername", func() {
					Expect(conf.GrafanaUsername).To(Equal("admin"))
				})

				It("Generates a secure random string for InfluxDBPassword", func() {
					Expect(conf.InfluxDBPassword).To(beARandomPassword())
				})

				It("Sets the default value for the InfluxDBUsername", func() {
					Expect(conf.InfluxDBUsername).To(Equal("admin"))
				})

				It("Sets the default value for the MultiAZRDS", func() {
					Expect(conf.MultiAZRDS).To(Equal(false))
				})

				It("Generates a random RSA private key for PrivateKey", func() {
					Expect(conf.PrivateKey).To(HavePrefix("-----BEGIN RSA PRIVATE KEY-----"))
				})

				It("Sets the default value for the Project", func() {
					Expect(conf.Project).To(Equal("test"))
				})

				It("Generates a random RSA public key for PublicKey", func() {
					Expect(conf.PublicKey).To(HavePrefix("ssh-rsa"))
				})

				It("Sets the default value for the RDSDefaultDatabaseName", func() {
					Expect(conf.RDSDefaultDatabaseName).To(Equal("bosh"))
				})

				It("Sets the default value for the RDSInstanceClass", func() {
					Expect(conf.RDSInstanceClass).To(Equal("db.t2.medium"))
				})

				It("Generates a secure random string for RDSPassword", func() {
					Expect(conf.RDSPassword).To(beARandomPassword())
				})

				It("Generates a secure random string for the RDSUsername", func() {
					Expect(conf.RDSUsername).To(MatchRegexp("^admin[a-z0-9]{20}$"))
				})

				It("Sets the default value for the Region", func() {
					Expect(conf.Region).To(Equal("eu-west-1"))
				})

				It("Sets the default value for the TFStatePath", func() {
					Expect(conf.TFStatePath).To(Equal("terraform.tfstate"))
				})
			})
		})
	})
})

func beARandomPassword() types.GomegaMatcher {
	return MatchRegexp("^[a-z0-9]{20}$")
}
