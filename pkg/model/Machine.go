package model

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

const NoMachineName = "Unknown"
const NoUserName = "Unknown"

var NoMachine = Machine{}
var NoBastion = BastionMachine{}

type Machine struct {
	Id      string
	Name    string
	User    string
	Keyfile string
	Bastion BastionMachine
}

type BastionMachine struct {
	Url     string
	User    string
	Keyfile string
}

func SerializeMachine(machine Machine) string {
	json, err := json.Marshal(machine)

	if err != nil {
		log.Printf("Error serializing machine %v", machine)
		return ""
	}

	return base64.StdEncoding.EncodeToString(json)
}

func DeserializeMachine(serializedMachine string) (Machine, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(serializedMachine)

	if err != nil {
		log.Printf("Error decoding machine %v", serializedMachine)
		return NoMachine, err
	}

	machine := Machine{}

	if err := json.Unmarshal(jsonBytes, &machine); err != nil {
		log.Printf("Error unmarshalling machine %v", jsonBytes)
		return NoMachine, err
	}

	return machine, nil
}
