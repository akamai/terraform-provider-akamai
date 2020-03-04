# Snippets with Terraform

If we want to use our Terraform provider for anything useful, we need to provide Property Manager rules as raw json. This is configured in the akamai_property resource using the "rules" configuration parameter. This requires that a string of valid json is given to it (as opposed to a filename). To achieve this in a useful way, Terraform gives us a "template_file" data provider which we can use to provide the Property Manager rules from a file whilst interpolating Terraform variables.

```
data "template_file" "rules" {
 template = "${file("${path.module}/rules.json")}"
 vars = {
 origin = "${var.origin}"
 }
}

resource "akamai_property" "example" {
...
 rules = "${data.template_file.rules.rendered}"
}
``` 

This is getting closer to how we'd like to do things in a DevOps world but we can do better. Our json is now templated but it's still a large monolithic blob. This would be ok if all our properties were exactly the same but often different properties want different rule sets. It would be ideal if we could provide a base rule template & then import rule sets individually to give us more flexibility. 

Firstly, we need to create our directory structure....

```
rules/rules.json
rules/snippets/default.json
rules/snippets/performance.json
```

The "rules" directory contains a single file "rules.json" and a sub directory containing all rule snippets. Here, we would provide a basic template for our json.

```
"rules": {
  "name": "default",
  "children": [
    ${file("${snippets}/performance.json")}
    ],
    ${file("${snippets}/default.json")}
  ],
"options": {
    "is_secure": true
  }
},
  "ruleFormat": "v2018-02-27"
}
```

As you can see, we're pulling in just two snippets - one for the default rules and another for the performance rule. We could add more. Each snippet would be just json fragments for each section of the rule tree (perhaps from a central repo). To make this work, we need to define our "template_file" section like so ...

```
data "template_file" "rule_template" {
 template = "${file("${path.module}/rules/rules.json")}"
 vars = {
 snippets = "${path.module}/rules/snippets"
 }
}

data "template_file" "rules" {
 template = "${data.template_file.rule_template.rendered}"
 vars = {
 tdenabled = var.tdenabled
 ...
 }
}
```

This will tell Terraform to process the rules.json & pull each fragment that's referenced and then to pass its output through another template_file section to process it a second time. This is because the first pass creates the entire json and the second pass replaces the variables that we need for each fragment. We can utilize the rendered output in our property definition exactly the same as we did before.

```
resource "akamai_property" "example" {
...
rules = "${data.template_file.rules.rendered}"
}
```
