package ec2client

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
	"os"
)

var (
	ec2Client *ec2.Client
)

func Initialize(awsProfile string) error {
	if awsProfile != "" {
		os.Setenv("AWS_PROFILE", awsProfile)

		// --profile command line argument is stronger than those environment variables
		//
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		os.Unsetenv("AWS_DEFAULT_REGION")
	}

	var err error

	ec2Client, err = initializeEc2Client()
	return err
}

func initializeAwsConfig() (aws.Config, error) {
	config, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Printf("Unable to load SDK config, %v\n", err)
		return config, err
	}

	creds, err := config.Credentials.Retrieve(context.TODO())

	if err != nil {
		log.Printf("Unable to get credentials %v", err)
		return config, err
	}

	log.Printf("AWS Config initialized, region: %v, access key: %v\n", config.Region, creds.AccessKeyID)
	return config, nil
}

func initializeEc2Client() (*ec2.Client, error) {
	awsConfig, err := initializeAwsConfig()

	if err != nil {
		return nil, err
	}

	return ec2.NewFromConfig(awsConfig), nil
}

func describeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	instances, err := ec2Client.DescribeInstances(context.TODO(), input)

	if err != nil {
		log.Printf("Error getting aws instances: %v\n", err)
		return instances, err
	}

	return instances, nil
}

func DescribeInstances() (*ec2.DescribeInstancesOutput, error) {
	input := &ec2.DescribeInstancesInput{}
	return describeInstances(input)
}

func DescribeInstance(instanceId string) (*types.Instance, error) {
	input := &ec2.DescribeInstancesInput{}
	input.InstanceIds = append(input.InstanceIds, instanceId)
	output, error := describeInstances(input)

	if error != nil {
		return nil, error
	}

	if len(output.Reservations) != 1 {
		log.Printf("Excpecting to get exactly 1 reservation for instance id: %v, got %v", instanceId, len(output.Reservations))
		return nil, errors.New("DescribeInstance bad Reservations")
	}

	if len(output.Reservations[0].Instances) != 1 {
		log.Printf("Excpecting to get exactly 1 instance for instance id: %v, got %v", instanceId, len(output.Reservations[0].Instances))
		return nil, errors.New("DescribeInstance bad Instances")
	}

	return &output.Reservations[0].Instances[0], nil
}
