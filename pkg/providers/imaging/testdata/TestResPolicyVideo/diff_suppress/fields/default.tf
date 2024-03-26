provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_imaging_policy_video" "policy" {
  policy_id    = "test_policy"
  contract_id  = "test_contract"
  policyset_id = "test_policy_set"
  json = jsonencode({
    "breakpoints" : {
      "widths" : [
        320,
        640,
        1024,
        2048,
        5000
      ]
    },
    "hosts" : [
      "test1",
      "test2",
      "test3"
    ],
    "variables" : [
      {
        "defaultValue" : "test1"
        "name" : "var1",
        "type" : "string",
      },
      {
        "defaultValue" : "test2"
        "name" : "var2",
        "type" : "string",
      },
      {
        "defaultValue" : "test3"
        "name" : "var3",
        "type" : "string",
      }
    ]
  })
}