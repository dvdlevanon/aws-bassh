package loader

import (
	"aws-bassh/pkg/ec2client"
	"aws-bassh/pkg/model"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func LoadAllMachines(config model.GenerateConfig) (map[string]model.Machine, error) {
	instances, err := ec2client.DescribeInstances()

	if err != nil {
		return nil, err
	}

	return buildModelFromInstances(instances, config)
}

type buildModelContext struct {
	config       model.GenerateConfig
	allInstances *ec2.DescribeInstancesOutput
}

func buildModelFromInstances(instances *ec2.DescribeInstancesOutput, config model.GenerateConfig) (map[string]model.Machine, error) {
	machines := make(map[string]model.Machine)
	context := buildModelContext{
		config:       config,
		allInstances: instances,
	}

	for _, reservations := range instances.Reservations {
		for _, instance := range reservations.Instances {
			if instance.State.Name == types.InstanceStateNameTerminated {
				continue
			}

			machine := buildModelForMachine(instance, buildTagsMap(instance), context)

			machines[machine.Id] = machine
		}
	}

	return machines, nil
}

func buildTagsMap(instance types.Instance) map[string]string {
	tags := make(map[string]string)
	for _, tag := range instance.Tags {
		tags[*tag.Key] = *tag.Value
	}

	return tags
}

func buildModelForMachine(instance types.Instance, tags map[string]string, context buildModelContext) model.Machine {
	return model.Machine{
		Id:      *instance.InstanceId,
		Name:    findMachineName(instance, tags, context),
		User:    findUserName(instance, tags, context),
		Keyfile: findKeyFile(instance, context),
		Bastion: findBastion(instance, tags, context),
	}
}

func findMachineName(instance types.Instance, tags map[string]string, context buildModelContext) string {
	nameFromTag := getFirstTag(context.config.NameTags, tags)

	if nameFromTag != "" {
		return nameFromTag
	}

	return model.NoMachineName
}

func findUserName(instance types.Instance, tags map[string]string, context buildModelContext) string {
	userFromTag := getFirstTag(context.config.UserTags, tags)

	if userFromTag != "" {
		return userFromTag
	}

	// According to this doc:
	//	https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/connection-prereqs.html
	//
	switch distro := getDistroName(instance.ImageId); distro {
	case "amazon-linux":
		return "ec2-user"
	case "centos":
		return "centos"
	case "debian":
		return "admin"
	case "fedora":
		return "ec2-user"
	case "rhel":
		return "ec2-user"
	case "suse":
		return "ec2-user"
	case "ubuntu":
		return "ubuntu"
	default:
		return "ec2-user"
	}
}

func findKeyFile(instance types.Instance, context buildModelContext) string {
	if instance.KeyName != nil {
		return buildKeyfile(*instance.KeyName, context)
	}

	return ""
}

func findBastion(instance types.Instance, tags map[string]string, context buildModelContext) model.BastionMachine {
	bastionFromTags := getBastionFromTags(instance, tags, context)

	if bastionFromTags != model.NoBastion {
		return bastionFromTags
	}

	return findBastionInVpc(instance, context)
}

func getBastionFromTags(instance types.Instance, tags map[string]string, context buildModelContext) model.BastionMachine {
	bastionUrl := getFirstTag(context.config.BastionUrlTags, tags)
	bastionUser := getFirstTag(context.config.BastionUserTags, tags)
	bastionKey := getFirstTag(context.config.BastionKeyNameTags, tags)

	if bastionUrl == "" {
		return model.NoBastion
	}

	if bastionUser == "" {
		return model.NoBastion
	}

	if bastionKey == "" {
		bastionKey = *instance.KeyName
	}

	return model.BastionMachine{
		Url:     bastionUrl,
		User:    bastionUser,
		Keyfile: buildKeyfile(bastionKey, context),
	}
}

func findBastionInVpc(instance types.Instance, context buildModelContext) model.BastionMachine {
	if instance.VpcId == nil {
		return model.NoBastion
	}

	for _, reservations := range context.allInstances.Reservations {
		for _, bastionCandidate := range reservations.Instances {
			if bastionCandidate.VpcId == nil {
				continue
			}

			if bastionCandidate.PublicIpAddress == nil {
				continue
			}

			if *bastionCandidate.VpcId != *instance.VpcId {
				continue
			}

			if !containsBastionInName(bastionCandidate) {
				continue
			}

			return buildModelForBastionMachine(bastionCandidate, buildTagsMap(bastionCandidate), context)
		}
	}

	return model.NoBastion
}

func buildModelForBastionMachine(bastion types.Instance, tags map[string]string, context buildModelContext) model.BastionMachine {
	return model.BastionMachine{
		Url:     *bastion.PublicIpAddress,
		User:    findUserName(bastion, tags, context),
		Keyfile: findKeyFile(bastion, context),
	}
}

func containsBastionInName(instance types.Instance) bool {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			if strings.Contains(strings.ToLower(*tag.Value), "bastion") {
				return true
			}
		}
	}

	return false
}

func getFirstTag(tagNames []string, tags map[string]string) string {
	for _, tag := range tagNames {
		if _, found := tags[tag]; found {
			return tags[tag]
		}
	}

	return ""
}

func buildKeyfile(keyname string, context buildModelContext) string {
	return path.Join(context.config.KeysDirectory, keyname+".pem")
}
