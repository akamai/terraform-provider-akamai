variable "env" {
        default = "staging"
}

variable "username" {
	default = "test"
}

variable "password" {
	default = "test"
}

variable "customers" {
	type = map(object({
		username = string
		password = string
	}))
}
