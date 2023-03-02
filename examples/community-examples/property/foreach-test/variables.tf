variable "env" {
  type    = string
  default = "staging"
}

variable "customers" {
  type = map(object({
    username = string
    password = string
  }))
}
