#!/usr/bin/env bash

# This script executes a "dry run" simulating what the goreleaser GitHub action does for every release
# It is executed on last commit from develop branch.

SSH_PRV_KEY="$(cat ~/.ssh/id_rsa)"
SSH_PUB_KEY="$(cat ~/.ssh/id_rsa.pub)"
SSH_KNOWN_HOSTS="$(cat ~/.ssh/known_hosts)"
EDGERC="$(cat ~/.edgerc)"

echo "building image:"
cd build/internal/releaser
docker build -f Dockerfile -t "releaser" .

echo "starting container:"
docker container run --name releaser_container -d -it \
  -e SSH_PUB_KEY="${SSH_PUB_KEY}" \
  -e SSH_PRV_KEY="${SSH_PRV_KEY}" \
  -e SSH_KNOWN_HOSTS="${SSH_KNOWN_HOSTS}" \
  -e EDGERC="${EDGERC}" \
  --add-host git.source.akamai.com:100.78.0.6 "releaser"

echo "cloning repositories and executing goreleaser:"
docker exec releaser_container bash -c 'echo "$SSH_KNOWN_HOSTS" > /root/.ssh/known_hosts;
          echo "$SSH_PUB_KEY" > /root/.ssh/id_rsa.pub;
          echo "$SSH_PRV_KEY" > /root/.ssh/id_rsa;
          echo "$EDGERC" > /root/.edgerc;
          chmod 700 /root/.ssh;
          chmod 600 /root/.ssh/id_rsa;
          chmod 644 /root/.ssh/id_rsa.pub;
          chmod 644 /root/.ssh/known_hosts;
          cd /workspace;
          git clone ssh://git@git.source.akamai.com:7999/devexp/akamaiopen-edgegrid-golang.git;
          git clone ssh://git@git.source.akamai.com:7999/devexp/terraform-provider-akamai.git;
          cd terraform-provider-akamai;
          go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v2=../akamaiopen-edgegrid-golang/;
          git tag v10.0.0;
          goreleaser build --single-target --skip-validate --config ./.goreleaser.yml --output /root/.terraform.d/plugins/registry.terraform.io/akamai/akamai/10.0.0/linux_amd64/terraform-provider-akamai_v10.0.0'

echo "smoke test:"
docker exec -w '/workspace/terraform-provider-akamai/examples/akamai_cp_code' releaser_container bash -c '
          terraform init;
          terraform plan'

echo "cleaning up:"
docker container rm -f releaser_container
docker image rm releaser