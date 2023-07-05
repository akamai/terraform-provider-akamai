#!/usr/bin/env bash

VERSION="1.4.6"

[[ -n $(which terraform) ]] && echo "Terraform already installed" && exit 0

echo "Installing terraform $VERSION"
curl -fSL "https://releases.hashicorp.com/terraform/${VERSION}/terraform_${VERSION}_linux_amd64.zip" -o terraform.zip
sudo unzip terraform.zip -d /opt/terraform
sudo ln -s /opt/terraform/terraform /usr/bin/terraform
rm -f terraform.zip