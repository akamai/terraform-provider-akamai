#!/usr/bin/env bash

echo "==> Checking for linting errors..."

if ! which golint > /dev/null; then
    echo "==> Installing go lint..."
    go get -u github.com/golang/lint/golint
fi

lint_files=$(golint -set_exit_status $(go list ./... | grep -v ^/vendor/))

if [[ -n ${lint_files} ]]; then
    echo 'Linting errors found in the following places:'
    echo "${lint_files}"
    echo "Please handle returned errors. You can check directly with \`make lint\`"
    exit 1
fi

exit 0
