provider "akamai" {
    edgerc = "/root/.edgerc"
    dns_section = "dns"
}

locals {
  zone = "akavaiodeveloper.net"
}

resource "akamai_dnsv2_zone" "test_zone" {
    contractid = "C-1FRYVV3"
    zone = local.zone
    #type = "SECONDARY"
    masters = ["1.2.3.4" , "1.2.3.5"]
    type = "PRIMARY"
    comment =  "This is a test zone"
    gid     = "64867"
    sign_and_serve = false
}

resource "akamai_dnsv2_record" "a_record" {
    zone = local.zone
    name = "akavaiodeveloper.net"
    recordtype =  "A"
    active = true
    ttl =  300
    target = ["10.0.0.2","10.0.0.3"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "aaaa_record" {
    zone = local.zone
    name = "ipv6record.akavaiodeveloper.net"
    recordtype =  "AAAA"
    active = true
    ttl =  3600
    target = ["2001:0db8::ff00:0042:8329"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "afsdb_record" {
    zone = local.zone
    name = "afsdb.akavaiodeveloper.net"
    recordtype =  "AFSDB"
    active = true
    ttl =  3600
    subtype = 4
    target = ["example.com"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}



resource "akamai_dnsv2_record" "cname_record" {
    zone = local.zone
    name = "www.akavaiodeveloper.net"
    recordtype =  "CNAME"
    active = true
    ttl =  300
    target = ["api.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "dnskey_record" {
    zone = local.zone
    name = "dnskey.akavaiodeveloper.net"
    recordtype =  "DNSKEY"
    active = true
    ttl =  7200
    algorithm = 3
    flags = 257
    key = "Av//0/goGKPtaa28nQvPoUwVQ ... i/0hC+1CrmQkuuKtQt98WObuv7q8iQ=="
    protocol = 7
    target = ["dnskey.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "ds_record" {
    zone = local.zone
    name = "ds.akavaiodeveloper.net"
    recordtype =  "DS"
    active = true
    ttl =  7200
    algorithm = 7
    keytag = 30336
    digest = "909FF0B4DD66F91F56524C4F968D13083BE42380"
    digest_type = 1
    target = ["dnskey.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "hinfo_record" {
    zone = local.zone
    name = "hinfo.akavaiodeveloper.net"
    recordtype =  "HINFO"
    active = true
    ttl =  7200
    hardware = "INTEL-386"
    software = "Unix"
    target = ["hinfo.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "loc_record" {
    zone = local.zone
    name = "location.akavaiodeveloper.net"
    recordtype =  "LOC"
    active = true
    ttl =  7200
    target = ["51 30 12.748 N 0 7 39.611 W 0.00m 0.00m 0.00m 0.00m"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}


resource "akamai_dnsv2_record" "mx_record" {
    count = 3
    zone = local.zone
    name = "akavaiodeveloper.net"
    recordtype =  "MX"
    active = true
    ttl =  300
    target = ["smtp-${count.index}.akavaiodeveloper.net"]
    priority = "${count.index*10}"
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "naptr_record" {
    zone = local.zone
    name = "naptrrecord.akavaiodeveloper.net"
    recordtype =  "NAPTR"
    active = true
    ttl =  3600
    flagsnaptr = "S"
    order = 0
    preference = 10
    regexp = "!^.*$!sip:customer-service@example.com!"
    replacement = "."
    service = "SIP+D2U"
    target = ["naptr.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "ns_record" {
    zone = local.zone
    name = "ns.akavaiodeveloper.net"
    recordtype =  "NS"
    active = true
    ttl =  300
    target = ["use4.akam.net."]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "nsec3_record" {
    zone = local.zone
    name = "qdeo8lqu4l81uo67oolpo9h0nv9l13dh.akavaiodeveloper.net"
    recordtype =  "NSEC3"
    active = true
    ttl =  3600
    flags = 0
    algorithm = 1
    iterations = 1
    next_hashed_owner_name = "R2NUSMGFSEUHT195P59KOU2AI30JR96"
    salt = "EBD1E0942543A01B"
    type_bitmaps = "CNAME RRSIG"
    target = ["naptr.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "nsec3param_record" {
    zone = local.zone
    name = "qnsec3param.akavaiodeveloper.net"
    recordtype =  "NSEC3PARAM"
    active = true
    ttl =  3600
    flags = 0
    algorithm = 1
    iterations = 1
    salt = "EBD1E0942543A01B"
    //salt = "IVBEIMKFGA4TIMRVGQZUCMBRII======"
    target = ["qnsec3param.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "ptr_record" {
    zone = local.zone
    name = "ptr.akavaiodeveloper.net"
    recordtype =  "PTR"
    active = true
    ttl =  300
    target = ["ptr.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}


resource "akamai_dnsv2_record" "rp_record" {
    zone = local.zone
    name = "rp.akavaiodeveloper.net"
    recordtype =  "RP"
    active = true
    ttl =  7200
    mailbox = "admin.example.com."
    txt = "txt.example.com"
    target = ["txt.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

/*
resource "akamai_dnsv2_record" "rrsig_record" {
    zone = local.zone
    name = "rrsig.akavaiodeveloper.net"
    recordtype =  "RRSIG"
    expiration = "20120318104101"
    inception = "20120315094101"
    active = true
    ttl =  7200
    original_ttl =  3600
    algorithm = 7
    keytag = 63761
    signature = "909FF0B4DD66F91F56524C4F968D13083BE42380"
    signer = "akavaiodeveloper.net."
    labels = 3
    type_covered = "A"
    target = ["dnskey.akavaiodeveloper.net"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}
*/

resource "akamai_dnsv2_record" "spf_record" {
    zone = local.zone
    name = "spf.akavaiodeveloper.net"
    recordtype =  "PTR"
    active = true
    ttl =  7200
    target = ["v=spf"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "srv_record" {
    zone = local.zone
    name = "srv.akavaiodeveloper.net"
    recordtype =  "SRV"
    active = true
    ttl =  7200
    priority = 10
    weight  = 0
    port = 522
    target = ["target.akavaiodeveloper.net."]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

resource "akamai_dnsv2_record" "sshfp_record" {
    zone = local.zone
    name = "sshfp.akavaiodeveloper.net"
    recordtype =  "SSHFP"
    active = true
    ttl =  7200
    algorithm = 2
    fingerprint_type  = 1
    fingerprint = "123456789ABCDEF67890123456789ABCDEF67890"
    target = ["sshfp.akavaiodeveloper.net."]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}


resource "akamai_dnsv2_record" "txt_record" {
    zone = local.zone
    name = "text.akavaiodeveloper.net"
    recordtype =  "TXT"
    active = true
    ttl =  7200
    target = ["Hello world"]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

/*
resource "akamai_dnsv2_record" "ns_record" {
    zone = local.zone
    name = "akavaiodeveloper.net"
    recordtype =  "NS"
    active = true
    ttl =  3600
    target = ["use4.akam.net."]
    depends_on =  ["akamai_dnsv2_zone.test_zone"]
}
*/


data "akamai_authorities_set" "ns" {
  contractid = "C-1FRYVV3"
}


output "authorities" {
  value = "${join(",", data.akamai_authorities_set.ns.authorities)}"
}

/*
data "akamai_dns_record_set" "rs" {
  zone = "akavaiodeveloper.net"
  host = "akavaiodeveloper.net"
  record_type = "A"
  depends_on =  ["akamai_dnsv2_zone.test_zone"]
}

output "rs_addrs" {
  value = "${join(",", data.akamai_dns_record_set.rs.addrs)}"
}
*/
data "akamai_dns_record_set" "rs1" {
  zone = "akamaideveloper.com"
  host = "k1._domainkey.akamaideveloper.com"
  record_type = "CNAME"
}

output "rs1_addrs" {
  value = "${join(",", data.akamai_dns_record_set.rs1.rdata)}"
}

data "akamai_dns_record_set" "mx" {
  zone = "akamaideveloper.com"
  host = "akamaideveloper.com"
  record_type = "MX"
}

output "mx_addrs" {
  value = "${join(",", data.akamai_dns_record_set.mx.rdata)}"
}
