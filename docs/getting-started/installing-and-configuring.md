---
layout: "akamai"
page_title: "Installing and Configuring Terraform and the Akamai Provider"
description: |-
  Installing and Configuring Terraform and the Akamai Provider
---


# Installing and Configuring Terraform and the Akamai Provider

Before you can use Terraform to manage your Akamai infrastructure, you'll need to download, install, and configure the Terraform executable file and the Akamai provider. Fortunately, that's remarkably easy, especially the part that might sound daunting to you: installing and configuring the Terraform executable. In fact, that might be the easiest thing you'll do all day.

To download the Terraform executable, open a web browser, navigate to https://www.terraform.io and then click Download CLI:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/download-button.png)

On the **Download Terraform** page, click the appropriate link based on the operating system of the computer where Terraform will be installed:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/pick-your-os.png)

When you click a button, a .ZIP file will be copied to the downloads folder on your local computer. 
That file will have a name similar to `terraform_1.0.4_darwin_amd64`, although the actual file name will vary depending on the most recent version of Terraform and on your operating systems. 
For example, the Solaris version of Terraform will have a filename like `terraform_1.0.4_solaris_amd64`. 
Regardless, the downloaded ZIP file will have only one file, the Terraform executable.

> **Note**. And, again, the name of *that* file will vary depending on your operating system. 
For Windows, the file will be named `Terraform.exe`; for the Mac, the file will be named `terraform`.

Unzip the ZIP file and copy the Terraform executable to a folder of your choice. Typically you'll copy the file to a folder named **Terraform**; however, you can use any folder and any folder name that you want. No matter which folder you choose, it's recommended that you add the folder to your operating system path (assuming it isn't already *in* the path). Doing that enables you to call Terraform from any folder on your computer, something that can come in handy later on.

Believe it or not, that's all you have to do to install Terraform itself. To verify the installation, open a command window and navigate to the folder where you copied the Terraform executable. From the command prompt, type the following command and then press ENTER:

```
terraform version
```

If everything has gone according to plan, you should get back a response similar to this:

```
Terraform v1.0.4
on darwin_amd64
```

As soon as Terraform is up and running, the next step is to install the Akamai Terraform provider. To do that, go to the Web page https://registry.terraform.io/providers/akamai/akamai/latest and click the **Use Provider** button. You should see a dropdown similar to this (depending, of course, on the latest version of the Akamai provider):

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/how-to-use-the-provider.png)

In the dropdown, copy the text with the gray background to the clipboard.

Now we need to create out first Terraform configuration file (we won't talk go into any details about configuration files here; instead, we have an entire article devoted to that subject). Open your favorite text editor (or any text editor capable of saving plain-text files) and paste in the copied text:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/configuration-options.png)

After you've done that, save the file to the same folder where you put the Terraform executable (e.g., the **Terraform** folder). You can give this file any filename you want, as long as the name has a **.tf** file extension. For example, we saved our file as **akamai.tf**. You're more than welcome to use that exact same filename.

To this point all we've done is create a Terraform configuration file that references the Akamai Terraform provider, a provider we haven't even installed. But that's OK: when we try running Terraform, Terraform will recognize that the provider is missing and automatically download and install it for us. To get that to happen, run the following command from the command prompt:

```
terraform init
```

You should see output similar to this:

```
Initializing the backend...

Initializing provider plugins...

Finding akamai/akamai versions matching "1.6.1"...

Installing akamai/akamai v1.6.1...

Installed akamai/akamai v1.6.1 (signed by a HashiCorp partner, key ID A26ECDD8F0BCBA73)

Partner and community providers are signed by their developers.
If you'd like to know more about provider signing, you can read about it here:
https://www.terraform.io/docs/cli/plugins/signing.html

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

If that seems like too much output to process, don't worry about it: to make a long story short, all that text simply means that the Akamai provider has been successfully installed and is ready for use. If you take a peek at your Terraform folder, you should see several new files and folders. For example, on a Windows computer your Terraform folder will now more like this:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/terraform-folder.png)

So are we ready to start using Terraform to manage our Akamai infrastructure? No, not yet. For example, suppose we add a line to our akamai.tf file, a line that calls the [akamai_appsec_contracts_groups](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_contracts_groups) data source; ideally, that will return information about our Akamai contracts and contract groups. In other words, suppose **akamai.tf** now looks like this:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  # Configuration options
}

data "akamai_appsec_contracts_groups" "contracts_and_groups" {
}
```

And suppose we run the following command from the command prompt in order to test this data source:

```
terraform plan
```

Do we back contracts and groups information? Not exactly. Instead, that command is going to fail with the following error:

```
╷
│ Error: Akamai EdgeGrid configuration was not specified. Specify the configuration using system environment variables or the location and file name containing the edgerc configuration. Default location the provider checks for is the current user's home directory. Default configuration file name the provider checks for is .edgerc.
│
│   with provider["registry.terraform.io/akamai/akamai"],
│   on akamai.tf line 10, in provider "akamai":
│   10: provider "akamai" {
│
╵
```

Admittedly, that's a somewhat-long and somewhat-obtuse error message. What it comes down to, however, is this: we didn't provide any sort of authentication credentials and, as a result, we're denied access to all Akamai-managed resources. As you might expect, Terraform or no Terraform, you have to be authenticated before you can do anything else.

If you've run any of the Akamai APIs then you have a pretty good idea of what comes next: in order to generate the required authentication credentials you need to log on to Akamai Control Center and create an API client. To do all that, log on to Control Center, click the “hamburger” menu, and then, under **Account Admin**, click **Identity & access**:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/control-center-menu.png)

On the **Identity and Access Management** page, click **Create API Client** and then click **Quick**:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/create-api-client.png)

