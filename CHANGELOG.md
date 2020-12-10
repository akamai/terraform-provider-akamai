# RELEASE NOTES

## 1.0.0 (Dec 9, 2020) Provisioning redesign

#### BREAKING CHANGES:
* provider: configuring via an inline provider block (`property`, `dns`, or `gtm`) has been replaced with a more general `config` block that works the same way.
* There are several breaking changes in the 1.0 release.  You should consult the [Migration guide](docs/guides/1.0_migration.md) for details.
  * resources/akamai_property_activation no longer supports the following fields : activate.  version has gone from being optional to being a required field.
  * data-sources/akamai_property_rules removed in favor of using template JSON object to better work with other Akamai tools and documentation that is all JSON based.
  * resources/akamai_property_variables removed in favor of directly managing the variable segment as part of ruletree object.
  * resources/akamai_cp_code no longer auto-imports on create. If a conflict is detected will error out and to ignore simply import the resource.
  * resources/akamai_edge_hostname no longer supports the following fields : ipv4, ipv6. The revised resource allows setting ip_behavior directly.
  * resources/akamai_property no longer supports the following fields : cp_code, origin, variables, is_secure, contact. The revised resource simplifies the object structure and removes the ability to set the same value more than one way.
#### NOTES:
* provider/papi: changed attribute names in Provisioning to distinguish objects and names from id attributes.  In prior releases, "group" could represent a name, an id, or sometimes both. This release distinguishes them with distinct attribute names "group_name", "group_id" instead of "group"."
#### KNOWN BUGS:
* resources/akamai_property removing hostnames attribute can result in repeated noop update calls because in this case removal means the hostname relationships are un-managed leaving the attribute as empty is a better way to express this change.
#### FEATURES:
* data-sources/akamai_properties added to list properties accessible to the user.
* data-sources/akamai_property_contracts added to list contracts accessible to the user.
* data-sources/akamai_property_groups added to list groups accessible to the user.
* data-sources/akamai_property_products added to list products associated with a given contract.
* data-sources/akamai_property_rule_formats added to list rule_formats.
* data-sources/akamai_property_rules changed to output the structure of a particular rule version on the server. NOTE: this is NOT the same as the deprecated datasource used for rule formatting.
* data-sources/akamai_rules_template added to handle file based JSON templating for rules tree data management
#### ENHANCEMENTS:
* resources/akamai_property_activation aliased property to property_id. Returns these additional attributes : target_version, warnings, errors, activation_id, and status
#### BUG FIXES:
* provider: provider configuration validation requires an edgerc file configured and present even when environment variable-based configuration was used.
* provider: provider inline configuration support was re-introduced as a new config field.
* resources/akamai_property_activation activating and destroying activation for the same property multiple times in a row would fail on second destroy attempt and subsequent destroy attempts with "resource not found error" message.
* resources/akamai_property_activation wrong activation id read for property versions that had been activated and deactivated multiple time.
#### MINOR CHANGES:
* resources/akamai_property aliased property to property_id. contract to contract_id, and product to product_id and account to account_id.  Renamed version to latest_version.
* data-sources/akamai_contract aliased group to group_id and/or group_name.
* data-sources/akamai_cp_code aliased group to group_id and contract to contract_id.
* data-sources/akamai_group aliased name to group_name and contract to contract_id.

## 0.11.0 (Nov 19,2020)

#### NOTES:
* provider: Added support for application security API
#### BUG FIXES:
* provider: Updated edgegrid library to version 2.0.2. Ths should include the following fixes:
    * Re-enabled global account switch key support in edgerc files for reseller accounts.
    * PAPI - edgehostname updated returns - The System could not find cnameTo value
    * PAPI - property update return error - You provided an Etag that does not represent the last edit. Another edit has occurred, so check your request again before retrying.

## 0.10.2 (Oct 22,2020)
#### NOTES:
* Documentation formatting
#### KNOWN BUGS:
* provider: provider configuration validation requires an edgerc file configured and present even when environment variable-based configuration was used.
* provider: support for configuring the provider via an inline provider block (`property`, `dns`, or `gtm`) no longer works.  Users should use edgerc file or Terraform environment args to configure instead.

## 0.10.1 (Not released)

## 0.10.0 (Oct 20,2020)

