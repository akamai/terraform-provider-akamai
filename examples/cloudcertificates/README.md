# Examples

This directory contains a basic Cloud Certificate Manager (CCM) workflow.
The resources used in these examples are available to all users. 
However, if any of the write examples do not work for you, contact your account administrator about your privilege level.

## Run

To run the files, follow the general steps below.
Refer to each example file for more details.

1. Specify the location of your `.edgerc` file and the section header for the set of credentials you want to use. The default section is `default`.
2. Make any necessary changes to the attribute values, replacing placeholder data with your valid data to match your account privileges or requirements.
3. Open a terminal or shell and initialize the provider by running `terraform init`.
4. Run `terraform plan` to preview the changes and `terraform apply` to apply your changes.

## Sample files

Each example file contains calls to the Cloud Certificate Manager (CCM) subprovider and Property API (PAPI) subprovider endpoints. See the [PAPI Terraform integration](https://techdocs.akamai.com/terraform/docs/set-up-property-provisioning) and [CCM Terraform integration](https://techdocs.akamai.com/terraform/docs/ccm-integration-guide) documentation for complete guides.

| Asset                                   | Description                                                                                                                                                                                                                                                                                                                                 |
|--------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [Cloud Certificate](./cloudcertificate.tf)               | Creates a third-party cloud certificate and uploads a self-signed certificate.                            |
| [PAPI](./papi.tf)                            | Creates a property with a hostname pointing to the cloud certificate and activates the configuration on the `STAGING` and `PRODUCTION` environments.      |
| [Rules](./property-snippets/main.json)                            | Contains sample rules for the property.      |
