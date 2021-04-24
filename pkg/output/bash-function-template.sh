function {{ .FunctionName }}() {
	{{ .AwsbasshExec }} connect --profile "{{ .AwsProfile }}" --machine-data "{{ .MachineData }}" "{{ .ForceBastion }}" "$@"
}

