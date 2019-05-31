---
layout: "akamai"
page_title: "Akamai: cps enrollment"
sidebar_current: "docs-akamai-resource-cps-enrollment"
description: |-
  CPS Enrollment
---

# akamai_cps_enrollment


The `akamai_cps_enrollments` provides the resource for managing CPS enrollments. A CPS enrollment is the collection of settings for:

<ol>
<li>One or more active and pending X509 certificate(s).</li>
<li>A reference to the key pairs for all X509 certificates.</li>
<li>The settings for how SSL connections utilizing this certificate collection managed by Akamai.</li>
<li>Information regarding contact information used in making validation requests.</li>

A CPS enrollment is the most fundamental and definitive concept. It behaves as a core container for all the operations that clients can perform within CPS. CPS is a certificate life cycle management tool and a CPS enrollment is the agent in this tool that allows users display all the information about the process that certificate goes through from the time it was requested, through renewal or removal.


## Example Usage

Basic usage:

```hcl
resource "akamai_cps_enrollment" "cps_terraformdemo" {

    contract_id = "ctr_XXX"
    admin_contact = ["devrel@akamai.com"]
    validation_type = ‘third-party”
    techcontact = ["devrel@akamai.com"]
    ra = “third-party”
    enable_multi_stacked_certificates = “true”
    change_management = “true”
    csr {
        "cn": "www.example.com",
        "c": "US",
        "st": "MA",
        "l": "Cambridge",
        "o": "Akamai",
        "ou": "WebEx",
        "sans": [
        "san1.example.com",
        "san2.example.com",
        "san3.example.com",
        "san4.example.com"
        ] 
        },
    org {
        "name": "Akamai Technologies",
        "addressLineOne": "150 Broadway",
        "addressLineTwo": null,
        "city": "Cambridge",
        "region": "MA",
        "postalCode": "02142",
        "country": "US",
        "phone": "617-555-0111"
        }

}

```

## Argument Reference

The following arguments are supported:

`*contract_id` — (Required) The contract ID.
`*deploy_not_after` — (Optional) Do not deploy after this date (UTC).
`*deploy_not_before` — (Optional) Do not deploy before this date (UTC).
`*admincontact` — (Optional) The contact of the admin.
`*certificate_chain_type` — (Optional) The certificate trust chain type.
`*certificate_type` — (Required) The type of the certificate.
`*change_management` — (Optional) If you turn change management on for an enrollment, it stops CPS from deploying the certificate to the network until you acknowledge that you are ready to deploy the certificate.
`*csr` — (Required) Certificate Signing request (CSR) is a block of encoded text that is given to a Certificate Authority when applying for an SSL Certificate.
`*enable_multistacked_certificates` —(Optional,boolean) Enable Dual-Stacked certificate deployment for this enrollment.  Default:false
`*max_allowed_san_names` — (Optional) Maximum number of SAN names supported for this enrollment type.
`*max_allowed_wildcard_san_names` — (Optional) Maximum number of wildcard SAN names supported for this enrollment type. 
`*network_configurations` — (Optional) Settings that specify any network information and TLS Metadata you want CPS to use to push the completed certificate to the network.
`*org` — (Required) Your organization information.
`*pending_changes` — (Optional) Returns the Changes currently pending in CPS. The last item in the array is the most recent change.
`*ra` — (Required) The registration authority or certificate authority (CA) you want to use to obtain a certificate. A CA is a trusted entity that signs certificates and can vouch for the identity of a website. Either symantec, lets-encrypt, or third-party.
`*signature_algorithm` — (Optional) The SHA (Secure Hash Algorithm) function. NSA designed this function to produce a hash of certificate contents, which is used in a digital signature. Specify either SHA-1 or SHA-256. We recommend you use SHA–256.
`*tech_contact` — (Required) Your technical contact in akamai.
`*third_party` — (Optional) Specifies that you want to use a third party certificate. This is any certificate that is not issued through CPS.
`*validation_type` — (Optional) There are three types of validation. Domain Validation (DV), which is the lowest level of validation. The CA validates that you have control of the domain. CPS supports DV certificates issued by Let’s Encrypt, a free, automated, and open CA, run for public benefit. Organization Validation (OV), which is the next level of validation. The CA validates that you have control of the domain. Extended Validation (EV), which is the highest level of validation in which you must have signed letters and notaries sent to the CA before signing. You can also specify third party as a type of validation, if you want to use a signed certificate obtained by you from a CA not supported by CPS. Either dv, ev, ov, or third-party.
