# minimum basic info B
resource "akamai_iam_user" "test" {
  first_name     = "first name B"
  last_name      = "last name B"
  email          = "email-B@akamai.net"
  country        = "country B"
  phone          = "phone B"
  enable_tfa     = false
  send_otp_email = false

  auth_grants_json = "[{\"groupId\":0,\"groupName\":\"B\",\"isBlocked\":false,\"roleDescription\":\"\",\"roleName\":\"\"}]"
}
