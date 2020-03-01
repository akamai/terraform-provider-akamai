# Managing multiple similar environments with Terraform Workspaces

This is an alternative to the [environments-test](../environments-test) example for synchronizing property changes across multiple environments using [Terraform Workspaces](https://www.terraform.io/docs/state/workspaces.html).

The idea is simple: maintain one Terraform configuration file and split your Terraform state across workspaces.
Per-environment overrides are provided by environment-specific `tfvars` files.

## Initial Setup

* create `terraform.tfvars` and update `main.tf` to match your configuration
* for each environment that you plan to have, create `stages/${ENVNAME}.tfvars` (e.g. `stages/dev.tfvars`)
* create your workspaces, e.g.: `for stage in dev qa prod; terraform workspace new $stage; done`

## Operation Example

```bash
# switch to dev
terraform workspace select dev
# update without activating
terraform apply -var-file stages/$(terraform workspace show).tfvars -var staging=false -var production=false
# update and activate on staging
terraform apply -var-file stages/$(terraform workspace show).tfvars -var staging=true
# update and activate on staging & production
terraform apply -var-file stages/$(terraform workspace show).tfvars -var staging=true -var production=true
```
