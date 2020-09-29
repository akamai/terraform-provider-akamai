terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {}

// The following is a valid terraform dns zone configuration containing data and resource configurations for a vast majority of valid record types..
// Names and values provided are EXAMPLES and should be replaced with actual.
// NOTE: Records associated with secure zoness are commented out.
// NOTE2: SOA and NS Records are commented out. These records should be imported after a primary zone is created.

locals {
  zone = "example_primary_test.zone"
  contract = "9-ABCDEF"
  group = "12345"
}

data "akamai_authorities_set" "ns_auths" {
  contract = local.contract
}

/*
// Requires record to pre exist ...

data "akamai_dns_record_set" "test_recset" {
  zone = local.zone
  host = "txtrecord.${local.zone}"
  record_type = "TXT"
}
*/

resource "akamai_dns_zone" "test_zone" {
    zone = local.zone
    contract = local.contract
    group = local.group
    type = "PRIMARY"
    comment =  "This is a test primary zone"
    end_customer_id = "end_customer_id"
    sign_and_serve = false
}

//
// The following SOA and NS records can be uncommented AND imported (if needed) after primary zone creation is completed
//

/*

resource "akamai_dns_record" "soa_record" {
  zone        = local.zone
  name        = local.zone
  recordtype  = "SOA"
  ttl         = 86400
  name_server = "a1-98.akam.net."
  email_address = "hostmaster.${local.zone}"
  refresh = 3600
  retry = 600
  expiry = 604800
  nxdomain_ttl = 300
  depends_on = [ akamai_dns_zone.test_zone ]
}

resource "akamai_dns_record" "ns_record" {
    zone = local.zone
    target = data.akamai_authorities_set.ns_auths.authorities 
    name = local.zone
    recordtype = "NS"
    ttl = 86400
    depends_on = [ akamai_dns_zone.test_zone ]
}

*/

