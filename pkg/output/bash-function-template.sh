function {{ .FunctionName }}() {
	{{ .AwsbasshExec }} connect --profile "{{ .AwsProfile }}" \
		--machine-data "{{ .MachineData }}" \
		{{ if .ForceBastion }} --force-bastion {{ end }} \
		"$@"
}

