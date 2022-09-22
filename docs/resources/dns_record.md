---
layout: akamai
subcategory: Edge DNS
---

# akamai_dns_record

Use the `akamai_dns_record` resource to configure a DNS record that can integrate with your existing DNS infrastructure.

## Example usage

Here are examples of an A record and a CNAME record.

### An A record example

```
resource "akamai_dns_record" "origin" {
    zone = "origin.org"
    name = "origin.example.org"
    recordtype =  "A"
    active = true
    ttl =  30
    target = ["192.0.2.42"]
}
```

### CNAME Record Example

```
resource "akamai_dns_record" "www" {
    zone = "example.com"
    name = "www.example.com"
    recordtype =  "CNAME"
    active = true
    ttl =  600
    target = "origin.example.org.edgesuite.net"
}
```

## Argument reference [argument-reference]

This resource supports these arguments for all record types:

* `name` - (Required) The DNS record name. This is the node this DNS record is associated with. Also known as an owner name.
* `zone` - (Required) The domain zone, including any nested subdomains.  
* `recordType` - (Required) The DNS record type.  
* `ttl` - (Required) The time to live (TTL) is a 32-bit signed integer for the time the resource record is cached. <br /> A value of `0` means that the resource record is not cached. It's only used for the transaction in progress and may be useful for extremely volatile data.  

## Additional arguments by record type

This section lists additional required and optional arguments for specific record types.


### A record

An A record requires this argument:

* `target` - One or more IPv4 addresses, for example, 192.0.2.0.

### AAAA record

An AAAA record requires this argument:

* `target` - One or more IPv6 addresses, for example, 2001:0db8::ff00:0042:8329.

### AFSDB record

An AFSDB record requires these arguments:

* `target` - The domain name of the host having a server for the cell named by the owner name of the resource record.
* `subtype` - An integer between `0` and `65535` that indicates the type of service provided by the host.

### AKAMAICDN record

An AKAMAICDN record requires this argument:

* `target` - A DNS name representing the selected edge hostname and domain.

### AKAMAITLC record

No additional arguments are needed for AKAMAITLC records. This resource returns these computed attributes for this record type:

* `dns_name` - A valid DNS name.
* `answer_type` - The answer type.

### CAA record

A certificate authority authorization (CAA) record requires this argument:

* `target` - One or more certificate authority authorizations. Each authorization contains three attributes: flags, property tag, and property value.

Example:

```
target = ["0 issue \"caa1.example.net\"", "0 issuewild \"ca2.example.org\"", "0 issue ca1.example.net"]
```

### CERT record

A CERT record requires these arguments:

* `type_value` - A numeric certificate type value.
* `type_mnemonic` - A mnemonic certificate type value.
* `keytag` - A value computed for the key embedded in the certificate.
* `algorithm` - The cryptographic algorithm used to create the signature.
* `certificate` - Certificate data.

> **Note:** When entering the certificate type, you can enter `type_value`, `type_mnemonic`, or  both arguments. If you use both, `type_mnemonic` takes precedence.

### CNAME record

A CNAME record requires this argument:

* `target `- A domain name that specifies the canonical or primary name for the owner. The owner name is an alias.

### DNSKEY record

A DNSKEY record requires these arguments:

* `flags`
* `protocol` - Set to `3`. If the value isn't `3`, the DNSKEY resource record is treated as invalid during signature verification.
* `algorithm` - The public key's cryptographic algorithm. This algorithm determines the format of the public key field.
* `key` - A Base64 encoded value representing the public key. The format used depends on the `algorithm`.

### DS record

A DS record requires these arguments:

* `keytag` - The key tag of the DNSKEY record that the DS record refers to, in network byte order.
* `algorithm` - The algorithm number of the DNSKEY resource record referred to by the DS record.
* `digest_type` - Identifies the algorithm used to construct the digest.
* `digest` - A base 16 encoded DS record includes a digest of the DNSKEY record it refers to. The digest is conifgured the canonical form of the DNSKEY record's fully qualified owner name with the DNSKEY RDATA, and then applying the digest algorithm.

### HINFO record

