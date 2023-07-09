terraform {
  required_providers {
    payments = {
      source  = "jubril/custom/payments"
    }
  }
}

provider "payments" {
  api_token = "VprYNiCxdaEG9cwX3vqK"
  email = "hello@jubril.xyz"
  account_number = "3000763375"
}