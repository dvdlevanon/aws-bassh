
# AWSBassh

A command line utility which helps to connect AWS Ec2 instances. It get a list of instances from AWS, generate a bash function for each machine and write them into an output file. This file can then be included from .bashrc in order to expose those functions to all bash sessions.

## Installation

```bash
curl -L -o awsbassh https://github.com/dvdlevanon/aws-bassh/releases/download/0.0.1/awsbassh-0.0.1-x68_64
install -t /usr/local/bin awsbassh
```

## Building from source

Install awsbassh from source

```bash 
  git clone https://github.com/dvdlevanon/aws-bassh
  go build awsbassh.go
  ./awsbassh
```

## Usage/Examples

### Prerequisite
This tool uses AWS SDK for accessing Ec2 information, it uses AWS profiles for authentication. Make sure you have at least one working AWS profile before continue.
https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html

Make sure AWS profile is properly configured by running.
```bash
aws --profile <PROFILE_NAME> ec2 describe-instances
```

### Generating bash functions

Generate `output.sh` in the current working directory.
```bash
./awsbassh generate --profile <PROFILE_NAME> --keys <keys_directory>
```
The output file should contain a function per machine, the name of the function is the name of the machine (Name tag).

`<keys_directory>` is a directory containing the ssh private keys of the machines in this profile. The `generate` command takes the keyname as provided from AWS and append it to the `<keys_directory>` parameter.

### Connecting to a machine by its name
```bash
source output.sh
ec2_<type the machine name and press enter>
```

### Manage machines from different AWS profiles.
- Initialize a list of AWS profiles
- For each profile run: `./awsbassh generate --profile <PROFILE_NAME> --keys <keys_directory> --output-file .output/<PROFILE_NAME>.sh --prefix <PROFILE_NAME>_`
- Create a `aws-machines.sh` with a content similar to:
```bash
source .output/<PROFILE_NAME1>.sh 
source .output/<PROFILE_NAME2>.sh 
source .output/<PROFILE_NAME3>.sh 
```
- Add `source aws-machines.sh` to `.bashrc`
- Open a new shell, type the profile name and tab+tab to see a complete list of machines in this profile.

### Supporting Proxy (Bastion, jump server, e.g.)
If your machines are running inside a VPC and they don't have a public IP, or their ssh is blocked for outside connections via security group. Its common practice to use a proxy machine to connect the instances inside the VPC. `awsbassh` support connecting via a proxy server using SSH's `proxycommand`.

There are two ways to configure the proxy server, explicit (recommended) and implicit.

#### Explicit proxy server
Add those tags to the machines inside the VPC (tag names are configurable - see Configuration section)
- BastionKey - The keyname of the bastion server (Make sure it exists in `<keys_directory>`)
- BastionUrl - The IP or DNS of the proxy server
- BastionUser - The user name to use for the proxy server

#### Implicit proxy server (Not recommended)
If proxy tags are missing, and the machine doesn't have a public IP `awsbashh` looks for a machine with those conditions:
- Same VPC as the private machine
- Contains "bastion" in its name
- Have a public IP
If such machine is found, it automatically used as the proxy server for the private machine. Notice there is not way to control its private key or user name, they'll be the same as the private machine.

### Configuration
```bash
./awsbassh generate --help
Usage of generate:
  -bastion-key-tags string
    	A comma separated names of tags, for Bastion ssh key (default "BastionKey")
  -bastion-url-tags string
    	A comma separated names of tags, for Bastion url (default "BastionUrl")
  -bastion-user-tags string
    	A comma separated names of tags, for Bastion user (default "BastionUser")
  -force-bastion
    	Force connection via bastion, even if Public Ip available
  -keys string
    	A directory containing pem keys for the machines (default "keys")
  -name-tags string
    	A comma separated names of tags, for Machine name (default "Name")
  -output-file string
    	Bash output file (default "output.sh")
  -prefix string
    	Bash functions prefix (default "ec2_")
  -profile string
    	AWS Cli Profile to use
  -user-tags string
    	A comma separated names of tags, for SSH user (default "SSHUser")
```

## Authors

- [@dvdlevanon](https://www.github.com/dvdlevanon)


## Contributing

Contributions are always welcome!

## License

[GPL-3](https://choosealicense.com/licenses/gpl-3.0//)
