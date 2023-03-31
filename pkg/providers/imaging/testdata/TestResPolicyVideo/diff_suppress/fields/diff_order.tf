provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_imaging_policy_video" "policy" {
  policy_id    = "test_policy"
  contract_id  = "test_contract"
  policyset_id = "test_policy_set"
  json = jsonencode({
    "breakpoints" : {
      "widths" : [
        5000,
        2048,
        1024,
        640,
        320
      ]
    },
    "hosts" : [
      "test3",
      "test1",
      "test2"
    ],
    "variables" : [
      {
        "defaultValue" : "test2"
        "name" : "var2",
        "type" : "string",
      },
      {
        "defaultValue" : "test3"
        "name" : "var3",
        "type" : "string",
      },
      {
        "defaultValue" : "test1"
        "name" : "var1",
        "type" : "string",
      }
    ]
  })
}