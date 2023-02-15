# zone2hcl

Traspose AWS Route53 zones and records to Terraform HCL resource definitions.

## Usage

```bash
zone2hcl - Transform AWS Route53 Zones and RecordsSet to Terraform HCL

Usage:
  zone2hcl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  records     Generate a zones records set Terraform resources
  version     Print the version number of zone2hcl
  zone        Generate hosted zone Terraform resource
  zones       Generate hosted zones Terraform resource

Flags:
      --config string   config file (default is $HOME/.zone2hcl.yaml)
  -h, --help            help for zone2hcl

Use "zone2hcl [command] --help" for more information about a command.
```