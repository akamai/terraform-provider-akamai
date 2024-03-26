provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_user" "test" {
  first_name = "John"
  last_name  = "Smith"
  email      = "jsmith@example.com"
  country    = "country"
  phone      = "(111) 111-1111"
  enable_tfa = false

  auth_grants_json = "[{}]"
}


resource "akamai_iam_user" "test2" {
  first_name = "John"
  last_name  = "Smith"
  email      = "jsmith@example.com"
  country    = "country"
  phone      = "(111) 111-1111"
  enable_tfa = false

  auth_grants_json = "[{\"groupId\": 123,\"groupName\":\"group\",\"roleDescription\":\"${akamai_iam_user.test.id}\",\"roleName\":\"\"}]"
}
