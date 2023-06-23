provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_imaging_policy_image" "policy" {
  policy_id    = "test_policy"
  contract_id  = "test_contract"
  policyset_id = "test_policy_set"
  json = jsonencode({
    "hosts" : [
      "test1",
      "test2",
      "test3"
    ],
    "breakpoints" : {
      "widths" : [
        320,
        640,
        1024,
        2048,
        5000
      ]
    },
    "output" : {
      "allowedFormats" : [
        "jpeg",
        "webp",
        "avif",
        "png",
        "gif"
      ],
      "forcedFormats" : [
        "jpeg",
        "webp",
        "avif",
        "png",
        "gif"
      ],
      "perceptualQuality" : "mediumHigh"
    },
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