package main

import (
	"aws-bassh/pkg/connect"
	"aws-bassh/pkg/ec2client"
	"aws-bassh/pkg/loader"
	"aws-bassh/pkg/model"
	"aws-bassh/pkg/output"
	"log"
	"os"
)

func initialize(awsProfile string) bool {
	return ec2client.Initialize(awsProfile) == nil
}

func runGenerate() bool {
	generateConfig := model.MakeCommandLineGenerateConfig()

	if !initialize(generateConfig.AwsProfile) {
		return false
	}

	log.Printf("Generate config %+v\n", generateConfig)

	machines, err := loader.LoadAllMachines(generateConfig)

	if err != nil {
		return false
	}

	return output.WriteMachines(generateConfig, machines)
}

func runConnect() bool {
	connectConfig := model.MakeCommandLineConnectConfig()

	if !initialize(connectConfig.AwsProfile) {
		return false
	}

	log.Printf("Connect config %+v\n", connectConfig)

	return connect.SSH(connectConfig)
}

func main() {
	if len(os.Args) < 2 {
		log.Printf("expected 'generate' or 'connect' subcommands")
		os.Exit(1)
	}

	success := false

	switch os.Args[1] {
	case "generate":
		success = runGenerate()
	case "connect":
		success = runConnect()
	default:
		log.Printf("expected 'generate' or 'connect' subcommands")
	}

	if !success {
		os.Exit(1)
	}
}
