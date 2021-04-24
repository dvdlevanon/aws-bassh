package model

import (
	"flag"
	"os"
	"strings"
)

var (
	generateCmd = flag.NewFlagSet("generate", flag.ExitOnError)

	awsProfileGenerateParam   = generateCmd.String("profile", "", "AWS Cli Profile to use")
	outputFileParam           = generateCmd.String("output-file", "output.sh", "Bash output file")
	bashFunctionsPrefixParam  = generateCmd.String("prefix", "ec2_", "Bash functions prefix")
	keysDirectoryParam        = generateCmd.String("keys", "keys", "A directory containing pem keys for the machines")
	forceBastionGenerateParam = generateCmd.Bool("force-bastion", false, "Force connection via bastion, even if Public Ip available")
	nameTagsParam             = generateCmd.String("name-tags", "Name", "A comma separated names of tags, for Machine name")
	userTagsParam             = generateCmd.String("user-tags", "SSHUser", "A comma separated names of tags, for SSH user")
	bastionUrlTagsParam       = generateCmd.String("bastion-url-tags", "BastionUrl", "A comma separated names of tags, for Bastion url")
	bastionUserTagsParam      = generateCmd.String("bastion-user-tags", "BastionUser", "A comma separated names of tags, for Bastion user")
	bastionKeysTagsParam      = generateCmd.String("bastion-key-tags", "BastionKey", "A comma separated names of tags, for Bastion ssh key")
)

type GenerateConfig struct {
	AwsProfile      string
	OutputFile      string
	KeysDirectory   string
	BashAliasPrefix string
	ForceBastion    bool

	NameTags           []string
	UserTags           []string
	BastionUrlTags     []string
	BastionUserTags    []string
	BastionKeyNameTags []string
}

func MakeCommandLineGenerateConfig() GenerateConfig {
	generateCmd.Parse(os.Args[2:])

	return GenerateConfig{
		AwsProfile:         getAwsGenerateProfile(),
		OutputFile:         *outputFileParam,
		KeysDirectory:      *keysDirectoryParam,
		BashAliasPrefix:    *bashFunctionsPrefixParam,
		ForceBastion:       *forceBastionGenerateParam,
		NameTags:           strings.Split(*nameTagsParam, ","),
		UserTags:           strings.Split(*userTagsParam, ","),
		BastionUrlTags:     strings.Split(*bastionUrlTagsParam, ","),
		BastionUserTags:    strings.Split(*bastionUserTagsParam, ","),
		BastionKeyNameTags: strings.Split(*bastionKeysTagsParam, ","),
	}
}

func getAwsGenerateProfile() string {
	profile := *awsProfileGenerateParam

	if profile != "" {
		return profile
	}

	profile = os.Getenv("AWS_PROFILE")

	if profile != "" {
		return profile
	}

	return ""
}
