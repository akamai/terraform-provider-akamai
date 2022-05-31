---
layout: "akamai"
page_title: "Akamai: akamai_property_rules_template"
subcategory: "Property Provisioning"
description: |-
 Property rules template
---
 
# akamai_property_rules_template

The `akamai_property_rules_template` data source lets you define a rule tree. A rule tree is a nested block of property rules in JSON format that includes match criteria and behaviors. You can break a rule tree out into smaller JSON  template files that cover individual rules.

With this data source, you define which JSON template files to use for your property. You can also set values for variables. 

This data source uses the rule template format from the [Property Manager CLI](https://github.com/akamai/cli-property-manager#set-up-property-snippets).

~> You can define variables either by using the Property Manager CLI syntax or by using standard Terraform variables.

## How to work with JSON template files 

You have a few options when working with rule template files: 

* Use a single JSON file that includes all rules for the property.
* Create separate JSON template files for each rule and store them in the `property-snippets` directory.
* Reference individual template files directly in this data source.

### Use a single JSON rule tree file

If you’re using a single JSON rule tree file for your property, set the `template_file` argument with the absolute or relative path to the file. For example: `template_file = abspath("${path.root}/property-snippets/main.json")` 

You can also set up a single JSON template file that calls individual template files. To reference other template files, use `include` statements in the `children` array . For example:  

```json
    "children": [
      "#include:Performance.json",
      "#include:Offload.json"
    ]
```

You have to store all files in the directory listed in the `template_file` argument. 

### Create a set of JSON template files

If you have a set of JSON template files you want to call: 

1. Put them all in a directory called `property-snippets`. 
1. Make sure the `property-snippets` folder only contains `.json` files. 
1. Add the `template_dir` argument. For example: `template_dir = "property-snippets/"`.

~> This directory name is different from the one required for the Property Manager CLI, which is called `config-snippets`. 

### Reference template files individually

You can also pass in specific rule files with this data source. For this option, add the `template_data` argument and  use Terraform’s `jsonencode` function to add the supporting JSON syntax.  For example: 

```hcl
template_data = jsonencode({
      "rules": {
        "name": "default",
        "children": [
          "#include:rules.json"
        ]
      }
    })
```
 
## How to use property variables with a template

You can also add variables to a template by using a string like `“${env.<variableName>}"`. These property variables follow the file structure and syntax used when [creating a pipeline in the Property Manager CLI](https://github.com/akamai/cli-property-manager#create-and-set-up-a-new-pipeline). 

You’ll need to create a `variableDefinitions.json` file to define your variables and their default values. 

If working with multiple environments, you can also set up variables.json files to override these default values. Since the file name should always be `variables.json`, you’ll need to create a file for each environment and a separate folder to house it in.

~> Property variables are separate from Terraform variables. Terraform variables work as expected in this data source.

## Example usage: JSON template files

Here are some examples of how you can set up your JSON template files for use with this data source.

### Single JSON template that calls other templates

Here's an example of a JSON template file with nested templates:

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
      "is_secure": "${env.secure}"
    },
    "variables": "#include:akamai_variables.json",
    "comments": "The behaviors in the Default Rule apply to all requests for the property hostnames unless another rule overrides the Default Rule settings."
  }
}
```

### Individual JSON rule template file

Here’s a simple default rule example that you can include inside the `template_data` argument:

```json
{
  "name": "default",
  "children": [],
  "behaviors": "#include:behaviors_default.json",
  "criteria": [],
  "criteriaMustSatisfy": "all"
}
```

## Example usage: Add templates to the data source

Here are some examples of how you can call your JSON template files with this data source.

### Call a single JSON template file

This example shows how to define a specific JSON template using the `template_file` argument:

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
    hostnames   = {
      "example.org"     = "example.org.edgesuite.net"
      "www.example.org" = "example.org.edgesuite.net"
      "sub.example.org" = "sub.example.org.edgesuite.net"
    }
    rule_format = "v2020-03-04"
    rules       = data.akamai_property_rules_template.example.json
}
```

### Call individual template files with this data source

This second example shows how to call a specific JSON template using the `template_data` field:
 
```hcl
data "akamai_property_rules_template" "example" {
  template {
    template_data = jsonencode({
      "rules": {
        "name": "default",
        "children": [
          "#include:rules.json"
        ]
      }
    })
    template_dir = "property-snippets/"
  }
}
 
resource "akamai_property" "example" {
  name = "dev.example.com"
  contract_id = var.contractid
  group_id    = var.groupid
  rule_format = "v2020-03-04"
  rules       = data.akamai_property_rules_template.example.json
}
```
 
## Example usage: Variables

You can add variables individually or reference variable definition files.
 
### Define variables in the data source

This example shows two variables defined in the data source:
 
```hcl
data "akamai_property_rules_template" "akarules" {
  template_file = abspath("${path.root}/property-snippets/main.json")
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
 
### Call variable definition files

In this example, you store the variables in separate JSON files that follow the naming and formats used with the [Property Manager CLI pipeline](https://github.com/akamai/cli-property-manager#akamai-pipeline-workflow):
 

```hcl
data "akamai_property_rules_template" "akarules" {
  template_file       = abspath("${path.root}/property-snippets/main.json")
  var_definition_file = abspath("${path.root}/environments/variableDefinitions.json")
  var_values_file     = abspath("${path.root}/environments/dev.example.com/variables.json")
}
```

## Argument reference

* `template_file` - (Optional) The absolute path to your top-level JSON template file. The top-level template combines smaller, nested JSON templates to form your property rule tree. This argument conflicts with the `template` argument.  
* `template` - (Optional) The template you use in your configuration. This argument conflicts with the `template_file` argument.
  * `template_data` - (Required) The content of the JSON template as a string. 
  * `template_dir` - (Required) The absolute or relative path to the directory containing the template files. The path must end with `property-snippets`, the required directory name. For example: `template_dir = abspath("${path.root}/property-snippets/")`, or `template_dir = "property-snippets/"`
* `variables` - (Optional) The definition of one or more variables. This argument conflicts with the `var_definition_file` and `var_values_file` arguments. A `variables` block includes:
    * `name` - The name of the variable used in the template.
    * `type` - The type of variable: `string`, `number`, `bool`, or `jsonBlock`.
    * `value` - The value of the variable passed as a string.
* `var_definition_file` - (Optional) Required when using `var_values_file`. The absolute path to the file containing variable definitions and defaults. This argument conflicts with the `variables` argument. 
* `var_values_file` - (Optional) Required when using `var_definition_file`. The absolute path to the file containing variable values. This argument conflicts with the `variables` argument.  
 
## Attributes reference
 
This data source returns this attribute:
 
* `json` - The fully expanded template with variables and all nested templates resolved.