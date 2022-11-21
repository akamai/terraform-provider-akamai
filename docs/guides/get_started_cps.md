---
layout: akamai
page_title: Guide
subcategory: Certificate Provisioning
---

Increase your customers' trust, encrypt sensitive information, and improve SEO rankings using first- and third-party SSL/TLS Domain Validation certificates that securely deliver content to and from your site.

## Before you begin

* Understand the [basics of Terraform](https://learn.hashicorp.com/terraform?utm_source=terraform_io).
* Complete the steps in [Get started](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started) and [Set up your authentication](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/auth).
* Set up a [secure property](https://techdocs.akamai.com/terraform/docs/pm-rc-property).

## How it works

![CPS workflow process](https://techdocs.akamai.com/terraform-images/certificate-provisioning-system/img/cps-workflow.png)

## Get a contract ID

To identify what contract you'll work against, add the contracts data source to your Akamai configuration to get a list of available contracts associated with your account.

```
data "akamai_contracts" "my-contracts" {
}

output "property_match" {
value = data.akamai_contracts.my-contracts
}
```

## First-party enrollments

### Create a Domain Validated enrollment

<blockquote style="border-left-style: solid; border-left-color: #5bc0de; border-width: 0.25em; padding: 1.33rem; background-color: #e3edf2;"><img src="https://techdocs.akamai.com/terraform-images/img/note.svg" style="float:left; display:inline;" /><div style="overflow:auto;">You can import an existing enrollment(s) and make adjustments using the import command on the <code>akamai_cps_dv_enrollment</code> resource. Add a name your <code>enrollment_id</code> and <code>contract_id</code> as a comma-delimited string of the to the end of the command.<pre>$ terraform import akamai_cps_dv_enrollment.name 12345,A-12345</pre></div></blockquote>

Create an enrollment, the core configuration for CPS, and define your certificate's life cycle.

1. Add your contract ID to the [akamai_cps_dv_enrollment](../resources/cps_dv_enrollment.md) resource and provide information about your certificate.
1. Run `terraform plan` to check your syntax and review the total set of changes you're making and then run `terraform apply` to create your enrollment.

The response will return certificate signing request (CSR) challenges you need to validate control over your domains.

### Validate control over your domains

When you create or modify a DV enrollment, it triggers a certificate signing request (CSR). CPS automatically sends it to [Let's Encrypt](https://letsencrypt.org/) for signing, and Let's Encrypt sends back a challenge for each domain listed on your certificate.

To answer the challenge, prove that you have control over the domains listed in the CSR by redirecting your traffic to Akamai, and then we complete the challenge process for you by detecting the redirect.

* For `http_challenges`, create a file with a token and put it in the designated folder on your site. Once Akamai detects the file is in place, it asks Let's Encrypt to validate the domain.
* For `dns_challenges`, add a `TXT` record to the DNS configuration of your domain. If you're using the [Edge DNS subprovider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/edge-dns), you can create DNS records for the provided SANs from the same `config` file.

<blockquote style="border-left-style: solid; border-left-color: #5bc0de; border-width: 0.25em; padding: 1.33rem; background-color: #e3edf2;"><img src="https://techdocs.akamai.com/terraform-images/img/note.svg" style="float:left; display:inline;" /><div style="overflow:auto;">If the challenge token expires, run <code>terraform-apply</code> again to pull the latest token. Even though Terraform doesn't automatically notify you of any updates to tokens, you can set <a href="https://www.terraform.io/docs/language/values/outputs.html">outputs</a> for <code>dns_challenges</code> and <code>http_challenges</code>. If applicable, <code>terraform-apply</code> returns new values for those arrays.</div></blockquote>

### Send the validation acknowledgement

Once you complete the Let's Encrypt challenges, let CPS know your tokens are ready for validation.

1. Use the `akamai_cps_dv_validation` resource to send the acknowledgement to CPS.

    ```hcl
      resource "akamai_cps_dv_validation" "example" {
      enrollment_id = akamai_cps_dv_enrollment.example.id
      sans = akamai_cps_dv_enrollment.example.sans
      }
    ```

1. Run `terraform plan` to check your syntax and review the total set of changes you're making and then run `terraform apply` to apply your changes to your infrastructure.

## Third-party enrollments

### Create a third-party enrollment

To create a third-party enrollment, use your contract ID to define your certificate's life cycle and create your enrollment with the [akamai_cps_third_party_enrollment](../resources/cps_third_party_enrollment.md) resource.

Once you've set up the resource, run `terraform plan` to check your syntax and review the total set of changes you're making.

Run `terraform apply` to apply your changes to your infrastructure and begin domain validation.

### Get CSR and keys

When you create or modify a third-party DV enrollment, it triggers a PEM-formatted certificate signing request (CSR) with all the information the certificate authority (CA) needs to issue your certificate. CPS encodes the CSR with a private key using either the RSA or the ECDSA algorithm.

1. Send your enrollment ID in the [akamai_cps_csr](../data-sources/cps_csr.md) data source to get the CSR(s) for your enrollment.

    ```hcl

    data "akamai_cps_csr" "example" {
      enrollment_id = 12345
    }
    ```

1. Run `terraform plan` to check your syntax and review the total set of changes you're making and then `terraform apply` to apply your changes to your infrastructure.

    If you're using dual-stacked certificates, you'll see data for both ECDSA and RSA keys.

    <blockquote style="border-left-style: solid; border-left-color: #5bc0de; border-width: 0.25em; padding: 1.33rem; background-color: #e3edf2;"><img src="https://techdocs.akamai.com/terraform-images/img/note.svg" style="float:left; display:inline;" /><div style="overflow:auto;">Dual-stacked certificates are enabled by default for third-party enrollments.</div></blockquote>

1. Send the CSR(s) to a CA for signing.

    <blockquote style="border-left-style: solid; border-left-color: #50af51; border-width: 0.25em; padding: 1.33rem; background-color: #f3f8f3;"><img src="https://techdocs.akamai.com/terraform-images/img/tip.svg" style="float:left; display:inline;" /><div style="overflow:auto;">We recommend the using <code>acme</code> provider along with <code>tls_private_key</code>, <code>acme_registration</code>, <code>acme_certificate</code>. <code>acme</code> also supports maintenance of challenges through <code>dns_challenge</code>.</div></blockquote>

### Upload your signed certificate

Once you have your signed certificate, upload it and your trust chain to CPS.

1. Use the [akamai_cps_upload_certificate](../data-sources/cps_upload_certificate) to upload your certificate for deployment. These certificates must be in PEM format. If they were returned to you in a different format, convert them to PEM format using an SSL converter.

    ```hcl
    resource "akamai_cps_upload_certificate" "upload_cert" {
      enrollment_id                          = 12345
      certificate_ecdsa_pem                  = example_cert_ecdsa.pem
      trust_chain_ecdsa_pem                  = example_trust_chain_ecdsa.pem
      acknowledge_post_verification_warnings = true
      acknowledge_change_management          = true
      wait_for_deployment                    = true
    }
    ```

1. Run `terraform plan` to check your syntax and review the total set of changes you're making and then `terraform apply` to apply your changes to your infrastructure and upload your certificates. We recommend doing this in stage using the [`ChangeManagement` object](https://techdocs.akamai.com/cps/reference/change-management) first, activating it for the production environment after you've tested.

## Manage enrollments

### Modify Subject Alternate Names

In existing enrollments, you can add, modify, and remove existing Subject Alternate Names for your domain. These operations require another domain validation check.

### Modify deployed enrollments

You can edit your network deployment settings for a certificate that is in progress or active on the network.

1. Get the needed certificate by sending the enrollment ID in the [akamai_cps_deployments](../data-sources/cps_deployments).

    ```hcl
    data "akamai_cps_deployments" "example" {
      enrollment_id = 12345
    }
    ```

1. Make edits.
1. Use the [akamai_cps_upload_certificate](../data-sources/cps_upload_certificate) to redeploy your certificate to the network.