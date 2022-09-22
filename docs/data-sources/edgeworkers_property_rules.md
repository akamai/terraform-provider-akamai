---
layout: akamai
subcategory: EdgeWorkers
---

# akamai_edgeworkers_property_rules

The `akamai_edgeworkers_property_rules` data source lets you configure a rule tree through the use of JSON template files. A rule tree is a nested block of property rules in JSON format that include match criteria and behaviors.

With this data source you define the location of the JSON template files and provide information about any user-defined variables included within the templates.

The template format used in this data source matches those used in the [Property Manager CLI](https://github.com/akamai/cli-property-manager#add-a-new-snippet).

You can pass user-defined variables by supplying either:

* paths to `variableDefinitions.json` and `variables.json` with syntax used in Property Manager CLI, or
* a set of Terraform variables.

## Referencing sub-files from a template

You can split each template out into a series of smaller template files.
To add them to this data source, you need to include them in the currently loaded file, which corresponds to the value in the `template_file` argument.
For example, to include `example-file.json` from the `property-snippets` directory, use this syntax including the quotes: `"#include:example-file.json"`.  
Make sure the `property-snippets` folder contains only `.json` files.
All files are resolved relative to the directory that contains the starting template file.

### Example usage: template files

Here's an example of what a JSON-based template file with its nested templates might look like:

property-snippets/main.json:
```json
{
    "rules":{
        "name":"default",
        "name":"edgeWorker",
        "options":{
            "enabled":true,
            "edgeWorkerId":"42"
        }
    }
}
```

You can then define a Terraform configuration file like this, which pulls in the `main.json` file above and uses it with a property:

```hcl
data "akamai_property_rules_template" "example" {
  template_file = abspath("${path.root}/property-snippets/main.json")
  variables {
    name  = "secure"
    value = "false"
    type  = "bool"
  }
  variables {
    name  = "caching_ttl"
    value = "3d"
    type  = "string"
  }
}

resource "akamai_property" "example" {
  name        = "dev.example.com"
  contract_id = var.contractid
  group_id    = var.groupid
  hostnames = {
    "example.org"     = "example.org.edgekey.net"
    "www.example.org" = "example.org.edgekey.net"
    "sub.example.org" = "sub.example.org.edgekey.net"
  }
  rule_format = "v2020-03-04"
  rules       = data.akamai_property_rules_template.example.json
}
```

## Inserting variables in a template

You can also add variables to a template by using a string like `"${env.<variableName>}"`. You'll need the quotes here too. These variables follow the format used in the [Property Manager CLI](https://github.com/akamai/cli-property-manager#update-the-variabledefinitions-file).

### Example usage: variables

This first example shows two variables passed in data source definition:

```hcl
data "akamai_property_rules_template" "akarules" {
  template_file = abspath("${path.root}/rules/main.json")
  variables {
    name  = "enabled"
    value = "true"
    type  = "bool"
  }
  variables {
    name  = "name"
    value = "default"
    type  = "string"
  }
}
```

In this second example, the variables defined refer to files shared with a [Property Manager CLI pipeline](https://github.com/akamai/cli-property-manager#akamai-pipeline-workflow):

```hcl
data "akamai_edgeworkers_property_rules" "akarules" {
  template_file       = abspath("${path.root}/property-snippets/main.json")
  var_definition_file = abspath("${path.root}/environments/variableDefinitions.json")
  var_values_file     = abspath("${path.root}/environments/dev.example.com/variables.json")
}
```

## Argument reference

* `edgeworker_id` - (Required) Unique identifier of an EdgeWorker ID.

## Attributes reference

This data source returns this attribute:

* `json` - The property rule that specifies which EdgeWorker ID to enable.
