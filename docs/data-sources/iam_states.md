---
layout: "akamai"
page_title: "Akamai: akamai_iam_states"
subcategory: "Identity and Access Management"
description: |-
 IAM States
---

# akamai_iam_states

Use `akamai_iam_states` to list US states or Canadian provinces. If `country=USA` you may enter a value of `TBD` if you don't know a user's state. Administrators use this data source to set a user's state.

## Example usage

Basic usage:

```hcl
data "akamai_iam_states" "states" {
  country = "canada"
}

output "supported_states" {
  value = data.akamai_iam_states.states
}
```

## Argument reference

These attributes are supported:

* `country` — (required, string) Sepcifies USA or Canada.

## Attributes reference

These attributes are returned:

* `states` — A list of states.

[API Reference](https://techdocs.akamai.com/iam-api/reference/get-common-states)