resource "akamai_dns_record" "a_record" {
    zone = local.zone
    name = "a_${local.zone}"
    recordtype =  "A"
    ttl =  300
    target = ["10.0.0.2","10.0.0.3"]
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "aaaa_record" {
    zone = local.zone
    name = "ipv6record.${local.zone}"
    recordtype =  "AAAA"
    ttl =  3600
    target = ["4001:ab8:85b3:0:0:8a1e:370:7225"]
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "afsdb_record" {
    zone = local.zone
    name = "afsdb.${local.zone}"
    recordtype =  "AFSDB"
    ttl =  3600
    subtype = 4
    target = ["example.com"]
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "cname_record" {
    zone = local.zone
    name = "www.${local.zone}"
    recordtype =  "CNAME"
    ttl =  300
    target = ["api.${local.zone}"]
    depends_on = [akamai_dns_zone.test_zone]
}

// The following ** commented out** records are used for secure configurations. 

/*
resource "akamai_dns_record" "dnskey_record" {
    zone = local.zone
    name = "dnskey.${local.zone}"
    recordtype =  "DNSKEY"
    active = true
    ttl =  7200
    algorithm = 3
    flags = 257
    key = "Av//0/goGKPtaa28nQvPoUwVQ ... i/0hC+1CrmQkuuKtQt98WObuv7q8iQ=="
    protocol = 7
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "ds_record" {
    zone = local.zone
    name = "ds.${local.zone}"
    recordtype =  "DS"
    ttl =  7200
    algorithm = 7
    keytag = 30336
    digest = "909FF0B4DD66F91F56524C4F968D13083BE42380"
    digest_type = 1
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "nsec3_record" {
    zone = local.zone
    name = "qdeo8lqu4l81uo67oolpo9h0nv9l13dh.${local.zone}"
    recordtype =  "NSEC3"
    active = true
    ttl =  3600
    flags = 0
    algorithm = 1
    iterations = 1
    next_hashed_owner_name = "R2NUSMGFSEUHT195P59KOU2AI30JR96"
    salt = "EBD1E0942543A01B"
    type_bitmaps = "CNAME RRSIG"
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "nsec3param_record" {
    zone = local.zone
    name = "qnsec3param.${local.zone}"
    recordtype =  "NSEC3PARAM"
    active = true
    ttl =  3600
    flags = 0
    algorithm = 1
    iterations = 1
    salt = "EBD1E0942543A01B"
    //salt = "IVBEIMKFGA4TIMRVGQZUCMBRII======"
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "rrsig_record" {
    zone = local.zone
    name = "rrsig.${local.zone}"
    recordtype =  "RRSIG"
    expiration = "20120318104101"
    inception = "20120315094101"
    active = true
    ttl =  7200
    original_ttl =  3600
    algorithm = 7
    keytag = 63761
    signature = "909FF0B4DD66F91F56524C4F968D13083BE42380"
    signer = ".${local.zone}."
    labels = 3
    type_covered = "A"
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}

*/

resource "akamai_dns_record" "hinfo_record" {
    zone = local.zone
    name = "hinfo.${local.zone}"
    recordtype =  "HINFO"
    ttl =  7200
    hardware = "INTEL-386"
    software = "Unix"
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "loc_record" {
    zone = local.zone
    name = "location.${local.zone}"
    recordtype =  "LOC"
    ttl =  7200
    target = ["51 30 12.748 N 0 7 39.611 W 0.00m 0.00m 0.00m 0.00m"]
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "mx_record_self_contained" {
    zone = local.zone
    target = ["0 smtp-0.example.com.", "10 smtp-1.example.com."]
    name = "mx_record_self_contained.${local.zone}"
    recordtype = "MX"
    ttl = 300
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "mx_record_pri_increment" {
    zone = local.zone
    target = ["smtp-1.example.com.", "smtp-2.example.com.", "smtp-3.example.com."]
    priority = 10
    priority_increment = 10
    name = "mx_pri_increment.${local.zone}"
    recordtype = "MX"
    ttl = 900
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "mx_record_instances" {
    zone = local.zone
    name = "mx_record.example.${local.zone}"
    recordtype =  "MX"
    ttl =  500
    count = 3
    target = ["smtp-${count.index}.example.com."]
    priority = count.index*10
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "naptr_record" {
    zone = local.zone
    name = "naptrrecord.${local.zone}"
    recordtype =  "NAPTR"
    ttl =  3600
    flagsnaptr = "S"
    order = 0
    preference = 10
    regexp = "!^.*$!sip:customer-service@example.com!"
    replacement = "."
    service = "SIP+D2U"
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "ptr_record" {
    zone = local.zone
    name = "ptr.${local.zone}"
    recordtype =  "PTR"
    ttl =  300
    target = ["ptr.${local.zone}"]
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "rp_record" {
    zone = local.zone
    name = "rp.${local.zone}"
    recordtype =  "RP"
    ttl =  7200
    mailbox = "admin.example.com"
    txt = "txt.example.com"
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "spf_record" {
    zone = local.zone
    name = "spf.${local.zone}"
    recordtype =  "SPF"
    ttl =  7200
    target = ["v=spf"]
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "srv_record" {
    zone = local.zone
    name = "srv.${local.zone}"
    recordtype =  "SRV"
    ttl =  7200
    priority = 10
    weight  = 0
    port = 522
    target = ["target.${local.zone}"]
    depends_on = [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "sshfp_record" {
    zone = local.zone
    name = "sshfp.${local.zone}"
    recordtype =  "SSHFP"
    ttl =  7200
    algorithm = 2
    fingerprint_type  = 1
    fingerprint = "123456789ABCDEF67890123456789ABCDEF67890"
    target = []
    depends_on = [akamai_dns_zone.test_zone]
}


resource "akamai_dns_record" "txt_record" {
    zone = local.zone
    name = "text.${local.zone}"
    recordtype =  "TXT"
    ttl =  7200
    target = ["Hello world"]
    depends_on =  [akamai_dns_zone.test_zone]
}

resource "akamai_dns_record" "caa1" {
  zone        = akamai_dns_zone.test_zone.zone
  name        = akamai_dns_zone.test_zone.zone
  recordtype  = "CAA"
  target      = ["0 issue \"letsencrypt.org\""]
  ttl         = 2000
  depends_on  = [akamai_dns_zone.test_zone]
}

