provider "akamai" {
    edgerc = "~/.edgerc"
    fastdns_section = "dns"
}

resource "akamai_fastdns_zone" "test_zone" {
  hostname = "akamaideveloper.net"

  a {
    name = "web"
    ttl = 900
    active = true
    target = "1.2.3.4"
  }
  a {
    name = "www"
    ttl = 600
    active = true
    target = "5.6.7.8"
  }

  cname {
    name = "www-test"
    ttl = 600
    active = true
    target = "blog.akamaideveloper.net."
  }
}
