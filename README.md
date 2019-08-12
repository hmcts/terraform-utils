# Terraform utils

Terraform utilities to help with the maintenance of Terraform related code.

Currently available utilities:
- `tfwhitelist`: allows to scan all Terraform resources and modules used in a project and verify that they match a given whitelist.  


## tfwhitelist

To match Terraform infrastructure against a whitelist of resources and modules, the following command can be 
used:

```shell script
tf-utils --whitelist <terraform-infra-dir-path> <whitelist-file-path>...<whitelist-file-path>
```

where the first argument is a Terraform definitions directory and the second argument is a list of whitelist files which are merged as first step. 
A whitelist is a json file containing allowed resources and modules. Having multiple files merged allows to have a global whitelist which can be 
specialised in case that is needed in isolated specific cases.
For example:

```json
{
  "resources": [
    {"type": "azurerm_key_vault_secret"},
    {"type": "azurerm_resource_group"}
  ],
  "module_calls": [
    {"source":  "git@github.com:hmcts/cnp-module-webapp?ref=master"},
    {"source":  "git@github.com:hmcts/cnp-module-postgres?ref=master"}
  ]
}
```

This tool uses the [terraform-config-inspect](https://github.com/hashicorp/terraform-config-inspect) library by Hashicorp. 