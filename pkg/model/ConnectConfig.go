package model

import (
	"flag"
	"os"
	"strings"
)

var (
	connectCmd = flag.NewFlagSet("connect", flag.ExitOnError)

	awsProfileConnectParam = connectCmd.String("profile", "", "AWS Cli Profile to use")
	machineDataParam       = connectCmd.String("machine-data", "", "Base64 serialized machine information")
	forceBastionParam      = connectCmd.Bool("force-bastion", false, "Force connection via bastion, even if Public Ip available")
	usePublicDnsParam      = connectCmd.Bool("use-public-dns", false, "Use public dns instead of public ip")
	sshParamsParam         = connectCmd.String("ssh-params", "-o StrictHostKeyChecking=no -q", "Extra ssh parameters")
	sshCommandsParam       = connectCmd.String("ssh-commands", "", "SSH Commands to run after the ssh connection is established")
	sshUserNameParam       = connectCmd.String("ssh-user", "", "Use this ssh user for connection")
)

type ConnectConfig struct {
	AwsProfile     string
	Machine        Machine
	ForceBastion   bool
	UsePublicDns   bool
	ExtraSSHParams []string
	SSHCommands    []string
	SSHUserName    string
}

func MakeCommandLineConnectConfig() ConnectConfig {
	connectCmd.Parse(os.Args[2:])

	return ConnectConfig{
		AwsProfile:     getAwsConnectProfile(),
		Machine:        getMachine(*machineDataParam),
		ForceBastion:   *forceBastionParam,
		UsePublicDns:   *usePublicDnsParam,
		ExtraSSHParams: strings.Split(*sshParamsParam, " "),
		SSHUserName:    *sshUserNameParam,
		SSHCommands:    strings.Split(*sshCommandsParam, " "),
	}
}

func getAwsConnectProfile() string {
	profile := *awsProfileConnectParam

	if profile != "" {
		os.Setenv("AWS_PROFILE", profile)
		return profile
	}

	profile = os.Getenv("AWS_PROFILE")

	if profile != "" {
		return profile
	}

	return ""
}

func getMachine(serializedMachine string) Machine {
	if *machineDataParam == "" {
		return NoMachine
	}

	machine, err := DeserializeMachine(*machineDataParam)

	if err != nil {
		return NoMachine
	}

	return machine
}