A HINFO record requires these arguments:

* `hardware` - The type of hardware the host uses. A machine name or CPU type may be up to 40 characters long and include uppercase letters, digits, hyphens, and slashes. The entry needs to start and to end with an uppercase letter.
* `software` - The type of software the host uses. A system name may be up to 40 characters long and include uppercase letters, digits, hyphens, and slashes. The entry needs to start with an uppercase letter and end with an uppercase letter or a digit.

### HTTPS Record

The following fields are required:

* `svc_priority` - Service priority associated with endpoint. Value mist be between 0 and 65535. A piority of 0 enables alias mode.
* `svc_params` - Space separated list of endpoint parameters. Not allowed if service priority is 0.
* `target_name` - Domain name of the service endpoint.

### LOC record

A LOC record requires this argument:

* `target` - A geographical location associated with a domain name.

### MX record

An MX record supports these arguments:

* `target` - (Required) One or more domain names that specify a host willing to act as a mail exchange for the owner name.
* `priority` - (Optional) The preference value given to this MX record in relation to all other MX records. When a mailer needs to send mail to a certain DNS domain, it first contacts a DNS server for that domain and retrieves all the MX records. It then contacts the mailer with the lowest preference value. This value is ignored if an embedded priority exists in the target.
* `priority_increment` - (Optional) An auto priority increment when multiple targets are provided with no embedded priority.

