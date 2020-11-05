---
layout: "akamai"
page_title: "Akamai: akamai_rules_template"
subcategory: "Provisioning"
description: |-
 Property Rules Template
---

# akamai_property_rules_template

The `akamai_property_rules_template` data source allows you to configure a nested block of property rules, criteria, and behaviors. 
The rule tree is composed of a set of templates, which allow to nest other templates as well as interpolate user-defined variables.

The template format used in this data source matches the format used in [Property Manager CLI](https://learn.akamai.com/en-us/learn_akamai/getting_started_with_akamai_developers/developer_tools/getstartedpmcli.html#addanewsnippet)
User-defined variables can be passed either by supplying paths to `variableDefinitions.json` and `variables.json` with syntax used in PM CLI __or__ by supplying variables as set of terraform variables.

## Referencing sub-files from a template
Each template can include other template files by including them in the currently loaded file.  For example to include
`example-file.json` from template directory by using the following syntax `"#include:example-file.json"` including 
quotes.  All files are resolved in relation to the template directory.  

## Inserting variables in a template
Variables can also be included in a template by using a string like `“${env.<variableName>}"` including quotes.  These
are variables passed into the template call and are in contract to terraform variables which should resolve normally.

## Example Usage

### Variables passed in data source definition:
```hcl
data "akamai_property_rules_template" "akarules" {
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

### Variables defined in files:
```hcl
data "akamai_property_rules_template" "akarules" {
  template_file = abspath("${path.root}/rules/rules.json")
  var_definition_file = abspath("${path.root}/variables/variableDefinitions.json")
  var_values_file = abspath("${path.root}/variables/variables.json")
}
```

### Example of template in use:
Given a template file like the following (assuming all sub templates mentioned exist) usage is as follows:

templates/main.json:
```json
{
  "rules": {
    "name": "default",
    "children": [
      "#include:Default_GKE_Origin.json",
      "#include:Origin_Selection.json",
      "#include:Performance.json",
      "#include:Offload.json",
      "#include:Restrict_Method.json",
      "#include:Set_PMUSER_VARS.json",
      "#include:Modify_cache_key.json"
    ],
    "behaviors": "#include:behaviors_default.json",
    "options": {
      "is_secure": “${env.secure}"
    },
    "variables": "#include:akamai_variables.json",
    "comments": "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings."
  }
}
```
You then can define a terraform file like the following to expand the template file and use it with a property :
```hcl-terraform
data "akamai_rules_template" "example" {
  template_file = abspath("${path.root}/templates/main.json")
  variables {
      name = "secure"
      value = "false"
      type = "bool"
  }
  variables {
      name = "caching_ttl"
      value = "3d"
      type = "string"
  }
}

resource "akamai_property_version" "example" {
    contract_id = var.contractid
    group_id    = var.groupid
    property_id = var.propertyid
    hostnames = {
      "example.org" = "example.org.edgesuite.net"
      "www.example.org" = "example.org.edgesuite.net" 
      "sub.example.org" = "sub.example.org.edgesuite.net"
    }
    rule_format = "v2020-03-04"
    rules       = data.akamai_rules_template.example.json
}
```

## Argument Reference

### Input arguments

#### Required arguments
* `template_file` - (Required) the absolute path to the top-level template file in which other templates might be nested.

#### Optional arguments
* `variables` - (Optional) a definition of a variable. There can be 0 or more `variables` arguments passed. 
This argument conflicts with `variable_definition_file` and `variable_values_file`. `variable` block consists of:
    * `name` - the name of the variable used in template.
    * `type` - the type of the variable - must be one of `string`, `number`, `bool` or `jsonBlock`
    * `value` - the value of the variable passed as string.
* `variable_definition_file` - (Optional) the absolute path to file containing variable definitions and defaults following the syntax used in PM CLI. 
This argument conflicts with `variables` and is required when `variable_values_file` is set.
* `variable_values_file` - (Optional) the absolute path to file containing variable values following the syntax used in PM CLI. This argument conflicts with `variables` argument.

## Attributes Reference

The following are the return attributes:

* `json` — The fully expanded template with variables and all sub-templates resolved.

