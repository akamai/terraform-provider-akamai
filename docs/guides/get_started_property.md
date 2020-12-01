---
title: "[]{#_r00ussmr24 .anchor}Get Started with Properties"
---

**General notes for reviewers:**

**Need to go back and update code examples with the latest from the
repo.**

**Will add cross linking after the content is more settled.**

[[Prerequisites]{.ul}](#prerequisites)

> [[Retrieve the product ID]{.ul}](#retrieve-the-product-id)

[[Make some decisions]{.ul}](#make-some-decisions)

> [[Variables available with the Provider
> module]{.ul}](#variables-available-with-the-provider-module)

[[Add an edge hostname]{.ul}](#add-an-edge-hostname)

> [[Standard TLS examples]{.ul}](#standard-tls-examples)
>
> [[Non-secure hostname]{.ul}](#non-secure-hostname)
>
> [[Secure hostname]{.ul}](#secure-hostname)
>
> [[Domain Suffixes for edge hostname
> types]{.ul}](#domain-suffixes-for-edge-hostname-types)

[[Migrate a property to
Terraform]{.ul}](#migrate-a-property-to-terraform)

[[Create a property]{.ul}](#create-a-property)

[[Set up property rules]{.ul}](#set-up-property-rules)

[[Apply your property changes]{.ul}](#apply-your-property-changes)

[[Activate your property]{.ul}](#activate-your-property)

> [[Create your property activation
> resource]{.ul}](#create-your-property-activation-resource)
>
> [[Test and deploy your property
> activation]{.ul}](#test-and-deploy-your-property-activation)

[[How you can use property
resources]{.ul}](#how-you-can-use-property-resources)

> [[Dynamic Rule Trees Using
> Templates]{.ul}](#dynamic-rule-trees-using-templates)
>
> [[Leverage template_file and snippets to render your configuration
> file]{.ul}](#leverage-template_file-and-snippets-to-render-your-configuration-file)
>
> [[Snippets with Terraform
> (placeholder)]{.ul}](#snippets-with-terraform-placeholder)

You can use the Provisioning resources and data sources to create,
deploy, activate, and manage properties, edge hostnames, and content
provider codes (CP codes).

For more information about Property Manager see:

-   [API
    > documentation](https://developer.akamai.com/api/core_features/property_manager/v1.html)

-   Property Manager
    > d[ocumentation](https://learn.akamai.com/en-us/products/core_features/property_manager.html)
    > page

Prerequisites
=============

To create a property there are a number of dependencies you must first
meet:

-   **Complete the tasks in Get Started.** You need to complete the
    > tasks in the [[Get
    > Started]{.ul}](https://docs.google.com/document/d/1h2U2wu71OEi4dp4eQbjk5750B83zpbyg9mN-rr3Pvg0/edit#heading=h.yv0piszainqd)
    > section of the main Akamai Terraform Provider page before
    > continuing with this section.

-   **Retrieve the Product ID**. The [[Akamai Product
    > ID]{.ul}](https://registry.terraform.io/docs/providers/akamai/g/appendix#common-product-ids)
    > for the product you are using (Ion, DSA, etc.)

-   **Edge hostname:** The Akamai edge hostname for your property. You
    > can [create a new one or reuse an existing one](#_vdthk5hj4s9k).

-   **Origin hostname:** The origin hostname you want your property to
    > point to. Your property should point to an origin hostname you
    > create.

-   **Rules configuration**: The rules.json file contains the base rules
    > for the property. (learn how to leverage the rules.json tree from
    > an existing property
    > [here](https://registry.terraform.io/docs/providers/akamai/g/faq#migrating-a-property-to-terraform))

Retrieve the product ID
-----------------------

When setting up properties, you need to retrieve the ID for the specific
Akamai product you are using. Here's a list of common product IDs:

  **Edge Hostname Type**      **Domain Suffix**
  --------------------------- -----------------------------
  Web Performance Solutions   
  Dynamic Site Accelerator    prd_Site_Accel
  Ion Standard                prd_Fresca
  Ion Premier                 prd_SPM
  Dynamic Site Delivery       prd_Site_Del
  Rich Media Accelerator      prd_Rich_Media_Accel
  IoT Edge Connect            prd_IoT
  Security Solutions          
  Kona Site Defender          prd_Site_Defender
  Media Delivery Solutions    
  Download Delivery           prd_Download_Delivery
  Object Delivery             prd_Object_Delivery
  Adaptive Media Delivery     prd_Adaptive_Media_Delivery

Note that if you have previously used the Property Manager API or CLI
\`set-prefixes\` toggle option, you might have to remove the \"prd\_\"
prefix from your entry.

Make some decisions
===================

Decision: Migrate or create a property

**Make this more of a decision\...like using PM variables**

Variables available with the Provider module
--------------------------------------------

You'll work with three types of variables with the Akamai Provider:

-   **Terraform variables.** Terraform has its own set of \[optional
    > variables\](https://www.terraform.io/docs/commands/environment-variables.html)
    > that let you customize how it works.

-   **Property Manager variables.** For the Provisioning module, you can
    > use the Property Manager's \[variable
    > functionality\]([[https://developer.akamai.com/api/core_features/property_manager/v1.html\#variables]{.ul}](https://developer.akamai.com/api/core_features/property_manager/v1.html#variables)).
    > With Property Manager, you can either use built-in variables or
    > create your own.

Add an edge hostname
====================

You use the akamai_edge_hostname resource to reuse an existing edge
hostname or create a new one. For more information, go to the
[[akamai_edge_hostname]{.ul}
[resource]{.ul}](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/edge_hostname)
page.

Standard TLS examples
---------------------

### Non-secure hostname

The following example will create a Standard TLS edge hostname for
example.com:

resource \"akamai_edge_hostname\" \"example\" { group =
\"\${data.akamai_group.default.id}\" contract =
\"\${data.akamai_contract.default.id}\" product = \"prd_SPM\"
edge_hostname = \"example.com.edgesuite.net\" }

**Note:** Notice that for the group and contract IDs we're using
variables to reference the default group and contract. At runtime,
Terraform will automatically replace these variables with the actual
values.

### Secure hostname

To create a secure hostname, you need to

1.  Add the certificate argument.

2.  Enter the certificate enrollment ID from the [Certificate
    > Provisioning System CLI](https://github.com/akamai/cli-cps). Your
    > entry will looks something like this:

> resource \"akamai_edge_hostname\" \"example\" { group =
> \"\${data.akamai_group.default.id}\" contract =
> \"\${data.akamai_contract.default.id}\" product = \"prd_SPM\"
> edge_hostname = \"example.com.edgesuite.net\" certificate =
> \"\<CERTIFICATE ENROLLMENT ID\>\" }

3.  In the akamai_property or akamai_property_rules resource, set the
    > is_secure flag to true. If you set the value in the either of
    > these resources, it *overrides* the value in the rules.json.

Domain Suffixes for edge hostname types
---------------------------------------

To create different hostname types, you need to change the domain suffix
for the edge_hostname attribute. Here are the common domain suffixes for
edge hostnames:

  **Edge Hostname Type**   **Domain Suffix**
  ------------------------ -------------------
  Enhanced TLS             edgekey.net
  Standard TLS             edgesuite.net
  Shared Cert              akamaized.net

Migrate a property to Terraform
===============================

**FOR REVIEWERS: SHOULD THIS BE HERE, OR IS THIS A BASIC STEP FOR ALL
CUSTOMERS?**

If you have an existing property you would like to migrate to Terraform
we recommend the following process:

1.  Export your rules.json from your existing property (using the API,
    > CLI, or Control Center)

2.  Create a terraform configuration that pulls in the rules.json

3.  Assign a temporary hostname for testing (hint: you can use the edge
    > hostname as the public hostname to allow testing without changing
    > any DNS)

4.  Activate the property and test thoroughly

5.  Once testing has concluded successfully, update the configuration to
    > assign the production public hostnames

6.  Activate again

Once this second activation completes Akamai will automatically route
all traffic to the new property and will deactivate the original
property entirely if all hostnames are no longer pointed at it.

Create a property
=================

You use the [akamai_property
resource](https://registry.terraform.io/docs/providers/akamai/r/property)
to represent your property. Add this new block to your akamai.tf file
after the provider block.

To define the entire configuration, start by opening the resource block
and give it a name. In this case we're going to use the name
\"example\".

Next, we set the name of the property, contact email, product ID, group
ID, CP code, property hostname, and edge hostnames.

Finally, we set up the property rules: First, we specify the [[rule
format
argument]{.ul}](https://registry.terraform.io/docs/providers/akamai/r/property#rule_format),
then add the rules.json data information. In this case, we're using a
variable for the rule.

Once you're done, your property should look like this:

resource \"akamai_property\" \"example\" { name = \"xyz.example.com\" \#
Property Name contact = \[\"user\@example.org\"\] \# User to notify of
de/activations product = \"prd_SPM\" \# Product Identifier (Ion) group =
\"\${data.akamai_group.default.id}\" \# Group ID variable contract =
\"\${data.akamai_contract.default.id}\" \# Contract ID variable
hostnames = { \# Hostname configuration \# \"public hostname\" = \"edge
hostname\" \"example.com\" = \"example.com.edgesuite.net\"
\"www.example.com\" = \"example.com.edgesuite.net\" } rule_format =
\"v2018-02-27\" \# Rule Format rules =
\"\${data.local_file.rules.content}\" \# JSON Rule tree }

Set up property rules
=====================

A property contains the delivery configuration, or rule tree, which
determines how requests are handled. This rule tree is usually
represented using JSON and is often referred to as rules.json.

You can specify the rule tree as a JSON string using the [rules argument
of the akamai_property
resource](https://registry.terraform.io/docs/providers/akamai/r/property#rules).

We recommend storing the rules.json as a JSON file on disk and using
Terraform's local_file data source to ingest it. For example, if our
file is called rules.json, we might create a local_file data source
called rules. We specify the path to rules.json using the filename
argument:

data \"local_file\" \"rules\" { filename = \"rules.json\" }

We can now use \${data.local_file.rules.content} to reference the file
contents in the akamai_property.rules argument.

Apply your property changes
===========================

To actually create our property, we need to instruct Terraform to apply
the changes outlined in the plan. To do this, run this command in the
terminal:

\$ terraform apply

Once the command completes your new property is created. You can verify
this in [Akamai Control Center](https://control.akamai.com/) or via the
[Akamai CLI](https://developer.akamai.com/cli). However, you still have
to activate the property configuration, so let's do that next!

Activate your property
======================

Create your property activation resource
----------------------------------------

To activate your property we need to create a new
[akamai_property_activation
resource](https://registry.terraform.io/docs/providers/akamai/r/property_activation).
This resource manages property activations, letting you specify the
property version to activate and the network to activate it on.

You need to set these arguments:

-   property ID and version arguments, which you can set from the
    > akamai_property resource

-   network to STAGING or PRODUCTION.

-   contact the email addresses to send activation updates to.

-   activate argument to true to kick off the activation process

Here's an example:

resource \"akamai_property_activation\" \"example\" { property =
\"\${akamai_property.example.id}\" version =
\"\${akamai_property.example.version}\" network = \"STAGING\" contact =
\[\"user\@example.org\"\] activate = true }

Test and deploy your property activation
----------------------------------------

Like you did with the property, you should first test the
[akamai_property_activation](https://registry.terraform.io/docs/providers/akamai/r/property_activation)
resource with this command:

\$ terraform plan

This time you will notice how the property is not being modified while
the activation is being added to the plan.

If everything looks good, run this command to start the activation:

\$ terraform apply

This will activate the property on the staging network.

How you can use property resources
==================================

This section includes information on different ways to set up your
property resources.

Dynamic Rule Trees Using Templates
----------------------------------

**THIS IS NOT UP TO DATE. SEE LATEST VERSION IN FAQ.**

If you wish to inject terraform data into your rules.json, for example
an origin address, you can use Terraform templates to do so like so:

First decide where your origin value will come from, this could be
another Terraform resource such as your origin cloud provider, or it
could be a terraform input variable like this:

variable \"origin\" { }

Because we have not specified a default, a value is required when
executing the config. We can then reference this variable using
\${vars.origin} in our template data source:

data \"template_file\" \"init\" { template =
\"\${file(\"rules.json\")}\" vars = { origin = \"\${vars.origin}\" } }

Then in our rules.json we would have:

{

\"name\": \"origin\",

\"options\": {

\"hostname\": \"\*\*\${origin}\*\*\",

\...

}

},

You can also inject entire JSON blocks using the same mechanism:

{

\"rules\": {

\"behaviors\": \[

\${origin}

\]

}

}

Leverage template_file and snippets to render your configuration file
---------------------------------------------------------------------

The 'rules' argument within the akamai_property resource enables
leveraging the full breadth of the Akamai's property management
capabilities. This requires that a valid json string is passed on as
opposed to a filename. Terraform enables this via the \"local_file\"
data source that loads the file.

data \"local_file\" \"rules\" { filename =
\"\${path.module}/rules.json\" } resource \"akamai_property\"
\"example\" { \.... rules =
\"\${data.local_file.terraform-demo.content}\" }

Microservices driven and DevOps users typically want additional
flexibility - using delegating snippets of the configuration to
different users, and inserting variables within code. Terraform\'s
\"template_file\" is provides that additoinal value. Use the example
below to construct a template_file data resource that helps maintain a
rules.json with variables inside it:

data \"template_file\" \"rules\" { template =
\"\${file(\"\${path.module}/rules.json\")}\" vars = { origin =
\"\${var.origin}\" } } resource \"akamai_property\" \"example\" { \...
rules = \"\${data.template_file.rules.rendered}\" } \"rules\": {
\"name\": \"default\", \"children\": \[
\${file(\"\${snippets}/performance.json\")} \],
\${file(\"\${snippets}/default.json\")} \], \"options\": {
\"is_secure\": true } }, \"ruleFormat\": \"v2018-02-27\" }

More advanced users want different properties to use different rule
sets. This can be done by maintaining a base rule set and then importing
individual rule sets. To do this we first create a directory structure -
something like:

rules/rules.json

rules/snippets/routing.json

rules/snippets/performance.json

...

The \"rules\" directory contains a single file \"rules.json\" and a sub
directory containing all rule snippets. Here, we would provide a basic
template for our json.

\"rules\": {

\"name\": \"default\",

\"children\": \[

\${file(\"\${snippets}/performance.json\")}

\],

\${file(\"\${snippets}/routing.json\")}

\],

\"options\": {

\"is_secure\": true

}

},

\"ruleFormat\": \"v2018-02-27\"

Then remove the \"template_file\" section we added earlier and replace
it with:

data \"template_file\" \"rule_template\" { template =
\"\${file(\"\${path.module}/rules/rules.json\")}\" vars = { snippets =
\"\${path.module}/rules/snippets\" } } data \"template_file\" \"rules\"
{ template = \"\${data.template_file.rule_template.rendered}\" vars = {
tdenabled = var.tdenabled } }

This enables Terraform to process the rules.json & pull each fragment
that\'s referenced and then to pass its output through another
template_file section to process it a second time. This is because the
first pass creates the entire json and the second pass replaces the
variables that we need for each fragment. As before, we can utilize the
rendered output in our property definition.

resource \"akamai_property\" \"example\" { \.... rules =
\"\${data.template_file.rules.rendered}\" }

Snippets with Terraform (placeholder)
-------------------------------------

**FOR REVIEWERS: THIS IS A PLACEHOLDER. IF RELEVANT, WOULD LIKE TO TALK
ABOUT USING SNIPPETS GENERALLY ACROSS ALL AKAMAI PROVIDER RESOURCES AND
DATA SOURCES.**

We provide Property Manager rules as raw JSON when setting up Akamai
properties with the Akamai Provider.

You can do this in the akamai_property resource with the \"rules\"
configuration parameter. This requires valid JSON, as opposed to a
filename. Terraform gives us a \"template_file\" data provider which we
then use to provide the Property Manager rules from a file while
interpolating Terraform variables.

data \"template_file\" \"rules\" {

template = \"\${file(\"\${path.module}/rules.json\")}\"

vars = {

origin = \"\${var.origin}\"

}

}

resource \"akamai_property\" \"example\" {

\...

rules = \"\${data.template_file.rules.rendered}\"

}

Our json is now templated but it\'s still a large monolithic blob. We
could still use it if all our properties were exactly the same. However,
it is common for different properties to want different rule sets. It
would be ideal if we could provide a base rule template and then import
rule sets individually to give us more flexibility.

First, let's create our directory structure:

rules/rules.json

rules/snippets/default.json

rules/snippets/performance.json

The \"rules\" directory contains a single file \"rules.json\" and a sub
directory containing all rule snippets. Let's provide a basic template
for our JSON.

\"rules\": {

\"name\": \"default\",

\"children\": \[

\${file(\"\${snippets}/performance.json\")}

\],

\${file(\"\${snippets}/default.json\")}

\],

\"options\": {

\"is_secure\": true

}

},

\"ruleFormat\": \"v2018-02-27\"

}

We\'re pulling in two snippets: one for the default rules and another
for the performance rule. We could add more. Each snippet would be just
JSON fragments for each section of the rule tree, perhaps from a central
repository. To make this work, we need to define our \"template_file\"
section like this:

data \"template_file\" \"rule_template\" {

template = \"\${file(\"\${path.module}/rules/rules.json\")}\"

vars = {

snippets = \"\${path.module}/rules/snippets\"

}

}

data \"template_file\" \"rules\" {

template = \"\${data.template_file.rule_template.rendered}\"

vars = {

tdenabled = var.tdenabled

\...

}

}

The template now lets Terraform process the rules.json and pull each
fragment that\'s referenced. It can then pass the output through another
template_file section to process it a second time: the first pass
creates the entire JSON and the second pass replaces the variables that
we need for each fragment. We can utilize the rendered output in our
property definition:

resource \"akamai_property\" \"example\" {

\...

rules = \"\${data.template_file.rules.rendered}\"

}
