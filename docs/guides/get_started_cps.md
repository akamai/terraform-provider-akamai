---
layout: "akamai"
page_title: "Akamai: Get Started with Certificate Provisioning"
description: |-
  Get Started with Akamai Certificate Provisioning using Terraform
---

# Get Started with Certificate Provisioning

When setting a secure website, you need to ensure that the delivery of content to and from that site is secure. Certificate Provisioning System (CPS) provides the SSL/TLS certificates that authenticate the secure connection the browsers make during a secure delivery.

## Prerequisites

1. Complete the tasks in the
[Get Started with the Akamai Terraform Provider](../guides/get_started_provider.md)
guide.
1. Get familiar with:

    * SSL/TLS certificates
    * Certificate authorities (CAs)
    * How Akamai obtains certificates on the requester’s behalf, which includes the generation of public/private key pairs and certificate signing requests (CSRs)
    * DNS

    If you have questions about these concepts, contact your Akamai account representative.

## CPS Workflow

With the Certificate Provisioning module, you can [create](#validate-your-domain) new Domain Validation (DV) certificate enrollments or [make changes](#modify-subject-alternate-names-sans) to the existing ones.

### Validate your domain

Create a new enrollment and prove that you control the domain you want to secure.

* [Get the contract ID](#get-the-contract-ID). You need your contract ID to create an enrollment.
* [Create a DV enrollment](#create-an-enrollment). An enrollment is a core container for all the operations you perform within CPS. Currently, you can create DV certificates with this module.
* [Validate control over your domains](#validate-control-over-your-domains). Complete the HTTP or DNS challenge to prove that you control the domain names you requested the certificate for.
* [Send the validation acknlowledgement](#send-the-validation-acknowledgement).

## Get the contract ID

When setting up enrollments, you need to retrieve the Akamai [`contract ID`](../data-sources/contract.md).

-> **Note** If you're currently using prefixes with your IDs, you might have to remove the `ctr_` prefix from your entry. For more information about prefixes, see the [ID prefixes](https://developer.akamai.com/api/core_features/property_manager/v1.html#prefixes) section of the Property Manager API (PAPI) documentation.

## Create a DV enrollment

You use the [akamai_cps_dv_enrollment](../resources/cps_dv_enrollment.md) resource to create a new enrollment for the Domain Validation (DV) certificate type.

Once you set up the `akamai_cps_dv_enrollment` resource, run `terraform apply`. Terraform shows an overview of changes, so you can still go back and modify the configuration, or confirm to proceed. See [Command: apply](https://www.terraform.io/docs/commands/apply.html)

## Validate control over your domains

You need to prove that you have control over each of the domains listed in the certificate. When you create or modify a DV enrollment that generates a certificate signing request (CSR), CPS automatically sends it to [Let’s Encrypt](https://letsencrypt.org/) for signing. Let’s Encrypt sends back a challenge for each domain listed on your certificate. You prove that you have control over the domains listed in the CSR by redirecting your traffic to Akamai. This allows Akamai to complete the challenge process for you by detecting the redirect and answering Let’s Encrypt’s challenge.

Complete one of the challenges returned by [akamai_cps_dv_enrollment](../resources/cps_dv_enrollment.md) resource.

    * For `http_challenges`, create a file with a token and put it in the designated folder on your site. Once Akamai detects the file is in place, it asks Let's Encrypt to validate the domain.
    * For `dns_challenges`, add a `TXT` record to the DNS configuration of your domain. If you're using [Akamai DNS Provider](../guides/get_started_dns_zone.md), you can create DNS records for the provided SANs from the same `config` file.

~> **Note** If the challenge token expires, run `terraform-apply` again to pull the latest token. Terraform doesn't notify you of any changes to tokens.

## Send the validation acknowledgement

Once you complete the Let’s Encrypt challenges, use the `akamai_cps_dv_validation` resource to send the acknowledgement to CPS and inform it that tokens are ready for validation.

## Modify Subject Alternate Names (SANs)

In the already created enrollments, you can modify existing SANs, or add and remove SANs for your domain. These operations require another domain validation check.
