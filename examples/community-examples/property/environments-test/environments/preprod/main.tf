module "property" {
  source  = "../../modules/property"
  network = "staging"
  env     = "preprod-envtest"
}
