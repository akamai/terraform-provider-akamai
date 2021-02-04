---
layout: "akamai"
page_title: "Akamai: DNS Record"
subcategory: "DNS"
description: |-
  DNS record
---

# akamai_dns_record

Use the `akamai_dns_record` resource to configure a DNS record that can integrate with your existing DNS infrastructure.

## Example usage

Basic usage:

```hcl
### A Record Example
resource "akamai_dns_record" "origin" {
    zone = "origin.org"
    name = "origin.example.org"
    recordtype =  "A"
    active = true
    ttl =  30
    target = ["192.0.2.42"]
}

### CNAME Record Example
resource "akamai_dns_record" "www" {
    zone = "example.com"
    name = "www.example.com"
    recordtype =  "CNAME"
    active = true
    ttl =  600 
    target = "origin.example.org.edgesuite.net"
}
```

## Argument reference

This resource supports these arguments:

* `name` - (Required) The DNS record name. This is the node this DNS record is associated with. Also known as an owner name. 
* `zone` - (Required) The domain zone, including any nested subdomains.  
* `recordType` - (Required) The DNS record type.  
* `ttl` - (Required, Boolean) The time to live. This is a 32-bit signed integer for the time the resource record is cached. <br /> A value of `0` means that the resource record is not cached. It's only used for the transaction in progress and may be useful for extremely volatile data.  

## Required fields per record type

This section lists required arguments by record type. These are needed in addition to the ones listed in the argument reference above.


### A record

An A record requires this argument:

* target - One or more IPv4 addresses, for example, 1.2.3.4.

### AAAA record

An AAAA record requires this argument:

* target - One or more IPv6 addresses, for example, 2001:0db8::ff00:0042:8329.

### AFSDB record

An A record requires these arguments:

* target - The domain name of the host having a server for the cell named by the owner name of the resource record.
* subtype- An integer between 0 and 65535, indicating the type of service provided by the host.

### AKAMAICDN Record

The following field is required:

* target - DNS name representing selected Edge Hostname name+domain.

### AKAMAITLC Record

No additional fields are required. The following fields are Computed.

* dns_name - valid DNS name.
* answer_type - answer type.

### CAA Record

The following field are required:

* target - One or more CA Authorizations. Each authorization contains three attributes: flags, property tag and property value.

Example:
```
target = ["0 issue \"caa1.example.net\"", "0 issuewild \"ca2.example.org\"", "0 issue ca1.example.net"]
```

### CERT Record

The following fields are required:

* type_value - numeric certificate type value
* type_mnemonic - mnemonic certificate type value.
* keytag - value computed for the key embedded in the certificate
* algorithm - identifies the cryptographic algorithm used to create the signature.
* certificate - certificate data

Note: Type can be configured either as a numeric OR menmonic value. With both set, type_mnemonic takes precedence.

### CNAME Record

The following field is required:

* target - A domain name that specifies the canonical or primary name for the owner. The owner name is an alias.

### DNSKEY Record

The following fields are required:

* flags
* protocol - Must have the value 3. The DNSKEY resource record must be treated as invalid during signature verification if it contains a value other than 3.
* algorithm - The public key’s cryptographic algorithm that determines the format of the public key field.
* key - Base 64 encoded value representing the public key, the format of which depends on the algorithm.

### DS Record

The following fields are required:

* keytag - The key tag of the DNSKEY resource record referred to by the DS record, in network byte order.
* algorithm - The algorithm number of the DNSKEY resource record referred to by the DS record.
* digest_type - Identifies the algorithm used to construct the digest.
* digest - The base 16 encoded DS record refers to a DNSKEY RR by including a digest of that DNSKEY RR. The digest is calculated by concatenating the canonical form of the fully qualified owner name of the DNSKEY RR with the DNSKEY RDATA, and then applying the digest algorithm.

### HINFO Record

The following fields are required:

* hardware - Type of hardware the host uses. A machine name or CPU type may be up to 40 characters taken from the set of uppercase letters, digits, and the two punctuation characters hyphen and slash. It must start with a letter, and end with a letter.
* software - Type of software the host uses. A system name may be up to 40 characters taken from the set of uppercase letters, digits, and the two punctuation characters hyphen and slash. It must start with a letter, and end with a letter or digit.

### LOC Record

The following field is required:

* target - A geographical location associated with a domain name.

### MX Record

The following field is required:

* target - One or more domain names that specify a host willing to act as a mail exchange for the owner name.

