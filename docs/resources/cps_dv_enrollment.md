---
layout: akamai
subcategory: Certificate Provisioning System
---

# akamai_cps_dv_enrollment

Use the `akamai_cps_dv_enrollment` resource to create an enrollment for a Domain Validated (DV) certificate. This resource includes all information about your certificate life cycle, from the time you request it, through removal or automatic renewal. You can treat an enrollment as a core container for all the operations you perform within CPS.

You can use this resource with [`akamai_dns_record`](../resources/dns_record.md) or other third-party DNS provider to create DNS records, and [`akamai_cps_dv_validation`](../resources/cps_dv_validation.md) to complete the certificate validation.

<blockquote style="border-left-style: solid; border-left-color: #5bc0de; border-width: 0.25em; padding: 1.33rem; background-color: #e3edf2;"><img src="https://techdocs.akamai.com/terraform-images/img/note.svg" style="float:left; display:inline;" /><div style="overflow:auto;">If you need to enroll a third-party certificate, use the <code><a href="https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/cps_third_party_enrollment">akamai_cps_third_party_enrollment</a></code> resource.</div></blockquote>


## Example usage

Basic usage:

```hcl
resource "akamai_cps_dv_enrollment" "example" {
  contract_id = "ctr_1-AB123"
  acknowledge_pre_verification_warnings = true
  common_name = "cps-test.example.net"
  sans = ["san1.cps-test.example.net","san2.cps-test.example.net"]
  secure_network = "enhanced-tls"
  sni_only = true
  admin_contact {
    first_name = "x1"
    last_name = "x2"
    phone = "123123123"
    email = "x1x2@example.net"
    address_line_one = "150 Broadway"
    city = "Cambridge"
    country_code = "US"
    organization = "Akamai"
    postal_code = "02142"
    region = "MA"
    title = "Administrator"
  }
  tech_contact {
    first_name = "x3"
    last_name = "x4"
    phone = "123123123"
    email = "x3x4@akamai.com"
    address_line_one = "150 Broadway"
    city = "Cambridge"
    country_code = "US"
    organization = "Akamai"
    postal_code = "02142"
    region = "MA"
    title = "Administrator"
  }
  certificate_chain_type = "default"
  csr {
    country_code = "US"
    city = "Cambridge"
    organization = "Akamai"
    organizational_unit = "Dev"
    preferred_trust_chain = "intermediate-a"
    state = "MA"
  }
  network_configuration {
    disallowed_tls_versions = ["TLSv1", "TLSv1_1"]
    clone_dns_names = false
    geography = "core"
    ocsp_stapling = "on"
    preferred_ciphers = "ak-akamai-default"
    must_have_ciphers = "ak-akamai-default"
    quic_enabled = false
  }
  signature_algorithm = "SHA-256"
  organization {
    name = "Akamai"
    phone = "123123123"
    address_line_one = "150 Broadway"
    city = "Cambridge"
    country_code = "US"
    postal_code = "02142"
    region = "MA"
  }
}

output "dns_challenges" {
  value = akamai_cps_dv_enrollment.example.dns_challenges
}

output "http_challenges" {
  value = akamai_cps_dv_enrollment.example.http_challenges
}

output "enrollment_id" {
  value = akamai_cps_dv_enrollment.example.id
}
```
## Argument reference

The following arguments are supported:

