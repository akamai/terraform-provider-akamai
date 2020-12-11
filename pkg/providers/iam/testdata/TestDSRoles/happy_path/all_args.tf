data "akamai_iam_roles" "test" {
  group_id       = 300
  get_actions    = true
  get_users      = true
  ignore_context = true
}