The following fields are optional depending on configuration type. See [DNS Getting Started Guide](../guides/get_started_dns_zone#working-with-mx-records) for more information.

* priority - The preference value given to the MX record among MX records. When a mailer needs to send mail to a certain DNS domain, it first contacts a DNS server for that domain and retrieves all the MX records. It then contacts the mailer with the lowest preference value. Ignored if an embedded priority exists in the target.
* priority_increment - auto priority increment when multiple targets are provided with no embedded priority.

### NAPTR Record

The following fields are required:

* order - A 16-bit unsigned integer specifying the order in which the NAPTR records MUST be processed to ensure the correct ordering of rules. Low numbers are processed before high numbers, and once a NAPTR is found whose rule “matches” the target, the client MUST NOT consider any NAPTRs with a higher value for order (except as noted below for the flagsnapter field).
* preference - A 16-bit unsigned integer that specifies the order in which NAPTR records with equal order values should be processed with low numbers being processed before high numbers.
* flagsnaptr - A <character-string> containing flags to control aspects of the rewriting and interpretation of the fields in the record. Flags are single characters from the set [A-Z0-9]. The case of the alphabetic characters is not significant.
* service - Specifies the services available down this rewrite path.
* regexp - A String containing a substitution expression that is applied to the original string held by the client in order to construct the next domain name to lookup.
* replacement - The next NAME to query for NAPTR, SRV, or address records depending on the value of the flags field. This MUST be a fully qualified domain name.

### NS Record

The following field is required:

* target - One or more domain names that specify authoritative hosts for the specified class and domain.

### NSEC3 Record

The following fields are required:

* algorithm - The cryptographic hash algorithm used to construct the hash-value.
* flags - Eight (8) one-bit flags that can be used to indicate different processing. All undefined flags must be zero.
* iterations - The number of additional times the hash function has been performed.
* salt - The base 16 encoded salt value, which is appended to the original owner name before hashing in order to defend against pre-calculated dictionary attacks.
* next_hashed_owner_name - Base 32 encoded. The next hashed owner name in hash order. This value is in binary format. Given the ordered set of all hashed owner names, the Next Hashed Owner Name field contains the hash of an owner name that immediately follows the owner name of the given NSEC3 RR.
* type_bitmaps - The resource record set types that exist at the original owner name of the NSEC3 RR.

### NSEC3PARAM Record

The following fields are required:

* algorithm - The cryptographic hash algorithm used to construct the hash-value.
* flags - Eight (8) one-bit flags that can be used to indicate different processing. All undefined flags must be zero.
* iterations - The number of additional times the hash function has been performed.
* salt - The base 16 encoded salt value, which is appended to the original owner name before hashing in order to defend against pre-calculated dictionary attacks.

### PTR Record

The following field is required:

* target - A domain name that points to some location in the domain name space.

### RP Record

The following fields are required:

* mailbox - A domain name that specifies the mailbox for the responsible person.
* txt - A domain name for which TXT resource records exist.

### RRSIG Record

The following fields are required:

* type_covered - The resource record set type covered by this signature.
* algorithm - The Algorithm Number field identifies the cryptographic algorithm used to create the signature.
* original_ttl - The TTL of the covered record set as it appears in the authoritative zone.
* expiration - The end point of this signature’s validity. The signature cannot be used for authentication past this point in time.
* inception - The start point of this signature’s validity. The signature cannot be used for authentication prior to this point in time.
* keytag - The Key Tag field contains the key tag value of the DNSKEY RR that validates this signature, in network byte order.
* signer - The owner of the DNSKEY resource record who validates this signature.
* signature - The base 64 encoded cryptographic signature that covers the RRSIG RDATA and covered record set. Format depends on the TSIG algorithm in use.
* labels - The Labels field specifies the number of labels in the original RRSIG RR owner name. The significance of this field is that a validator uses it to determine whether the answer was synthesized from a wildcard. If so, it can be used to determine what owner name was used in generating the signature.

### SPF Record

The following field is required:

* target - Indicates which hosts are, and are not, authorized to use a domain name for the “HELO” and “MAIL FROM” identities.

### SRV Record

The following fields are required:

* target - The domain name of the target host.
* priority - A 16-bit integer that specifies the preference given to this resource record among others at the same owner. Lower values are preferred.
* weight - A server selection mechanism, specifying a relative weight for entries with the same priority. Larger weights should be given a proportionately higher probability of being selected. The range of this number is 0–65535, a 16-bit unsigned integer in network byte order. Domain administrators should use Weight 0 when there isn’t any server selection to do, to make the RR easier to read for humans. In the presence of records containing weights greater than 0, records with weight 0 should have a very small chance of being selected.
* port - The port on this target of this service. The range of this number is 0–65535, a 16-bit unsigned integer in network byte order.

### SSHFP Record

The following fields are required:

* algorithm - Describes the algorithm of the public key. The following values are assigned: 0 = reserved; 1 = RSA; 2 = DSS, 3 = ECDSA.
* fingerprint_type - Describes the message-digest algorithm used to calculate the fingerprint of the public key. The following values are assigned: 0 = reserved, 1 = SHA-1, 2 = SHA-256.
* fingerprint - The base 16 encoded fingerprint as calculated over the public key blob. The message-digest algorithm is presumed to produce an opaque octet string output, which is placed as-is in the RDATA fingerprint field.

### SOA Record

The following fields are required:

* name_server - The domain name of the name server that was the original or primary source of data for this zone.
* email_address - A domain name that specifies the mailbox of this person responsible for this zone.
* serial - The unsigned version number between 0 and 214748364 of the original copy of the zone.
* refresh - A time interval between 0 and 214748364 before the zone should be refreshed.
* retry - A time interval between 0 and 214748364 that should elapse before a failed refresh should be retried.
* expiry - A time value between 0 and 214748364 that specifies the upper limit on the time interval that can elapse before the zone is no longer authoritative.
* nxdomain_ttl - The unsigned minimum TTL between 0 and 214748364 that should be exported with any resource record from this zone.

### TLSA Record

The following fields are required:

* usage - specifies the provided association that will be used to match the certificate presented in the TLS handshake.
* selector - specifies which part of the TLS certificate presented by the server will be matched against the association data. 
* match_type - specifies how the certificate association is presented.
* certificate - specifies the "certificate association data" to be matched.

### TXT Record

The following field is required:

* target - One or more character strings. TXT RRs are used to hold descriptive text. The semantics of the text depends on the domain where it is found.

