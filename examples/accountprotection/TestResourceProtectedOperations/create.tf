terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.0.0"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_apr_protected_operations" "test" {
  config_id          = 96033
  security_policy_id = "UCON_161669"
  operation_id       = "3bcb59e4-f53d-40c4-915c-142ebd8f8c03"
  protected_operation = jsonencode(
    {
      "apiEndPointId" : 820204,
      "traffic" : {
        "standard" : {
          "overrideThresholds" : false,
          "cautious" : {
            "action" : "monitor"
          },
          "strict" : {
            "action" : "monitor"
          },
          "aggressive" : {
            "action" : "deny"
          }
        },
        "nativeSdkIos" : {
          "overrideThresholds" : false,
          "bypassPreSdkVersion" : false,
          "cautious" : {
            "action" : "monitor"
          },
          "strict" : {
            "action" : "monitor"
          },
          "aggressive" : {
            "action" : "monitor"
          }
        },
        "nativeSdkAndroid" : {
          "overrideThresholds" : false,
          "bypassPreSdkVersion" : false,
          "cautious" : {
            "action" : "monitor"
          },
          "strict" : {
            "action" : "monitor"
          },
          "aggressive" : {
            "action" : "monitor"
          }
        }
      },
      "telemetryTypeStates" : {
        "standard" : {
          "enabled" : true
        },
        "inline" : {
          "enabled" : false,
          "disabledAction" : "none"
        },
        "nativeSdk" : {
          "enabled" : true
        }
      }
    }
  )
}
