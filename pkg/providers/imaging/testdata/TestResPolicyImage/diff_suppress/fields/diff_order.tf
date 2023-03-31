provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_imaging_policy_image" "policy" {
  policy_id    = "test_policy"
  contract_id  = "test_contract"
  policyset_id = "test_policy_set"
  json = jsonencode({
    "hosts" : [
      "test3",
      "test1",
      "test2"
    ],
    "breakpoints" : {
      "widths" : [
        5000,
        2048,
        1024,
        640,
        320
      ]
    },
    "output" : {
      "allowedFormats" : [
        "gif",
        "webp",
        "avif",
        "jpeg",
        "png"
      ],
      "forcedFormats" : [
        "gif",
        "webp",
        "avif",
        "jpeg",
        "png"
      ],
      "perceptualQuality" : "mediumHigh",
    },
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