* `contract_id` - (Required) A contract's ID, optionally with the `ctr_` prefix.
* `common_name` - (Required) The fully qualified domain name (FQDN) for which you plan to use your certificate. The domain name you specify here must be owned or have legal rights to use the domain by the company you specify as `organization`. The company that owns the domain name must be a legally incorporated entity and be active and in good standing.
* `allow_duplicate_common_name` - (Optional) Boolean. Set to `true` if you want to reuse a common name that's part of an existing enrollment. 
* `sans` - (Optional) Additional common names to create a Subject Alternative Names (SAN) list.
* `secure_network` - (Required) The type of deployment network you want to use. `standard-tls` deploys your certificate to Akamai's standard secure network, but it isn't PCI compliant. `enhanced-tls` deploys your certificate to Akamai's more secure network with PCI compliance capability.
* `sni_only` - (Required) Whether you want to enable SNI-only extension for the enrollment. Server Name Indication (SNI) is an extension of the Transport Layer Security (TLS) networking protocol. It allows a server to present multiple certificates on the same IP address. All modern web browsers support the SNI extension. If you have the same SAN on two or more certificates with the SNI-only option set, Akamai may serve traffic using any certificate which matches the requested SNI hostname. You should avoid multiple certificates with overlapping SAN names when using SNI-only. You can't change this setting once an enrollment is created.
* `acknowledge_pre_verification_warnings` - (Optional) Whether you want to automatically acknowledge the validation warnings of the current job state and proceed with the execution of a change.
* `admin_contact` - (Required) Contact information for the certificate administrator at your company.

    Requires these additional arguments:

      * `first_name` - (Required) The first name of the certificate administrator at your company.
      * `last_name` - (Required) The last name of the certificate administrator at your company.
      * `title` - (Optional) The title of the certificate administrator at your company.
      * `organization` - (Required) The name of your organization.
      * `email` - (Required) The email address of the administrator who you want to use as a contact at your company.
      * `phone` - (Required) The phone number of your organization.
      * `address_line_one` - (Required) The address of your organization.
      * `address_line_two` - (Optional) The address of your organization.
      * `city` - (Required) The city where your organization resides.
      * `region` - (Required) The region of your organization, typically a state or province.
      * `postal_code` - (Required) The postal code of your organization.
      * `country_code` - (Required) The code for the country where your organization resides.
