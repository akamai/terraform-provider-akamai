---
layout: "akamai"
page_title: "Akamai: property rules template"
subcategory: "Provisioning"
description: |-
  Property Rules template
---

# akamai_rules_template

The `akamai_rules_template` data source allows you to configure a nested block of property rules, criteria, and behaviors. 
The rule tree is composed of a set of templates, which allow to nest other templates as well as interpolate user-defined variables.

The template format used in this data source matches the format used in [Property Manager CLI](https://learn.akamai.com/en-us/learn_akamai/getting_started_with_akamai_developers/developer_tools/getstartedpmcli.html#addanewsnippet)
User-defined variables can be passed either by supplying paths to `variableDefinitions.json` and `variables.json` with syntax used in PM CLI __or__ by supplying variables as set of terraform variables.
## Example Usage

Basic usage:

Using variables passed in data source definition:
```hcl
data "akamai_rules_template" "akarules" {
  template_file = abspath("${path.root}/rules/rules.json")
  variables {
    name = "enabled"
    value = "true"
    type = "bool"
  }
  variables {
    name = "name"
    value = "default"
    type = "string"
  }
}
```

Using variables defined in files:
```hcl
data "akamai_rules_template" "akarules" {
  template_file = abspath("${path.root}/rules/rules.json")
  var_definition_file = abspath("${path.root}/variables/variableDefinitions.json")
  var_values_file = abspath("${path.root}/variables/variables.json")
}
```

## Argument Reference

### Input arguments

#### Required
* `template_file` - the absolute path to the top-level template file in which other templates might be nested.

#### Optional
* `variables` - a definition of a variable. There can be 0 or more `variables` arguments passed. 
This argument conflicts with `variable_definition_file` and `variable_values_file`. `variable` block consists of:
    * `name` - the name of the variable used in template.
    * `type` - the type of the variable - must be one of `string`, `number`, `bool` or `jsonBlock`
    * `value` - the value of the variable passed as string.
* `variable_definition_file` - the absolute path to file containing variable definitions and defaults following the syntax used in PM CLI. 
This argument conflicts with `variables` and is required when `variable_values_file` is set.
* `variable_values_file` - the absolute path to file containing variable values following the syntax used in PM CLI. This argument conflicts with `variables` argument.

#### Computed (returned)

* `json` — The resulting JSON rule tree