# Examples

This directory contains basic Rapid rules examples.

## Run

To run any of the files, follow the general steps described below.
Go to each example file for more detailed instructions.

1. Specify the location of your `.edgerc` file and the section header for the set of credentials you'd like to use. The default is `default`.
2. Perform any needed changes to the attribute values, replacing dummy data with your valid data to match your account privileges or needs.
3. Open a Terminal or shell instance and initialize the provider by running `terraform init`.
4. Run `terraform plan` to preview the changes and `terraform apply` to apply your changes.

## Sample files

| Resource                                    | Description                                                                                                                                                                                                                                                                                  |
|---------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [Terraform configuration](./main.tf)        | The `akamai_appsec_rapid_rules` resource enables and configures rapid rules. The corresponding data source returns information about a name, action, action lock, attack group, exceptions for your rapid rules as well as the default action for new rapid rules and rapid ruleset status.  |
| [Rule definitions](./rule_definitions.json) | A JSON file containing a rapid rule's definition, including a rule ID, rule action, action lock, and optionally exceptions for that rapid rule. |