* `certificate_chain_type` - (Optional) Certificate trust chain type.
* `csr` - (Required) 	When you create an enrollment, you also generate a certificate signing request (CSR) using CPS. CPS signs the CSR with the private key. The CSR contains all the information the CA needs to issue your certificate.

    Requires these additional arguments:

      * `country_code` - (Required) The country code for the country where your organization is located.
      * `city` - (Required) The city where your organization resides.
      * `organization` - (Required The name of your company or organization. Enter the name as it appears in all legal documents and as it appears in the legal entity filing.
      * `organizational_unit` - (Required) Your organizational unit.
      * `preferred_trust_chain` - (Optional) The preferred trust chain will be included by CPS with the leaf certificate in the TLS handshake. If the field does not have a value, whichever trust chain Akamai chooses will be used by default.
      * `state` - (Required) 	Your state or province.
* `network_configuration` - (Required) The network information and TLS Metadata you want CPS to use to push the completed certificate to the network.

    Requires these additional arguments:

      * `client_mutual_authentication` - (Optional) The configuration for client mutual authentication. Specifies the trust chain that is used to verify client certificates and some configuration options.

        Requires these additional arguments:

         * `send_ca_list_to_client` - (Optional) Whether you want to enable the server to send the certificate authority (CA) list to the client.
         * `ocsp_enabled` - (Optional) Whether you want to enable the Online Certificate Status Protocol (OCSP) stapling for client certificates.
         * `set_id` - (Optional) The identifier of the set of trust chains, created in [Trust Chain Manager](https://techdocs.akamai.com/trust-chain-mgr/docs/welcome-trust-chain-manager).
      * `disallowed_tls_versions` - (Optional) The TLS protocol version to disallow. CPS uses the TLS protocols that Akamai currently supports as a best practice.
      * `clone_dns_names` - (Optional) Whether CPS should direct traffic using all the SANs you listed in the SANs parameter when you created your enrollment.
      * `geography` - (Required) Lists where you can deploy the certificate. Either `core` to specify worldwide deployment (including China and Russia), `china+core` to specify worldwide deployment and China, or `russia+core` to specify worldwide deployment and Russia. You can only use the setting to include China and Russia if your Akamai contract specifies your ability to do so and you have approval from the Chinese and Russian government.
      * `must_have_ciphers` - (Optional) The ciphers to include for the enrollment while deploying it on the network. Defaults to `ak-akamai-default` when it is not set. For more information on cipher profiles, see [Akamai community](https://community.akamai.com/customers/s/article/SSL-TLS-Cipher-Profiles-for-Akamai-Secure-CDNrxdxm).
      * `ocsp_stapling` - (Optional) Whether to use OCSP stapling for the enrollment, either `on`, `off` or `not-set`. OCSP Stapling improves performance by including a valid OCSP response in every TLS handshake. This option allows the visitors on your site to query the Online Certificate Status Protocol (OCSP) server at regular intervals to obtain a signed time-stamped OCSP response. This response must be signed by the CA, not the server, therefore ensuring security. Disable OSCP Stapling if you want visitors to your site to contact the CA directly for an OSCP response. OCSP allows you to obtain the revocation status of a certificate.
      * `preferred_ciphers` - (Optional) Ciphers that you preferably want to include for the enrollment while deploying it on the network. Defaults to `ak-akamai-default` when it is not set. For more information on cipher profiles, see [Akamai community](https://community.akamai.com/customers/s/article/SSL-TLS-Cipher-Profiles-for-Akamai-Secure-CDNrxdxm).
      * `quic_enabled` - (Optional) Whether to use the QUIC transport layer network protocol.
* `signature_algorithm` - (Required) The Secure Hash Algorithm (SHA) function, either `SHA-1` or `SHA-256`.
* `tech_contact` - (Required) The technical contact within Akamai. This is the person you work closest with at Akamai and who can verify the certificate request. The CA calls this contact if there are any issues with the certificate and they can't reach the `admin_contact`.

    Requires these additional arguments:

      * `first_name` - (Required) The first name of the technical contact at Akamai.
      * `last_name` - (Required) The last name of the technical contact at Akamai.
      * `title` - (Optional) The title of the technical contact at Akamai.
      * `organization` - (Required) The name of the organization in Akamai where your technical contact works.
      * `email` - (Required) The email address of the technical contact at Akamai, accessible at the `akamai.com` domain.
      * `phone` - (Required) The phone number of the technical contact at Akamai.
      * `address_line_one` - (Required) The address for the technical contact at Akamai.
      * `address_line_two` - (Optional) The address for the technical contact at Akamai.
      * `city` - (Required) The address for the technical contact at Akamai.
      * `region` - (Required) The region for the technical contact at Akamai.
      * `postal_code` - (Required) The postal code for the technical contact at Akamai.
      * `country_code` - (Required) The country code for the technical contact at Akamai.
* `organization` - (Required) Your organization information.

    Requires these additional arguments:

      * `name` - (Required) The name of your organization.
      * `phone` - (Required) The phone number of the administrator who you want to use as a contact at your company.
      * `address_line_one` - (Required) The address of your organization.
      * `address_line_two` - (Optional) The address of your organization.
      * `city` - (Required) The city where your organization resides.
      * `region` - (Required) The region of your organization, typically a state or province.
      * `postal_code` - (Required) The postal code of your organization.
      * `country_code` - (Required) The code for the country where your organization resides.

### Deprecated arguments

* `enable_multi_stacked_certificates` - (Deprecated) Whether to enable an ECDSA certificate in addition to an RSA certificate. CPS automatically performs all certificate operations on both certificates, and uses the best certificate for each client connection to your secure properties. If you are pinning the certificates, you need to pin both the RSA and the ECDSA certificate.

## Attributes reference

The resource returns these attributes:

* `registration_authority` - (Required) This value populates automatically with the `lets-encrypt` certificate type and is preserved in the `state` file.
* `certificate_type` - (Required) This value populates automatically with the `san` certificate type and is preserved in the `state` file.
* `validation_type` - (Required) This value populates automatically with the `dv` validation type and is preserved in the `state` file.
* `id` - The unique identifier for this enrollment.
* `dns_challenges` - The validation challenge for the domains listed in the certificate. To successfully perform the validation, only one challenge for each domain must be complete, either `dns_challenges` or `http_challenges`.

    Returns these additional attributes:

      * `domain` - The domain to validate.
      * `full_path` - The URL where Akamai publishes `response_body` for Let's Encrypt to validate.
      * `response_body` - The data Let's Encrypt expects to find served at `full_path` URL.
* `http_challenges` - The validation challenge for the domains listed in the certificate. To successfully perform the validation, only one challenge for each domain must be complete, either `dns_challenges` or `http_challenges`.

    Returns these additional attributes:

      * `domain` - The domain to validate.
      * `full_path` - The URL where Akamai publishes `response_body` for Let's Encrypt to validate.
      * `response_body` - The data Let's Encrypt expects to find served at `full_path` URL.

## Import

Basic Usage:

```hcl
resource "akamai_cps_dv_enrollment" "example" {
# (resource arguments)
}
```

You can import your Akamai DV enrollment using a comma-delimited string of the enrollment ID and  
contract ID, optionally with the `ctr_` prefix. You have to enter the IDs in this order:

`enrollment_id,contract_id`

For example:

```shell
$ terraform import akamai_cps_dv_enrollment.example 12345,1-AB123
```
