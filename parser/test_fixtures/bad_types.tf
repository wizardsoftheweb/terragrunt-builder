variable "fails" {
  type = string
  default = {
    bad = "value"
  }
}

output "fails" {
  type = string
  value = {
    bad = "value"
  }
}
