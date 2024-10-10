provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_user" "test" {
  first_name = "John"
  last_name  = "Smith"
  email      = "jsmith@example.com"
  country    = "country"
  enable_tfa = false

  contact_type       = "contact type"
  job_title          = "job title"
  time_zone          = "timezone"
  secondary_email    = "secondary.email@example.com"
  mobile_phone       = "+49 98765 4321"
  address            = "123 B Street"
  city               = "B-Town"
  state              = "state"
  zip_code           = "zip"
  preferred_language = "language"
  session_timeout    = 2
  password           = ""
  auth_grants_json   = "[{\"groupId\":0,\"groupName\":\"group\",\"roleDescription\":\"\",\"roleName\":\"\"}]"
}
