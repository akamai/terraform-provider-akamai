terraform {
  required_version = ">= 1.0"
}

module "property" {
  source  = "../../modules/property"
  network = "staging"
  env     = "preprod-envtest"
}
