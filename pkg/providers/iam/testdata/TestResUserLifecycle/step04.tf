# all basic info B with notifications B
resource "akamai_iam_user" "test" {
  first_name     = "first name B"
  last_name      = "last name B"
  email          = "email@akamai.net"
  country        = "country B"
  phone          = "(111) 111-1111"
  enable_tfa     = false
  send_otp_email = false

  contact_type       = "contact type B"
  job_title          = "job title B"
  time_zone          = "Timezone B"
  secondary_email    = "secondary-email-B@akamai.net"
  mobile_phone       = "(111) 111-1111"
  address            = "123 B Street"
  city               = "B-Town"
  state              = "state B"
  zip_code           = "zip B"
  preferred_language = "language B"
  session_timeout    = 2

  enable_notifications          = false
  subscribe_new_users           = false
  subscribe_password_expiration = false
  subscribe_product_issues      = []
  subscribe_product_upgrades    = []

  auth_grants_json = "[{\"groupId\":0,\"groupName\":\"B\",\"isBlocked\":false,\"roleDescription\":\"\",\"roleName\":\"\"}]"
}
