terraform {
  required_version = ">= 0.12"
}

module "property" {
  source  = "../../modules/property"
  network = "staging"
  env     = "dev-envtest"
}