#### NOTES:
* provider: The backing edgegrid library was entirely rewritten.  Provider behavior should be preserved but there is chance of incidental changes due to the project size.
* resources/akamai_edge_hostname: edge_hostname field should be provided with an ending of edgesuite.net, edgekey.net, or akamaized.net.  If a required suffix is not provided then edgesuite.net is appended as default.
#### KNOWN BUGS:
* provider: provider configuration validation requires an edgerc file configured and present even when one should not be needed.
* provider: support for configuring the provider via an inline provider block (`property`, `dns`, or `gtm`) no longer works.  Users should use edgerc file or Terraform environment args to configure instead.
#### ENHANCEMENTS:
* provider: improved error handling and improved message consistency
* provider: release notes categorize updates according to Terraform best practices guide.
* resources/akamai_cp: support ids with and without prefixes
* resources/akamai_edge_hostnames: support ids with and without prefixes
* resources/akamai_property: support ids with and without prefixes
* resources/akamai_property_activation: support ids with and without prefixes
#### BUG FIXES:
* resources/akamai_property: [AT-42] Fix criteria_match values handling
* provider: fixed documentation to properly present guides and categories on Hashicorp Terraform registry site
* resources/edge_hostname: added error when neither IPV4 nor IPV6 is selected
* resources/akamai_property: comparisons in rule tree now properly ignore equivalent values with attribute order differences.
* data-sources/akamai_property_rules: comparisons in rule tree now properly ignore equivalent values with attribute order differences.
* provider: updated all error messages to better identify issues and actions required by user
* provider: fixed crash due to unexpected data types from unexpected API responses
* provider: fixed crash due to unexpected data types in Terraform files
* provider: errors now get reported using Terraform diagnostics allowing much more detail to be passed to user when an error occurs.

## 0.9.1 (Sept 02, 2020)
#### BREAKING CHANGES:
* [IMPORTANT] Dropped support for TF clients <= 0.11. Provider now built using Terraform sdk v2 library. Terraform dropped 0.11 client support as part of this update.  This change will make many new enhancements possible. ([See: Terraform v2 sdk](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html))
* resources/akamai_group: contract field (previously optional) now required to ensure contract and group agreement.

#### NOTES:
* [CHANGE] Individual edgerc file sections for different Akamai APIs (i.e., `property_section`, `dns_section`) has been deprecated in favor a common `config_section` used in conjuction with provider aliases ([See: Multiple Provider Configurations](https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-configurations))

#### KNOWN BUGS:
* provider: provider configuration validation requires an edgerc file configured and present even when one should not be needed.
* provider: support for configuring the provider via an inline provider block (`property`, `dns`, or `gtm`) no longer works.  Users should use edgerc file or Terraform environment args to configure instead.

#### BUG FIXES:
* [FIX] datasource akamai_group will no longer panic when contract not provided
* [ADD] Project re-organized to prepare for additional APIs to be included
* Fixed build job to compile sub-modules. Code is identical to 0.9.0 release

## 0.9.0 (August 26, 2020)
* [IMPORTANT] This build did not compile all modules properly so use 0.9.1 above instead.

## 0.8.2 (August 13, 2020)
* Initial release via the Terraform Registry. Otherwise identical to 0.8.1 release

## 0.8.1 (July 30, 2020)
* [FIX] Activation is executed, even without changes #139 (`akamai-property-activation`) ([#139](https://github.com/akamai/terraform-provider-template/issues/139))
* [FIX] Cannot find group when there are groups with the same name under multiple contract. #168 (`akamai-property-group`) ([#168](https://github.com/akamai/terraform-provider-template/issues/168))

