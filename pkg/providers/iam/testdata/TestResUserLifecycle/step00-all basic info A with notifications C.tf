resource "akamai_iam_user" "test" {
  first_name     = "first name A"
  last_name      = "last name A"
  email          = "email-A@akamai.net"
  country        = "country A"
  phone          = "phone A"
  enable_tfa     = true
  send_otp_email = true

  contact_type       = "contact type A"
  user_name          = "user name A"
  job_title          = "job title A"
  time_zone          = "Timezone A"
  secondary_email    = "secondary-email-A@akamai.net"
  mobile_phone       = "mobile phone A"
  address            = "123 A Street"
  city               = "A-Town"
  state              = "state A"
  zip_code           = "zip A"
  preferred_language = "language A"
  session_timeout    = 1

  enable_notifications          = true
  subscribe_new_users           = true
  subscribe_password_expiration = true
  subscribe_product_issues      = ["issues product"]
  subscribe_product_upgrades    = ["upgrades product"]
}
