# Managing dev, preprod, prod with Terraform

It's good practice to have an Akamai Property Manager configuration for each of your environments so that you can develop new rule sets in dev, then promote to preprod for QA testing before you push to prod. However, it's tricky to keep these environments synchronized when doing this manually in the GUI. Terraform, however, makes this simple. Here is just one example of how this can be done using Terraform modules. There are undoubtedly other ways to do this.

Firstly, we need somewhere to hold the property module, so we create a subdirectory "modules/property". The main terraform configuration files will live in there along with the Property Manager rules. The main Terraform configuration would look very similar to what you would do for a single configuration, except that any parameter that defines the environment needs to be a variable. Note, you don't need to worry about the names of the resources because in this case we're going to be keeping a state file per environment.

For example, here we define a CPCode specific for the environment

```
resource "akamai_cp_code" "test-wheep-co-uk" {
 product_id  = "prd_Download_Delivery"
 contract_id = data.akamai_contract.contract.id
 group_id    = data.akamai_group.group.id
 name        = "${var.env}.wheep.co.uk"
}
```

and here we define an edge hostname specific for the environment

```
resource "akamai_edge_hostname" "test-wheep-co-uk-edgesuite-net" {
 product_id    = "prd_Download_Delivery"
 contract_id   = data.akamai_contract.contract.id
 group_id      = data.akamai_group.group.id
 ip_behavior   = "IPV6_COMPLIANCE"
 edge_hostname = "${var.env}.wheep.co.uk.edgesuite.net"
}
```

Etc, etc.

We then need to define our variables in "variable.tf", without giving default values because at this point, there's no sensible defaults.

Once the module is created, we can use it by defining a directory for each environment.

```
environments/dev
environments/preprod
environments/prod
```

We can then place a "main.tf" in each environment directory that references the module and supplies the relevant variables.

For example

```
module "property" {
        source = "../../modules/property"
        env = "dev-envtest"
        network = "staging"
}
```

The workflow would be as follows:-

1) change the rules.tf in the module
2) Apply the changes in your dev environment
3) When you're ready, apply in your preprod env
4) And finally apply in your prod env once everything looks good

