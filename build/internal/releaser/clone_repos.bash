echo "$SSH_KNOWN_HOSTS" >/root/.ssh/known_hosts
echo "$SSH_PUB_KEY" >/root/.ssh/id_rsa.pub
echo "$SSH_PRV_KEY" >/root/.ssh/id_rsa
chmod 700 /root/.ssh
chmod 600 /root/.ssh/id_rsa
chmod 644 /root/.ssh/id_rsa.pub
chmod 644 /root/.ssh/known_hosts
cd /workspace
git clone ssh://git@git.source.akamai.com:7999/devexp/akamaiopen-edgegrid-golang.git
git clone ssh://git@git.source.akamai.com:7999/devexp/terraform-provider-akamai.git
