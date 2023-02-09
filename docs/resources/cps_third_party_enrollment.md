---
layout: akamai
subcategory: Certificate Provisioning System
---

# akamai_cps_third_party_enrollment

Use the `akamai_cps_third_party_enrollment` resource to create an enrollment for a third-party certificate. As with Domain Validated (DV) certificate enrollments, you can treat a third-party enrollment as a core container for all the operations you perform within CPS.

You can use this resource with:

* [`akamai_dns_record`](../resources/dns_record.md) or other third-party DNS provider to create DNS records
* [`akamai_cps_upload_certificate`](../resources/cps_upload_certificate.md) to complete the validation and activate the certificate on the staging and production networks. Set the `change_management` argument in this resource to `true` if you want to test and view the certificate on the staging network before deploying it to production.

<blockquote style="border-left-style: solid; border-left-color: #5bc0de; border-width: 0.25em; padding: 1.33rem; background-color: #e3edf2;"><img src="https://techdocs.akamai.com/terraform-images/img/note.svg" style="float:left; display:inline;" /><div style="overflow:auto;">If you need to enroll a DV certificate, use the <code><a href="https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/cps_dv_enrollment">akamai_cps_dv_enrollment</a></code> resource.</div></blockquote>

## Example usage

Basic usage:

```hcl

resource "akamai_cps_third_party_enrollment" "enrollment" {
  contract_id           = "C-0N7RAC7"
  common_name           = "*.example.com"
  secure_network        = "enhanced-tls"
  sni_only              = true
  auto_approve_warnings = [
    "DNS_NAME_LONGER_THEN_255_CHARS",
    "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS",
    "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"
  ]
  signature_algorithm   = "SHA-256"
  admin_contact {
    first_name       = "Mario"
    last_name        = "Rossi"
    phone            = "+1-311-555-2368"
    email            = "mrossi@example.com"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    organization     = "Example Corp."
    postal_code      = "02142"
    region           = "MA"
    title            = "Administrator"
  }
  tech_contact {
    first_name       = "Juan"
    last_name        = "Perez"
    phone            = "+1-311-555-2369"
    email            = "jperez@example.com"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    organization     = "Example Corp."
    postal_code      = "02142"
    region           = "MA"
    title            = "Administrator"
  }
  csr {
    country_code        = "US"
    city                = "Cambridge"
    organization        = "Example Corp."
    organizational_unit = "Corp IT"
    state               = "MA"
  }
  network_configuration {
    disallowed_tls_versions = ["TLSv1", "TLSv1_1"]
    clone_dns_names         = false
    geography               = "core"
    ocsp_stapling           = "on"
    preferred_ciphers       = "ak-akamai-2020q1"
    must_have_ciphers       = "ak-akamai-2020q1"
    quic_enabled            = false
  }
  organization {
    name             = "Example Corp."
    phone            = "+1-311-555-2370"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    country_code     = "US"
    postal_code      = "02142"
    region           = "MA"
  }
}

output "enrollment_id" {
  value = akamai_cps_third_party_enrollment.enrollment.id
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
* `acknowledge_pre_verification_warnings` - (Optional) Whether you want to automatically acknowledge the validation warnings related to the current job state and proceed with the change.
* `auto_approve_warnings` - (Optional) The list of post-verification warning IDs you want to automatically acknowledge. To retrieve the list of warnings, use the `akamai_cps_warnings` data source.
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
      * `must_have_ciphers` - (Optional) The ciphers to include for the enrollment while deploying it on the network. Defaults to `ak-akamai-2020q1` when it is not set. For more information on cipher profiles, see [Akamai community](https://community.akamai.com/customers/s/article/SSL-TLS-Cipher-Profiles-for-Akamai-Secure-CDNrxdxm).
      * `ocsp_stapling` - (Optional) Whether to use OCSP stapling for the enrollment, either `on`, `off` or `not-set`. OCSP Stapling improves performance by including a valid OCSP response in every TLS handshake. This option allows the visitors on your site to query the Online Certificate Status Protocol (OCSP) server at regular intervals to obtain a signed time-stamped OCSP response. This response must be signed by the CA, not the server, therefore ensuring security. Disable OSCP Stapling if you want visitors to your site to contact the CA directly for an OSCP response. OCSP allows you to obtain the revocation status of a certificate.
      * `preferred_ciphers` - (Optional) Ciphers that you preferably want to include for the enrollment while deploying it on the network. Defaults to `ak-akamai-2020q1` when it is not set. For more information on cipher profiles, see [Akamai community](https://community.akamai.com/customers/s/article/SSL-TLS-Cipher-Profiles-for-Akamai-Secure-CDNrxdxm).
      * `quic_enabled` - (Optional) Whether to use the QUIC transport layer network protocol.
* `signature_algorithm` - (Required) The Secure Hash Algorithm (SHA) function, either `SHA-1` or `SHA-256`. If you change this value, you may need to run the `terraform destroy` and `terraform apply` commands.
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
* `change_management` - (Optional) Boolean. Set to `true` to have CPS deploy first to staging for testing purposes. To deploy the certificate to production, use the `acknowledge_change_management` argument in the `akamai_cps_upload_certificate` resource. <br> If you don't use this option, CPS will automatically deploy the certificate to both networks.
* `exclude_sans` - (Optional) If set to `true`, then the SANs in the enrollment don't appear in the CSR that you send to your CA.

## Attributes reference

The resource returns this attribute:

* `id` - The unique identifier for this enrollment.

## Import

Basic Usage:

```hcl
resource "akamai_cps_third_party_enrollment" "example" {
# (resource arguments)
}
```

You can import your Akamai third-party enrollment using a comma-delimited string of the enrollment ID and  
contract ID, optionally with the `ctr_` prefix. You have to enter the IDs in this order:

`enrollment_id,contract_id`

For example:

```shell
$ terraform import akamai_cps_third_party_enrollment.example 12345,1-AB123
```
