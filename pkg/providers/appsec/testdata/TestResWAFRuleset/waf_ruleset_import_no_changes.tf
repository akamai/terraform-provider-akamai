provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_waf_ruleset" "import_no_changes" {
  config_id          = 63905
  security_policy_id = "p563_113699"

  # Configure individual WAF rules
  rules = [
    {
      rule_id     = 950003
      rule_action = "alert"
      condition_exception = jsonencode({
        advancedExceptions = {
          conditionOperator = "AND"
          conditions = [
            {
              type          = "pathMatch"
              paths         = ["/xxcatssssszzz*"]
              positiveMatch = true
            }
          ]
        }
      })
    }
  ]

  # Configure attack groups
  attack_groups = [
    {
      attack_group        = "CMD"
      attack_group_action = "alert"
      condition_exception = jsonencode({
        advancedExceptions = {
          conditionOperator = "AND"
          conditions = [
            {
              type          = "hostMatch"
              hosts         = ["cheppubakku.com"]
              positiveMatch = true
            },
            {
              type          = "ipMatch"
              ips           = ["21.1.1.1", "2.2.2.2", "3.3.3.3"]
              positiveMatch = true
            },
            {
              type          = "pathMatch"
              paths         = ["/test"]
              positiveMatch = true
            },
            {
              type          = "extensionMatch"
              extensions    = ["jpg", "kgf", "lsf"]
              positiveMatch = true
            },
            {
              type          = "extensionMatch"
              extensions    = ["abs"]
              positiveMatch = true
            },
            {
              type          = "filenameMatch"
              filenames     = ["abcd"]
              positiveMatch = true
            },
            {
              type          = "extensionMatch"
              extensions    = ["xyz", "ds"]
              positiveMatch = true
            },
            {
              type          = "ipMatch"
              ips           = ["22.2.2.3"]
              positiveMatch = true
              useHeaders    = true
            },
            {
              type          = "requestHeaderMatch"
              header        = "Max-Forwards"
              positiveMatch = true
              value         = "as"
              valueCase     = true
            },
            {
              type          = "pathMatch"
              paths         = ["/aaadasd"]
              positiveMatch = false
            },
            {
              type          = "pathMatch"
              paths         = ["/aaabcs", "/ssc"]
              positiveMatch = false
            }
          ]
        }
      })
    },
    {
      attack_group        = "WAT"
      attack_group_action = "alert"
      condition_exception = jsonencode({
        advancedExceptions = {
          conditionOperator = "AND"
          conditions = [
            {
              type          = "pathMatch"
              paths         = ["/xxicicibanksscatssssss*"]
              positiveMatch = true
            },
            {
              type          = "pathMatch"
              paths         = ["/xxdpogs"]
              positiveMatch = false
            }
          ]
        }
      })
    }
  ]
}
