# Examples

This directory contains basic CPS examples, including DV and third-party enrollment workflows. The CPS resources used in these examples are available to all users. 
But, if you find one of the write examples doesn't work for you, talk with your account's admin about your privilege level.

## Run

To run any of the files, follow the general steps described below.
Go to each example file for more detailed instructions.

1. Specify the location of your `.edgerc` file and the section header for the set of credentials you'd like to use. The default is `default`.
2. Perform any needed changes to the attribute values, replacing dummy data with your valid data to match your account privileges or needs.
3. Open a Terminal or shell instance and initialize the provider by running `terraform init`.
4. Run `terraform plan` to preview the changes and `terraform apply` to apply your changes.

## Sample files

The example in each file contains a call to one of the Certificate Provisioning System (CPS) subprovider endpoints. See the [CPS Terraform integration](https://techdocs.akamai.com/terraform/docs/cps-integration-guide) doc for a complete guide and its resources for each of the call used.

| Resource                                                   | Description                                                                                                                                                                                                                |
|------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [DV enrollment](./dv_enrollment/main.tf)                   | Creates a basic DV enrollment and performs DV validation for that certificate.                                                                                                                                             |
| [Third-party enrollment](./third_party_enrollment/main.tf) | Creates a basic third-party enrollment and fetches its CSRs. After manually signing the selected Certificate Signing Request (CSR), it uses the `akamai_cps_upload_certificate` resource to deploy the signed certificate. |
