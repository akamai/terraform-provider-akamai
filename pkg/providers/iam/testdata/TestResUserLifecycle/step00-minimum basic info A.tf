resource "akamai_iam_user" "test" {
  first_name     = "first name A"
  last_name      = "last name A"
  email          = "email-A@akamai.net"
  country        = "country A"
  phone          = "phone A"
  enable_tfa     = true
  send_otp_email = true
}