> **Note**. Clicking **Quick** creates an API client that gives you access to all the APIs that you're allowed to access. For our purposes, that's the fastest and easiest way to create an API client. In real life, however, you might be better off clicking **Advanced**, and creating an API client that has only the permissions needed to do whatever it is you need to do. That provides at least *some* measure of safety and security should anyone else manage to get hold of that API client.

After you click **Quick** an API client is created, and client credentials are displayed onscreen:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/copy-credential.png)

These four credentials – **client_secret**, **host**, **access_token**, and **client_token** – are required in order to run your Terraform configuration file. It's highly recommended that you click **Download** to download a copy of this information; that gives you a text file (e.g., **ID_gstemp.txt**) that can be opened in any application that can open text files:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/credentials-download.png)

At the very least, don't move off the API client page until you've done *something* with the credentialling information: as soon as you leave the page you'll no longer have access to the client secret. (You'll still be able to review the host name, access token, and client token, but the client secret will be gone for good. If you forget or misplace that secret, the only thing can do is repeat the process and create a new API client, with a new client secret.)

So what do you do with these credentials after you have them? There are at least two possibilities. For one, you can hardcode the credentials into each of your Terraform configuration files. You might recall that, at the moment, our configuration file includes the following block:

```
provider "akamai" {

  # Configuration option

}
```

One option available to us is to put our credentials into the provider block. For example:

```
provider "akamai" {
  config {
    client_secret = "w0Sz+KAZCwMM7gP5q4U5n5uxb9p42of3l8kBfTY="
    host          = "akab-s3aru73-iy2hqxpqat7buswj.luna.akamaiapis.net"
    access_token  = "akab-s7shuimyiux22adh-mscbvunw7rn5rrog"
    client_token  = "akab-pcyzhq6s27amjews-squvxwo3r5hpfgps"
  }
}
```

> **Note**. We've truncated some of the preceding values to ensure that they'd fit nicely on a single line. Note, too that the individual values must be enclosed in double quotes.

That approach works  fine, although there are a couple of drawbacks. For one, you'll need to add this credentialing information to every Terraform configuration file you write. That can be a bit tedious, to say the least. In addition, this complicates your ability to share Terraform files with your colleagues: after all, for each file you share you need to remove your credentials and your colleagues need to type in *their* credentials. Depending on how many files you share, that could be a lot of work, and a lot of opportunities for error.

Keep in mind, too, that these credentials won't last forever. When they do expire, you'll need to update each and every one of your Terraform configuration files.

So is there an alternative to hardcoding your credentials into each and every configuration file? Yes, there is. Instead of hardcoding the values into each file, create a text file and paste the credentials into the text file exactly as-is:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/edgerc-file.png)

When that's done,  save the file to your home directory using the filename **.edgerc**: note the period at the beginning of the filename and the lack of a file extension. We should add that you don't *have* to store the file in your home directory: the .edgerc file can be stored anywhere you want to store it. The value of using your home directory is that, when referencing the file, you don't have to provide any path information. That's because, by default, Akamai's Terraform provider always checks your home directory to see if it contains a .edgerc file.

> **Note**. And what if you do both, what if you have a .edgerc file in your home directory and you hardcode credential information in your Terraform configuration file? In that case, Terraform ignores the .edgerc file and uses the credentials found in the configuration file.

Oh, and if you have your .edgerc file in your home directory you can actually leave the provider block out of your configuration file: as we noted, the Akamai provider file automatically looks for the file in your home folder and, if found, uses those credentials. However, for educational reasons we'll modify our provider block to specify the location of the .edgerc file. For example:

```
provider "akamai" {
  edgerc = "~/.edgerc"
}
```

This same boilerplate provider block can be used in all your Terraform configuration files; furthermore, you can share those files with all your colleagues without: 1) exposing your credentials; and/or, 2) without requiring those colleagues to hardcode their credentials into the file. (Assuming, of course, that they have a .edgerc file in their home directory.) And what if your access token expires, or what if you need to change that access token? That's fine: all you have to do is update the .edgerc file: you won't have to change any of your Terraform configuration files.

> **Note**. True, that might not always be the case: it depends on what you're doing and how you're doing it. But it's good enough for now.

So does that mean that our Terraform configuration file will work now? Well, there's only one way to find out: let's re-run the `terraform plan` command and see what happens. Depending on how things are set up, you might get back a collection of groups and contracts; alternatively, you might get back a response similar to this:

```
No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

Either way this means you're in business: you're now ready to start writing your own Terraform configuration files.
