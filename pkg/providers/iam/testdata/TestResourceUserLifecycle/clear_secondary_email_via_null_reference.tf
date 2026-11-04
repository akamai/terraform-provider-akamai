provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

locals {
  # both entries are intentionally identical: the index stays unknown to Terraform until
  # apply, but the resulting secondary_email is nullified, which this test step relies on
  secondary_emails = [
    null,
    null
  ]
}

// min/max differ to satisfy the provider's validation; the chosen index is still
// unknown to Terraform at plan time
resource "random_integer" "secondary_email_choice" {
  min = 0
  max = 1
}

// This resource uses a random value to force an unknown secondary_email during the plan phase,
// simulating the pattern where secondary_email might be derived from an external source or computed value.
resource "akamai_iam_user" "test" {
  first_name = "John"
  last_name  = "Smith"
  email      = "jsmith@example.com"
  country    = "country"
  phone      = "(111) 111-1111"
  enable_tfa = false

  secondary_email  = local.secondary_emails[random_integer.secondary_email_choice.result]
  auth_grants_json = "[{}]"
}
