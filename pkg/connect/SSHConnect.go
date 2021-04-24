package connect

import (
	"aws-bassh/pkg/ec2client"
	"aws-bassh/pkg/model"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
	"os"
	"os/exec"
	"strings"
)

func SSH(config model.ConnectConfig) {
	instance, err := ec2client.DescribeInstance(config.Machine.Id)

	if err != nil {
		return
	}

	validateAndConnectToInstance(config, instance)
}

func validateAndConnectToInstance(config model.ConnectConfig, instance *types.Instance) {
	if !validateBeforeConnect(config, instance) {
		return
	}

	connectToInstance(config, instance)
}

func validateBeforeConnect(config model.ConnectConfig, instance *types.Instance) bool {
	if !checkMachineState(instance) {
		return false
	}

	if !validateKeyfile(config.Machine.Keyfile) {
		return false
	}

	if shouldUseBastion(config, instance) && !validateKeyfile(config.Machine.Bastion.Keyfile) {
		return false
	}

	return true
}

func checkMachineState(instance *types.Instance) bool {
	if instance.State.Name == types.InstanceStateNameRunning {
		return true
	}

	log.Printf("Invalid machine state: %v", instance.State.Name)
	return false
}

func validateKeyfile(keyfile string) bool {
	if _, err := os.Stat(keyfile); os.IsNotExist(err) {
		log.Printf("")
		log.Printf("Ssh private key is missing %v", keyfile)
		log.Printf("Please do the following steps in order to fix it:")
		log.Printf("  1. Run the 'generate' command again")
		log.Printf("  2. Specify the --keys parameter, it should point to a directory containing the ssh private key")
		log.Printf("  3. Try to connect again")
		log.Printf("")
		return false
	}

	return true
}

func connectToInstance(config model.ConnectConfig, instance *types.Instance) {
	var sshArgs = buildSshArgs(config, instance)

	logSshCommand(sshArgs)
	spawnSsh(sshArgs)
}

func buildSshArgs(config model.ConnectConfig, instance *types.Instance) []string {
	args := []string{}

	args = append(args, "-i", config.Machine.Keyfile)
	args = append(args, config.ExtraSSHParams...)

	if shouldUseBastion(config, instance) {
		args = append(args, buildBastionArgs(config, instance)...)
	}

	args = append(args, buildUserAddressArg(config, instance))

	return args
}

func shouldUseBastion(config model.ConnectConfig, instance *types.Instance) bool {
	return getPublicAddress(config, instance) == nil ||
		(config.ForceBastion && config.Machine.Bastion.Url != "")
}

func getPublicAddress(config model.ConnectConfig, instance *types.Instance) *string {
	if config.UsePublicDns {
		return instance.PublicDnsName
	} else {
		return instance.PublicIpAddress
	}
}

func buildUserAddressArg(config model.ConnectConfig, instance *types.Instance) string {
	if shouldUseBastion(config, instance) {
		return config.Machine.User + "@" + *instance.PrivateIpAddress
	} else {
		return config.Machine.User + "@" + *getPublicAddress(config, instance)
	}
}

func buildBastionArgs(config model.ConnectConfig, instance *types.Instance) []string {
	bastionArgs := []string{}

	bastionArgs = append(bastionArgs, "-o")
	bastionArgs = append(bastionArgs, generateBastionProxyCommand(config))

	return bastionArgs
}

func generateBastionProxyCommand(config model.ConnectConfig) string {
	return fmt.Sprintf("proxycommand ssh %s -W %s -f -i %s %s@%s",
		strings.Join(config.ExtraSSHParams, " "),
		"%h:%p",
		config.Machine.Bastion.Keyfile,
		config.Machine.Bastion.User,
		config.Machine.Bastion.Url,
	)
}

func logSshCommand(sshArgs []string) {
	var sb strings.Builder

	sb.WriteString(getSshExec())
	sb.WriteString(" ")

	for _, arg := range sshArgs {
		sb.WriteString("\"")
		sb.WriteString(arg)
		sb.WriteString("\" ")
	}

	sb.WriteString("\n")
	fmt.Println("")
	fmt.Println("SSH Command:")
	fmt.Println("")
	fmt.Println(sb.String())
	fmt.Println("")
}

func spawnSsh(sshArgs []string) {
	cmd := exec.Command(getSshExec(), sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func getSshExec() string {
	return "ssh"
}
