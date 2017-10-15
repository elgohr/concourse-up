package iaas

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
)

// AWSClient is the concrete implementation of IClient on AWS
// using Route53 & S3
type AWSClient struct {
	region  string
	route53 Route53
}

func newAWS(region string) (IClient, error) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		return nil, errors.New("env var AWS_ACCESS_KEY_ID not found")
	}

	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		return nil, errors.New("env var AWS_SECRET_ACCESS_KEY not found")
	}

	route53Client, err := newRoute53Client()

	if err != nil {
		return nil, err
	}

	return &AWSClient{region: region, route53: route53Client}, nil
}

func (client *AWSClient) Region() string {
	return client.region
}

func (client *AWSClient) IaaS() string {
	return "AWS"
}

func (client *AWSClient) DeleteVMsInVPC(vpcID string) error {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return err
	}

	filterName := "vpc-id"
	ec2Client := ec2.New(sess, &aws.Config{Region: &client.region})

	resp, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: &filterName,
				Values: []*string{
					&vpcID,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	instancesToTerminate := []*string{}
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Printf("Terminating instance %s\n", *instance.InstanceId)
			instancesToTerminate = append(instancesToTerminate, instance.InstanceId)
		}
	}

	if len(instancesToTerminate) == 0 {
		return nil
	}

	_, err = ec2Client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: instancesToTerminate,
	})
	return err
}

func (client *AWSClient) FindLongestMatchingHostedZone(subDomain string) (string, string, error) {
	hostedZones := []*route53.HostedZone{}
	err := client.route53.ListHostedZonesPages(&route53.ListHostedZonesInput{}, func(output *route53.ListHostedZonesOutput, _ bool) bool {
		hostedZones = append(hostedZones, output.HostedZones...)
		return true
	})
	if err != nil {
		return "", "", err
	}

	longestMatchingHostedZoneName := ""
	longestMatchingHostedZoneID := ""
	for i := 0; i < len(hostedZones); i++ {
		domain := strings.TrimRight(*hostedZones[i].Name, ".")

		id := *hostedZones[i].Id
		if strings.HasSuffix(subDomain, domain) {
			if len(domain) > len(longestMatchingHostedZoneName) {
				longestMatchingHostedZoneName = domain
				longestMatchingHostedZoneID = id
			}
		}
	}

	if longestMatchingHostedZoneName == "" {
		return "", "", fmt.Errorf("No matching hosted zone found for domain %s", subDomain)
	}

	longestMatchingHostedZoneID = strings.Replace(longestMatchingHostedZoneID, "/hostedzone/", "", -1)

	return longestMatchingHostedZoneName, longestMatchingHostedZoneID, err
}

func (client *AWSClient) MockProvider(backendStub interface{}) {
	client.route53 = backendStub.(Route53)
}

// NewRoute53Client returns a new Route53 client
func newRoute53Client() (Route53, error) {
	session, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return nil, err
	}
	return route53.New(session), nil
}
