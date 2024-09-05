#!/usr/bin/env bash

VERSION="${TERRAFORM_VERSION:-1.9.5}"
VERSION="${VERSION#v}"

if [[ -n $(which terraform) && "$(terraform --version | sed 1q | cut -f2 -d" " | cut -c2-)" == "$VERSION" ]]; then
    echo "Terraform $VERSION is installed" && exit 0
fi

echo "Installing terraform $VERSION"
curl -fSL "https://releases.hashicorp.com/terraform/${VERSION}/terraform_${VERSION}_$(go env GOOS)_$(go env GOARCH).zip" -o terraform.zip
unzip terraform.zip -d /usr/local/bin
rm -f terraform.zip