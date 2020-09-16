#!/usr/bin/env bash

# Check gofmt
echo "==> Checking for unchecked errors..."

if ! which errcheck > /dev/null; then
    echo "==> Installing errcheck..."
    go get -u github.com/kisielk/errcheck
fi

err_files=$($GOBIN/errcheck -ignoretests \
                     -ignore 'github.com/hashicorp/terraform/helper/schema:Set' \
                     -ignore 'bytes:.*' \
                     -ignore 'io:Close|Write' \
                     $(go list ./...) 2>&1 )

if [[ -n ${err_files} ]]; then
    echo 'Unchecked errors found in the following places:'
    echo "${err_files}"
    echo "Please handle returned errors. You can check directly with \`make errcheck\`"
    exit 1
fi

exit 0
