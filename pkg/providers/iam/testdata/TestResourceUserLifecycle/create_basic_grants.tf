provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_iam_user" "test" {
  first_name = "John"
  last_name  = "Smith"
  email      = "jsmith@example.com"
  country    = "country"
  phone      = "(111) 111-1111"
  enable_tfa = false

  auth_grants_json = jsonencode([
    {
      subGroups = [
        {
          groupId   = 2
          isBlocked = false
        },
        {
          groupId   = 1
          isBlocked = false
        }
      ]
    }
  ])
}
