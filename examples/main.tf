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

variable "bank_code" {
    type = string
}

variable "recipient" {
  type = string
}




provider "payments" {
  api_token      = var.api_token
  email          = var.email
  account_number = var.account_number
}

//create variables 


# data "payments_banks" "banks" {}

# output "banks" {
#   value = data.payments_banks.banks
# }

resource "payments_bank_transfer" "tf" {
  amount      = "2000"
  account_number  = var.recipient
  bank_code   = var.bank_code
}

output "status" {
  value = payments_bank_transfer.tf.message
}

output "tr_ref" {
  value = payments_bank_transfer.tf.request_reference
}
  



