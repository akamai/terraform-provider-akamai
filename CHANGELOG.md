## 0.1.3 (Unreleased)

* [FIX] Correct ordering of values for `SRV` records (`akamai_dns_record`) [GH-17]
* [FIX] IPV4-only hostnames no longer fail (`akamai_edge_hostname`) [GH-21]
* [FIX] Don't try to deactive any version but the current one (`akamai_property_activation`) [GH-21]
* [FIX] Fix crash in DNS record validation [GH-27]
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
