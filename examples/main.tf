terraform {
  required_providers {
    payments = {
      version = "0.1.0"
      source  = "jubril.xyz/custom/payments"
    }
  }
}


variable "api_token" {
  type = string
}

variable "email" {
  type = string

}
variable "account_number" {
    type = string
}



provider "payments" {
  api_token      = var.api_token
  email          = var.email
  account_number = var.account_number
}

//create variables 


data "payments_banks" "banks" {}

output "banks" {
  value = data.payments_banks.banks
}

