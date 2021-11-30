resource "akamai_iam_user" "test" {
  first_name = "first name A"
  last_name  = "last name A"
  email      = "email@akamai.net"
  country    = "country A"
  phone      = "(000) 000-0000"
  enable_tfa = true

  auth_grants_json = "[{\"groupId\":0,\"groupName\":\"A\",\"roleDescription\":\"\",\"roleName\":\"\"}]"
}