## 0.8.0 (July 13, 2020)
* [FIX] Corrected Error 401 [Signature does not match] during  new primary zone creation (`akamai-dns`) ([#163](https://github.com/terraform-providers/terraform-provider-template/issues/163))
* [ADD] Updated Getting Started Primary Zone creation description. Added FAQ for Primary zone (`akamai-dns`)
* [FIX] SRV record priority value of 0 not allowed (`akamai-dns`) ([#165](https://github.com/terraform-providers/terraform-provider-template/issues/165))
* [ADD] Initial support for correlation ID in logging (`akamai-property`)

## 0.7.2 (June 11, 2020)
* [FIX] Corrected AAAA record handling of short and long IPv6 notation (`akamai-dns`)
## 0.7.1 (June 01, 2020)
* [FIX] Error after upgrading to 0.7.0 regarding MX records (`akamai-dns`) ([#154](https://github.com/terraform-providers/terraform-provider-template/issues/154))
* [FIX]Error 422 on SOA Record Apply After Creating a Primary Zone (`akamai-dns`) ([#155](https://github.com/terraform-providers/terraform-provider-template/issues/155))
## 0.7.0 (May 21, 2020)
* [ADD] User Agent support for Terraform version and provider version and SDK update
* [FIX] Bugs in Zone Create and Exists (`akamai_dns`) ([#151](https://github.com/terraform-providers/terraform-provider-template/issues/151))
## 0.6.0 (May 18, 2020)
* [ADD] Support the creation of DNS records of type AKAMAICDN (`akamai_dns`) ([#53](https://github.com/terraform-providers/terraform-provider-template/issues/53))
* [ADD] Support akamai_dns_record Import (`akamai_dns`) ([#69](https://github.com/terraform-providers/terraform-provider-template/issues/69))
* [FIX] Cannot remove a backup_cname from GTM property (`akamai_gtm`) ([#124](https://github.com/terraform-providers/terraform-provider-template/issues/124))
* [ADD] DNS Alias Zone Support (`akamai_dns`) ([#125](https://github.com/terraform-providers/terraform-provider-template/issues/125))
* [ADD] DNS TSIG Key support (`akamai_dns`) ([#126](https://github.com/terraform-providers/terraform-provider-template/issues/126))
* [ADD] DNS SOA, AKAMAITLC Record Support (`akamai_dns`) ([#127](https://github.com/terraform-providers/terraform-provider-template/issues/127))
* [FIX] Inverted Parameters - DNS Record Type NAPTR (`akamai_dns`) ([#130](https://github.com/terraform-providers/terraform-provider-template/issues/130))
* [FIX] Inverted Parameters - DNS Record Type NSEC3 (`akamai_dns`) ([#131](https://github.com/terraform-providers/terraform-provider-template/issues/131))
* [FIX] Inverted Parameters - DNS Record Type NSEC3PARAM (`akamai_dns`) ([#132](https://github.com/terraform-providers/terraform-provider-template/issues/132))
* [FIX] Inverted Parameters - DNS Record Type RRSIG (`akamai_dns`) ([#133](https://github.com/terraform-providers/terraform-provider-template/issues/133))
* [FIX] Inverted Parameters - DNS Record Type DS (`akamai_dns`) ([#134](https://github.com/terraform-providers/terraform-provider-template/issues/134))
* [ADD] DNS CAA, TLSA, CERT Record Support (`akamai_dns`) ([#148](https://github.com/terraform-providers/terraform-provider-template/issues/148))

## 0.5.0 (March 06, 2020)
* [FIX] Release edgehostnames and products caching edge library v0.9.10 (`akamai_property`)

## 0.4.0 (March 03, 2020)
* [FIX] Release contract group and cpcode caching edge library v0.9.9 (`akamai_property`) 

## 0.3.0 (March 02, 2020)
* [FIX] Provider produced inconsistent final plan #88 add contract group and cpcode caching edge library v0.9.9 (`akamai_property`) ([#88](https://github.com/terraform-providers/terraform-provider-template/issues/88))

## 0.2.0 (February 28, 2020)
* [FIX] Bug - Origin values customhostheader #93 (`akamai_property`) ([#93](https://github.com/terraform-providers/terraform-provider-template/issues/93))
* [FIX] akamai 0.1.5 - err: rpc error: code = Unavailable desc = transport is closing #87 (`akamai_property`) ([#87](https://github.com/terraform-providers/terraform-provider-template/issues/87))
* [FIX] Errors in documentation: akamai_contract and akamai_cp_code #52 (`akamai_property`) ([#52](https://github.com/terraform-providers/terraform-provider-template/issues/52))
* [FIX] Provider produced inconsistent final plan #88 (`akamai_property`) ([#88](https://github.com/terraform-providers/terraform-provider-template/issues/88))
* [FIX] akamai_property_activation creation crashing with Error: rpc error: code = Unavailable desc = transport is closing #102 (`akamai_property`) ([#102](https://github.com/terraform-providers/terraform-provider-template/issues/102))
* [ADD] Add Support for GTM domains and contained elements (domain, datacenter, property, resource, cidrmap, geographicmap, asmap)

## 0.1.5 (January 06, 2020)

* [FIX] Criteria is always end up using must satisfy "all" (`akamai_property`) ([#81](https://github.com/terraform-providers/terraform-provider-template/issues/81))
* [FIX] Provider produced inconsistent final plan (`akamai_property_variables`) ([#82](https://github.com/terraform-providers/terraform-provider-template/issues/82))
* [FIX] Cannot create multiple types of records with the same name (`akamai_dns_record`) ([#11](https://github.com/terraform-providers/terraform-provider-template/issues/11))
* [FIX] akamai_property_activation resource - changing network field causes deactivation of version in staging (`akamai_property_activation`) ([#51](https://github.com/terraform-providers/terraform-provider-template/issues/51))
* [FIX] Multiple MX records creation issue (`akamai_dns_record`) ([#57](https://github.com/terraform-providers/terraform-provider-template/issues/57))

## 0.1.4 (December 06, 2019)
* [FIX] Add support for update of rules state (`akamai_property`) ([#66](https://github.com/terraform-providers/terraform-provider-template/issues/66))
* [FIX] Add support for masters being optional (`akamai_dns_zone`) ([#61](https://github.com/terraform-providers/terraform-provider-template/issues/61))
* [FIX] Create edge hostname 400 error Bad Request Request parameter Slot Number (`akamai_property`) ([#56](https://github.com/terraform-providers/terraform-provider-template/issues/56))
* [FIX] TXT record - State update failure due to sha verification issue (`akamai_dns_zone`) ([#58](https://github.com/terraform-providers/terraform-provider-template/issues/58))
## 0.1.3 (August 12, 2019)

* [FIX] Correct ordering of values for `SRV` records (`akamai_dns_record`) ([#17](https://github.com/terraform-providers/terraform-provider-template/issues/17))
* [FIX] IPV4-only hostnames no longer fail (`akamai_edge_hostname`) ([#21](https://github.com/terraform-providers/terraform-provider-template/issues/21))
* [FIX] Don't try to deactive any version but the current one (`akamai_property_activation`) ([#21](https://github.com/terraform-providers/terraform-provider-template/issues/21))
* [FIX] Fix crash in DNS record validation ([#27](https://github.com/terraform-providers/terraform-provider-template/issues/27))
* [FIX] SiteShield behavior translated correctly to JSON ([#10](https://github.com/terraform-providers/terraform-provider-template/issues/10)] [[#40](https://github.com/terraform-providers/terraform-provider-template/issues/40))
* [FIX] Property rules correctly update (all rules now removed correctly) ([#30](https://github.com/terraform-providers/terraform-provider-template/issues/30))
* [FIX] Property Hostnames correctly update (all hostnames are now removed correctly) ([#44](https://github.com/terraform-providers/terraform-provider-template/issues/44))
* [FIX] Property activation was using the activation ID to fetch the property ([#35](https://github.com/terraform-providers/terraform-provider-template/issues/35))
* [FIX] Ensure property supports `is_secure` for Enhanced TLS ([#42](https://github.com/terraform-providers/terraform-provider-template/issues/42))
* [FIX] Multiple fixes to provider configuration for auth configuration. ([#46](https://github.com/terraform-providers/terraform-provider-template/issues/46))
* [FIX] Ensure the latest version is activated when no `akamai_property_activation.version` is set ([#45](https://github.com/terraform-providers/terraform-provider-template/issues/45))
* [FIX] Multiple records (e.g. using `count`) should now be created correctly ([#11](https://github.com/terraform-providers/terraform-provider-template/issues/11))
* [CHANGE] `akamai_property_rules` has been changed to a data source to ensure dependant resources update correctly, the existing resource now emits an error in all operations ([#47](https://github.com/terraform-providers/terraform-provider-template/issues/47))
* [ADD] Make zone type (primary or secondary) case-insensitive ([#29](https://github.com/terraform-providers/terraform-provider-template/issues/29))

## 0.1.2 (July 26, 2019)

* [FIX] Fixed handling of CPCode behavior in rules.json
* [FIX] Fixed hostname complexity, now a simple `{"public.host" = "edge.host"}` map
* [FIX] Fixed accidental deactivations
* [ADD] Added explicit property and dns credential blocks to provider config
* [ADD] Added better validation to `akamai_dns_record`

## 0.1.1 (July 09, 2019)

* [FIX] Bug fixes

## 0.1.0 (June 19, 2019)

* Initial release
