---
layout: "akamai"
page_title: "Akamai: Enrollments"
subcategory: "Certificate Provisioning System"
description: |-
  Certificate Provisioning System enrollments policies.
---

# akamai_cps_enrollments

Use the `akamai_cps_enrollments` data source to return data for all of a specific contract's enrollments. 

## Basic usage

This example shows how to set up a user:
```hcl
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc = "../config/edgerc"
  config_section = "shared_dns"
}

data "akamai_cps_enrollments" "test_enrollments_list" {
  contract_id = var.contract_id
}

output "dv_output" {
  value = data.akamai_cps_enrollments.test_enrollments_list
}
```


## Argument reference

This data source supports this argument:

* `contract_id` - (Required) A contract's ID, optionally with the `ctr_` prefix.

## Attributes reference

This data source returns these attributes:

* `enrollments` 
  * `enrollment_id` 
  * `common_name` - The fully qualified domain name (FQDN) used for the certificate. 
  * `sans` - Additional common names in a Subject Alternative Names (SAN) list.
  * `secure_network` - The type of deployment network used. `standard-tls` deploys your certificate to Akamai's standard secure network, but it isn't PCI compliant. `enhanced-tls` deploys your certificate to Akamai's more secure network with PCI compliance capability.
  * `sni_only` - Whether you enabled SNI-only extension for the enrollment. Server Name Indication (SNI) is an extension of the Transport Layer Security (TLS) networking protocol. It allows a server to present multiple certificates on the same IP address. All modern web browsers support the SNI extension. If you have the same SAN on two or more certificates with the SNI-only option set, Akamai may serve traffic using any certificate which matches the requested SNI hostname.
  * `admin_contact` - Contact information for the certificate administrator at your company.
  * `certificate_chain_type` - Certificate trust chain type.
  * `csr` - When you create an enrollment, you also generate a certificate signing request (CSR) using CPS. CPS signs the CSR with the private key. The CSR contains all the information the CA needs to issue your certificate.
    * `country_code` - The country code for the country where your organization is located.
    * `city` - The city where your organization resides.
    * `organization` - The name of your company or organization.
    * `organizational_unit` - Your organizational unit.
    * `state` - Your state or province.
  * `enable_multi_stacked_certificates` - If present, an ECDSA certificate is enabled in addition to an RSA certificate. CPS automatically performs all certificate operations on both certificates, and uses the best certificate for each client connection to your secure properties. 
  * `network_configuration` - The network information and TLS Metadata you want CPS to use to push the completed certificate to the network.
    * `client_mutual_authentication` - If present, shows the configuration for client mutual authentication. Specifies the trust chain that is used to verify client certificates and some configuration options.
      * `send_ca_list_to_client` - If present, the server is enabled to send the certificate authority (CA) list to the client.
      * `ocsp_enabled` - If present, the Online Certificate Status Protocol (OCSP) stapling is enabled for client certificates.
      * `set_id` - The identifier of the set of trust chains, created in [Trust Chain Manager](https://techdocs.akamai.com/trust-chain-mgr/docs/welcome-trust-chain-manager).
    * `disallowed_tls_versions` - The TLS protocol version that is not trusted. CPS uses the TLS protocols that Akamai currently supports as a best practice.
    * `clone_dns_names` - If present, CPS directs traffic using all the SANs listed in the SANs parameter when the enrollment was created.
    * `geography` - A list of where you can deploy the certificate. Either `core` to specify worldwide deployment (including China and Russia), `china+core` to specify worldwide deployment and China, or `russia+core` to specify worldwide deployment and Russia. 
    * `must_have_ciphers` - If present, shows ciphers included for enrollment when deployed on the network. The default is `ak-akamai-default` when it is not set. For more information on cipher profiles, see [Akamai community](https://community.akamai.com/customers/s/article/SSL-TLS-Cipher-Profiles-for-Akamai-Secure-CDNrxdxm).
    * `ocsp_stapling` - If present, the enrollment is using OCSP stapling. OCSP stapling improves performance by including a valid OCSP response in every TLS handshake. This option allows the visitors on your site to query the Online Certificate Status Protocol (OCSP) server at regular intervals to obtain a signed time-stamped OCSP response. Possible values are `on`, `off`, or `not-set`.
    * `preferred_ciphers` - If present, shows the ciphers that you prefer to include for the enrollment while deploying it on the network. The default is `ak-akamai-default` when its not set. For more information on cipher profiles, see [Akamai community](https://community.akamai.com/customers/s/article/SSL-TLS-Cipher-Profiles-for-Akamai-Secure-CDNrxdxm).
    * `quic_enabled` - If present, uses the QUIC transport layer network protocol.
  * `signature_algorithm` - If present, shows the Secure Hash Algorithm (SHA) function, either `SHA-1` or `SHA-256`.
  * `tech_contact` - The technical contact within Akamai. This is the person you work closest with at Akamai and who can verify the certificate request. The CA calls this contact if there are any issues with the certificate and they can't reach the `admin_contact`.
  * `organization` - The name of the organization in Akamai where your technical contact works.
    * `name` - The name of the technical contact at Akamai.
    * `phone` - The phone number of the technical contact at Akamai.
    * `address_line_one` - The address for the technical contact at Akamai.
    * `address_line_two` - The address for the technical contact at Akamai.
    * `city` - The address for the technical contact at Akamai.
    * `region` - The region for the technical contact at Akamai.
    * `postal_code` - The postal code for the technical contact at Akamai.
    * `country_code` - The country code for the technical contact at Akamai.
  * `certificate_type` - Populates automatically with the `san` certificate type and is preserved in the `state` file.
  * `validation_type` - Populates automatically with the `dv` validation type and is preserved in the `state` file.
  * `registration_authority` - Populates automatically with the `lets-encrypt` certificate type and is preserved in the `state` file.
  * `pending_changes` - If `true`, there are changes currently pending in CPS. To view pending changes, use the `data_akamai_cps_enrollment` data source.