cd /workspace/terraform-provider-akamai
if git diff origin/master -- go.mod | grep -E "\+\sgithub.com/akamai/AkamaiOPEN-edgegrid-golang/" &&
 sed -En 's:\t+(github.com/akamai/AkamaiOPEN-edgegrid-golang/v.+) (v.+)$:\1@\2:p' go.mod | xargs go mod download; then
  echo "New EdgeGrid version found in go.mod and it released on GitHub, go.sum should be up-to-date"
else
  echo "No new EdgeGrid version found in go.mod or it's not yet public, cleaning go.sum and use local EdgeGrid"
  go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v14=../akamaiopen-edgegrid-golang/
  go mod tidy # ensure no changes for goreleaser before hook - it enough to check go.mod and go.sum only once it's complete (so it has correct edgegrid-golang entry)
  #  fake commit to ensure clean git state
  git config --global user.email "you@example.com"
  git config --global user.name "Your Name"
  git diff --quiet go.sum go.mod || (git add go.sum go.mod && git commit -m "ensure clean git state")
fi
git tag v999.0.0
goreleaser build --single-target --config ./.goreleaser.yml --output /root/.terraform.d/plugins/registry.terraform.io/akamai/akamai/999.0.0/linux_amd64/terraform-provider-akamai_v999.0.0
