---
layout: "akamai"
page_title: "Akamai: Appendix"
sidebar_current: "docs-akamai-guide-appendix"
description: |-
  Appendix
---

# Appendix

## Domain Suffixes for Different Edge Hostname Types

Each type of edge hostname has its own domain suffix. Knowing which one to use is important when setting the cnameTovalue:

| Edge Hostname Type | Domain Suffix |
|--------------------|---------------|
| Enhanced TLS       | `edgekey.net` |
| Standard TLS       | `edgesuite.net` |
| Shared Cert        | `akamaized.net` |
| Non-TLS            | `edgesuite.net` |

## Secure Hostnames

For secure hostnames you must include the certificate enrollment ID in your [`akamai_edge_hostname` resource](/r/edge_hostname.html).

1. Retrieve the enrollment-id from the [CPS CLI](https://github.com/akamai/cli-cps) 
2. Enter the ID as the certificate attribute. 

## Common Product IDs

Leveraging Product IDs in your setup requires you to retrieve the ID for the specific Akamai product you are using. The following is a list of commonly used product IDs for different products:

| Product | Code |
|---|---|
| Ion/Ion Premier          | `prd_SPM`        |
| Ion Standard             | `prd_Fresca`     |
| Dynamic Site Accelerator | `prd_Site_Accel` |
| Rich Media Accelerator   | `prd_Rich_Media_Accel` |
| Ion Standard             | `prd_Fresca`           |