See [Working with MX records](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/edge-dns#working-with-mx-records) in the [DNS Getting Started Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/edge-dns) for more information.

### NAPTR record

An NAPTR record requires these arguments:

* `order` - A 16-bit unsigned integer specifying the order in which the NAPTR records need to be processed to ensure the correct ordering of rules. Low numbers are processed before high numbers. Once a NAPTR is found whose rule matches the target, the client shouldn't consider any NAPTRs with a higher value for order (except as noted below for the flagsnapter field).
* `preference` - A 16-bit unsigned integer that specifies the order in which NAPTR records with equal order values are processed. Low numbers are processed before high numbers.
* `flagsnaptr` - A character string containing flags that control how fields in the record are rewritten and interpreted. Flags are single alphanumeric characters.
* `service` - Specifies the services available down this rewrite path.
* `regexp` - A regular expression string containing a substitution expression. This substitution expression is applied to the original client string in order to construct the next domain name to lookup.
* `replacement` - Depending on the value of the `flags` attribute, the next NAME to query for NAPTR, SRV, or address records. Enter a fully qualified domain name as the value.

### NS record

An NS record requires these arguments:

* `target` - One or more domain names that specify authoritative hosts for the specified class and domain.

### NSEC3 record

An NSEC3 record requires these arguments:

* `algorithm` - The cryptographic hash algorithm used to construct the hash-value.
* `flags` - Eight one-bit flags you can use to indicate different processing. All undefined flags must be zero.
* `iterations` - The number of additional times the hash function has been performed.
* `salt` - The base 16 encoded salt value, which is appended to the original owner name before hashing. Used to defend against pre-calculated dictionary attacks.
* `next_hashed_owner_name` - Base 32 encoded. The next hashed owner name in hash order. This value is in binary format. Given the ordered set of all hashed owner names, the Next Hashed Owner Name field contains the hash of an owner name that immediately follows the owner name of the given NSEC3 RR.
* `type_bitmaps` - The resource record set types that exist at the original owner name of the NSEC3 RR.

### NSEC3PARAM record

An NSEC3PARAM record requires these arguments:

* `algorithm` - The cryptographic hash algorithm used to construct the hash-value.
* `flags` - Eight one-bit flags that can be used to indicate different processing. All undefined flags must be zero.
* `iterations` - The number of additional times the hash function has been performed.
* `salt` - The base 16 encoded salt value, which is appended to the original owner name before hashing in order to defend against pre-calculated dictionary attacks.

### PTR record

A PTR record requires this argument:

* `target` - A domain name that points to some location in the domain name space.

### RP record

An RP record requires these arguments:

* `mailbox` - A domain name that specifies the mailbox for the responsible person.
* `txt` - A domain name for which TXT resource records exist.

### RRSIG record

An RRSIG record requires these arguments:

* `type_covered` - The resource record set type covered by this signature.
* `algorithm` - Identifies the cryptographic algorithm used to create the signature.
* `original_ttl` - The TTL of the covered record set as it appears in the authoritative zone.
* `expiration` - The end point of this signature's validity. The signature can`t be used for authentication past this point in time.
* `inception` - The start point of this signature's validity. The signature can`t be used for authentication prior to this point in time.
* `keytag` - The Key Tag field contains the key tag value of the DNSKEY RR that validates this signature, in network byte order.
* `signer` - The owner of the DNSKEY resource record who validates this signature.
* `signature` - The base 64 encoded cryptographic signature that covers the RRSIG RDATA and covered record set. Format depends on the TSIG algorithm in use.
* `labels` - The Labels field specifies the number of labels in the original RRSIG RR owner name. The significance of this field is that a validator uses it to determine whether the answer was synthesized from a wildcard. If so, it can be used to determine what owner name was used in generating the signature.

### SPF record

An SPF record requires this argument:

* `target` - Indicates which hosts are, and are not, authorized to use a domain name for the “HELO” and “MAIL FROM” identities.

### SRV record

An SRV record requires these arguments:

* `target` - The domain name of the target host.
* `priority` - A 16-bit integer that specifies the preference given to this resource record among others at the same owner. Lower values are preferred.
* `weight` - A server selection mechanism that specifies a relative weight for entries with the same priority. Larger weights are given a proportionately higher probability of being selected. The range of this number is 0–65535, a 16-bit unsigned integer in network byte order. Domain administrators should use Weight 0 when there isn't any server selection to do, to make the RR easier to read for humans. In the presence of records containing weights greater than 0, records with weight 0 should have a very small chance of being selected.
* `port` - The port on this target of this service. The range of this number is 0–65535, a 16-bit unsigned integer in network byte order.

### SSHFP record

An SSHFP record requires these arguments:

* `algorithm` - Describes the algorithm of the public key. The following values are assigned: `0` is reserved, `1` is for RSA, `2` is for DSS, and `3` is for ECDSA.
* `fingerprint_type` - Describes the message-digest algorithm used to calculate the fingerprint of the public key. The following values are assigned: 0 = reserved, 1 = SHA-1, 2 = SHA-256.
* `fingerprint` - The base 16 encoded fingerprint as calculated over the public key blob. The message-digest algorithm is presumed to produce an opaque octet string output, which is placed as-is in the RDATA fingerprint field.

### SOA record

An SOA record requires these arguments:

* `name_server` - The domain name of the name server that was the original or primary source of data for this zone.
* `email_address` - A domain name that specifies the mailbox of this person responsible for this zone.
* `serial` - The unsigned version number between 0 and 214748364 of the original copy of the zone.
* `refresh` - A time interval between 0 and 214748364 before the zone should be refreshed.
* `retry` - A time interval between 0 and 214748364 that should elapse before a failed refresh should be retried.
* `expiry` - A time value between 0 and 214748364 that specifies the upper limit on the time interval that can elapse before the zone is no longer authoritative.
* `nxdomain_ttl` - The unsigned minimum TTL between 0 and 214748364 that should be exported with any resource record from this zone.

### SVCB record

An SVCB record requires these arguments:

* `svc_priority` - Service priority associated with endpoint. Value mist be between 0 and 65535. A piority of 0 enables alias mode.
* `svc_params` - Space separated list of endpoint parameters. Not allowed if service priority is 0.
* `target_name` - Domain name of the service endpoint.

### TLSA record

A TLSA record requires these arguments:

* `usage` - Specifies the association used to match the certificate presented in the TLS handshake.
* `selector` - Specifies the part of the TLS certificate presented by the server that is matched against the association data.
* `match_type` - Specifies how the certificate association is presented.
* `certificate` - Specifies the certificate association data to be matched.

### TXT record

A TXT record requires this argument:

* `target` - One or more character strings. TXT resource records hold descriptive text. The semantics of the text depends on the domain where it is found.

