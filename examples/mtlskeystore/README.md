# Examples

This directory contains basic mTLS Keystore examples, including `AKAMAI` and `THIRD-PARTY` client certificates integration workflows. 
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

Each example file contains calls to the mTLS Keystore subprovider and Property API (PAPI) subprovider endpoints. See the  [PAPI Terraform integration](https://techdocs.akamai.com/terraform/docs/set-up-property-provisioning) and [mTLS Keystore Terraform integration](https://techdocs.akamai.com/terraform/docs/manage-client-certificates) documentation for complete guides.

| Asset                                   | Description                                                                                                                                                                                                                                                                                                                                 |
|--------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [Akamai-signed client certificate](./akamai/main.tf)               | Creates an `AKAMAI` client certificate and enforces the mTLS Keystore configuration by using the `mtls_origin_keystore` behavior in the property rules configuration.                                                                                                                                                                                                                              |
| [Third-party signed client certificate](./third_party/main.tf)                            | Creates a `THIRD-PARTY` client certificate with a self-signed certificate and enforces the mTLS Keystore configuration by using the `mtls_origin_keystore` behavior in the property rules configuration.      |
