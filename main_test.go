package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var (
	cliPath string
)

var _ = Describe("concourse-up", func() {
	var ldFlags []string

	BeforeSuite(func() {
		compilationVars := map[string]string{}

		file, err := os.Open("compilation-vars.json")
		Expect(err).To(Succeed())
		defer file.Close()

		err = json.NewDecoder(file).Decode(&compilationVars)
		Expect(err).To(Succeed())

		ldFlags = []string{
			fmt.Sprintf("-X main.ConcourseUpVersion=%s", "0.0.0"),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.ConcourseStemcellURL=%s", compilationVars["concourse_stemcell_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.ConcourseStemcellVersion=%s", compilationVars["concourse_stemcell_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.ConcourseStemcellSHA1=%s", compilationVars["concourse_stemcell_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.ConcourseReleaseURL=%s", compilationVars["concourse_release_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.ConcourseReleaseVersion=%s", compilationVars["concourse_release_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.ConcourseReleaseSHA1=%s", compilationVars["concourse_release_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.RiemannReleaseURL=%s", compilationVars["riemann_release_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.RiemannReleaseVersion=%s", compilationVars["riemann_release_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.RiemannReleaseSHA1=%s", compilationVars["riemann_release_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.GrafanaReleaseURL=%s", compilationVars["grafana_release_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.GrafanaReleaseVersion=%s", compilationVars["grafana_release_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.GrafanaReleaseSHA1=%s", compilationVars["grafana_release_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.InfluxDBReleaseURL=%s", compilationVars["influxdb_release_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.InfluxDBReleaseVersion=%s", compilationVars["influxdb_release_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.InfluxDBReleaseSHA1=%s", compilationVars["influxdb_release_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.GardenReleaseURL=%s", compilationVars["garden_release_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.GardenReleaseVersion=%s", compilationVars["garden_release_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.GardenReleaseSHA1=%s", compilationVars["garden_release_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorStemcellURL=%s", compilationVars["director_stemcell_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorStemcellSHA1=%s", compilationVars["director_stemcell_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorStemcellVersion=%s", compilationVars["director_stemcell_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorCPIReleaseURL=%s", compilationVars["director_bosh_cpi_release_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorCPIReleaseVersion=%s", compilationVars["director_bosh_cpi_release_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorCPIReleaseSHA1=%s", compilationVars["director_bosh_cpi_release_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorReleaseURL=%s", compilationVars["director_bosh_release_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorReleaseVersion=%s", compilationVars["director_bosh_release_version"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/bosh.DirectorReleaseSHA1=%s", compilationVars["director_bosh_release_sha1"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/fly.DarwinBinaryURL=%s", compilationVars["fly_darwin_binary_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/fly.LinuxBinaryURL=%s", compilationVars["fly_linux_binary_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/fly.WindowsBinaryURL=%s", compilationVars["fly_windows_binary_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/director.DarwinBinaryURL=%s", compilationVars["director_darwin_binary_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/director.LinuxBinaryURL=%s", compilationVars["director_linux_binary_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/director.WindowsBinaryURL=%s", compilationVars["director_windows_binary_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/terraform.DarwinBinaryURL=%s", compilationVars["terraform_darwin_binary_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/terraform.LinuxBinaryURL=%s", compilationVars["terraform_linux_binary_url"]),
			fmt.Sprintf("-X github.com/EngineerBetter/concourse-up/terraform.WindowsBinaryURL=%s", compilationVars["terraform_windows_binary_url"]),
		}

		cliPath, err = Build("github.com/EngineerBetter/concourse-up", "-ldflags", strings.Join(ldFlags, " "))
		Expect(err).ToNot(HaveOccurred(), "Error building source")
	})

	AfterSuite(func() {
		CleanupBuildArtifacts()
	})

	It("displays usage instructions on --help", func() {
		command := exec.Command(cliPath, "--help")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred(), "Error running CLI: "+cliPath)
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Concourse-Up - A CLI tool to deploy Concourse CI"))
		Expect(session.Out).To(Say("deploy, d   Deploys or updates a Concourse"))
		Expect(session.Out).To(Say("destroy, x  Destroys a Concourse"))
	})

	Context("When a compile-time variable is missing", func() {
		It("Returns an error", func() {
			ldFlagsMising := ldFlags[0 : len(ldFlags)-1]

			cliPath, err := Build("github.com/EngineerBetter/concourse-up", "-ldflags", strings.Join(ldFlagsMising, " "))
			Expect(err).ToNot(HaveOccurred(), "Error building source")

			command := exec.Command(cliPath, "--help")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred(), "Error running CLI: "+cliPath)

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("Compile-time variable terraform.WindowsBinaryURL not set, please build with: `go build -ldflags \"-X github.com/EngineerBetter/concourse-up/terraform.WindowsBinaryURL=SOME_VALUE\"`"))
		})
	})
})
