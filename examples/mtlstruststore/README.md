# Examples

This directory contains a basic mTLS Truststore workflow, including setting up CPS and PAPI integrations. The resources used in these examples are available to all users. 
However, if any of the write examples do not work for you, contact your account administrator about your privilege level.

## Run

To run the files, follow the general steps below.
Refer to each example file for more details.

1. Specify the location of your `.edgerc` file and the section header for the set of credentials you want to use. The default section is `default`.
2. Make any necessary changes to the attribute values, replacing placeholder data with your valid data to match your account privileges or requirements.
3. Open a terminal or shell and initialize the provider by running `terraform init`.
4. Run `terraform plan` to preview the changes and `terraform apply` to apply your changes.

## Sample files

Each example file contains calls to the Certificate Provisioning System (CPS) subprovider, Property API (PAPI) subprovider, and mTLS Truststore subprovider endpoints. See the [CPS Terraform integration](https://techdocs.akamai.com/terraform/docs/cps-integration-guide), [PAPI Terraform integration](https://techdocs.akamai.com/terraform/docs/set-up-property-provisioning), and [mTLS Truststore Terraform integration](https://techdocs.akamai.com/terraform/docs/manage-certificate-authority-sets) documentation for complete guides.

| Asset                                   | Description                                                                                                                                                                                                                                                                                                                                 |
|--------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [mTLS Truststore](./mtls.tf)               | Creates a self-signed certificate used to create and activate a CA set on `STAGING` and `PRODUCTION` environments.                                                                                                                                                                                                                              |
| [CPS](./cps.tf)                            | Creates a third-party enrollment with `client_mutual_authentication` enabled, referencing the activated CA set and fetching its CSRs. After self-signing the selected certificate signing request (CSR), it uses the `akamai_cps_upload_certificate` resource to deploy the signed certificate.      |
| [Rules](./rules.tf)                        | Creates property rules with the `enforce_mtls_settings` behavior, referring to the activated CA set.                                                                                                                                                                                                                                        |
| [PAPI](./papi.tf)                          | Creates an edge hostname, CP code, and property, and activates that property on `STAGING` and `PRODUCTION` environments, enforcing the mTLS Truststore configuration.                                                                                                                                                                                                    |
