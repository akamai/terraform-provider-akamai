cd /workspace/terraform-provider-akamai
go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v11=../akamaiopen-edgegrid-golang/
git tag v10.0.0
goreleaser build --single-target --skip=validate --config ./.goreleaser.yml --output /root/.terraform.d/plugins/registry.terraform.io/akamai/akamai/10.0.0/linux_amd64/terraform-provider-akamai_v10.0.0
