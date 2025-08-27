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

resource "akamai_apr_user_risk_response_strategy" "test" {
  config_id = 96033
  user_risk_response_strategy = jsonencode(
    {
      "traffic" : {
        "standard" : {
          "cautious" : {
            "threshold" : 0
          },
          "strict" : {
            "threshold" : 51
          },
          "aggressive" : {
            "threshold" : 66
          }
        },
        "inline" : {
          "cautious" : {
            "threshold" : 0
          },
          "strict" : {
            "threshold" : 51
          },
          "aggressive" : {
            "threshold" : 76
          }
        },
        "nativeSdkIos" : {
          "cautious" : {
            "threshold" : 0
          },
          "strict" : {
            "threshold" : 51
          },
          "aggressive" : {
            "threshold" : 76
          }
        },
        "nativeSdkAndroid" : {
          "cautious" : {
            "threshold" : 0
          },
          "strict" : {
            "threshold" : 51
          },
          "aggressive" : {
            "threshold" : 76
          }
        }
      }
    }
  )
}
