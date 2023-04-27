Use our Terraform provider to provision and manage your Akamai configurations in Terraform.  

## Before you begin

- Understand the [basics of Terraform](https://learn.hashicorp.com/terraform?utm_source=terraform_io).
- Install our [Terraform CLI](https://github.com/akamai/cli-terraform) to export your Akamai configurations.

## How it works

![Terraform overview](https://techdocs.akamai.com/terraform-images/img/ext-tf-gs.png)

## Start your configuration

Your Akamai Terraform configuration starts with listing us as a required provider, stating which version of our provider to use, and initializing both a Terraform session and an install of our provider.  

1. Create a file named `akamai.tf` in your project directory. This is your base Akamai configuration file.
2. Add our provider to your file.

   ```
   terraform {
     required_providers {
       akamai = {
         source = "akamai/akamai"
         version = "3.6.0"
       }
     }
   }

   provider "akamai" {
     # Configuration options
   }
   ```

~> If you choose to use multiple configuration files to limit the scope of work or to share work across teams, our provider information only needs to be in one.

## Get authenticated

Authentication credentials for the majority of our API use a hash-based message authentication code or HMAC-SHA-256 created through an API client. We recommend each member of your team use their own client set up locally to prevent accidental exposure of credentials.

There are different types of API clients that grant access based on your need, role, or how many accounts you manage.

| API client type                                                                                | Description                                                                                                                                       |
| ---------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Basic](#create-a-basic-api-client)                                                            | Access all API associated with your account without any specific configuration. Individual service read/write permissions are based on your role. |
| [Advanced](https://techdocs.akamai.com/developer/docs/create-a-client-with-custom-permissions) | Configurable permissions to limit or narrow down scope of the API associated with your account.                                                   |
| [Managed](https://techdocs.akamai.com/developer/docs/manage-many-accounts-with-one-api-client) | Configurable permissions that work for multiple accounts.                                                                                         |

### Create a basic API client

1. Navigate to the [Identity and Access Management](https://control.akamai.com/apps/identity-management/#/tabs/users/list) section of Akamai Control Center and click **Create API Client**.

2. Click **Quick** and then **Download** in the Credentials section.

3. Open the downloaded file with a text editor and add `[default]` as a header above all text.

   ```bash
   [default]
   client_secret = C113nt53KR3TN6N90yVuAgICxIRwsObLi0E67/N8eRN=
   host = akab-h05tnam3wl42son7nktnlnnx-kbob3i3v.luna.akamaiapis.net
   access_token = akab-acc35t0k3nodujqunph3w7hzp7-gtm6ij
   client_token = akab-c113ntt0k3n4qtari252bfxxbsl-yvsdj
   ```

4. Save your credentials as an EdgeGrid resource file named `.edgerc` in your local home directory.

   ```bash
   // Linux
   $ /home/{username}/.edgerc

   // MacOS
   $ /Users/{username}/.edgerc

   // Windows
   C:\Users\{username}\.edgerc
   ```

5. In the provider block of your `akamai.tf` file, add in a pointer to your `.edgerc` file.

   ```hcl
   provider "akamai" {
     edgerc = "~/.edgerc"
     config_section = "default"
   }
   ```

## Initialize our provider

To install our provider and begin a Terraform session, run `terraform init`. The response log verifies your initialization along with a notice that the rest of the `terraform` commands should work.

## Add subprovider resources

Each of our subproviders use a set of resource objects that build out infrastructure components and data sources that provide information to and about those resources. Add these to your configurations manually or import them.

- Copy/paste or pull in our [examples](https://github.com/akamai/terraform-provider-akamai/tree/master/examples).
- Import a set of components using our [CLI for Terraform Provider](https://github.com/akamai/cli-terraform).
- Export and use your company's existing configurations.

Use the table to find information about the subprovider you’re using.

|Subprovider|Description|
|---|---|
|[Application Security](https://techdocs.akamai.com/terraform/v3.6/docs/configure-appsec)|Manage security configurations, security policies, match targets, rate policies, and firewall rules.|
|[Bot Manager](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-botman)|Identify, track, and respond to bot activity on your domain or in your app.|
|[Certificates](https://techdocs.akamai.com/terraform/v3.6/docs/cps-integration-guide)|Full life cycle management of SSL certificates for your ​Akamai​ CDN applications.|
|[Cloudlets](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-cloudlets)|Solve specific business challenges using value-added apps that complement ​Akamai​'s core solutions.|
|[DataStream](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-datastream)|Monitor activity on the ​Akamai​ platform and send live log data to a destination of your choice.|
|[Edge DNS](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-edgedns)|Replace or augment your DNS infrastructure with a cloud-based authoritative DNS solution.|
|[EdgeWorkers](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-edgeworkers)|Execute JavaScript functions at the edge to optimize site performance and customize web experiences.|
|[Global Traffic Management](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-gtm)|Use load balancing to manage website and mobile performance demands.|
|[Identity and Access Management](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-iam)|Create users and groups, and define policies that manage access to your Akamai applications.|
|[Image and Video Manager](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-ivm)|Automate image and video delivery optimizations for your website visitors.|
|[Network Lists](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-network-lists)|Automate the creation, deployment, and management of lists used in ​Akamai​ security products.|
|[Property](https://techdocs.akamai.com/terraform/v3.6/docs/set-up-property-provisioning)|Define rules and behaviors that govern your website delivery based on match criteria.|

### Get contract and group IDs

We also require a connection to your company's account, contract, and assets. This is done by adding a `contract_id` and a `group_id` to your configuration file.

Use the groups data source to get a list of the groups associated with your contracts.

```
data "akamai_groups" "my-ids" {
}

output "groups" {
value = data.akamai_groups.my-ids
}
```

The response breaks down all groups available to a particular contract. The values can be used directly or as variables. 

- Use them in a separate [variables file](https://developer.hashicorp.com/terraform/language/values/variables).
- Use them as [local variables](https://developer.hashicorp.com/terraform/language/values/locals) in the same configuration file.

## Test and go live

Check your configuration as you work by using a resource's data source and then activate your configuration in the stage environment for end-to-end testing.

When you're satisfied with the way things are running, activate your config for the production environment to analyze and respond to user requests.

To activate your configuration:

1. Run `terraform plan` to check your syntax.
2. Run `terraform apply` to execute all the actions in your configurations.

!> If multiple people are making changes to same configuration file, be mindful that running Terraform's apply command to activate your configuration will overwrite others' changes.
