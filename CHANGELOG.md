## 0.1.3 (Unreleased)

* [FIX] Correct ordering of values for `SRV` records (`akamai_dns_record`) [GH-17]
* [FIX] IPV4-only hostnames no longer fail (`akamai_edge_hostname`) [GH-21]
* [FIX] Don't try to deactive any version but the current one (`akamai_property_activation`) [GH-21]
* [FIX] Fix crash in DNS record validation [GH-27]
* [FIX] SiteShield behavior translated correctly to JSON [GH-10] [GH-40]
* [FIX] Property rules correctly update (all rules now removed correctly) [GH-30]
* [FIX] Property Hostnames correctly update (all hostnames are now removed correctly) [GH-44]
* [FIX] Property activation was using the activation ID to fetch the property [GH-35]
* [FIX] Ensure property supports `is_secure` for Enhanced TLS [GH-42]
* [FIX] Multiple fixes to provider configuration for auth configuration. [GH-46]
* [FIX] Ensure the latest version is activated when no `akamai_property_activation.version` is set [GH-45]
* [FIX] Multiple records (e.g. using `count`) should now be created correctly [GH-11]
* [CHANGE] `akamai_property_rules` has been changed to a data source to ensure dependant resources update correctly, the existing resource now emits an error in all operations [GH-47]
* [ADD] Make zone type (primary or secondary) case-insensitive [GH-29]

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
