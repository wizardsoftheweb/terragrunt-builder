variable "image_id" {
  type        = string
  description = "The id of the machine image (AMI) to use for the server."
  sensitive   = false
  default     = "test"
}

variable "other_id" {
  type        = string
  description = "The id of the machine image (AMI) to use for the server."
  sensitive   = true
}

output "random_output" {
  value = "1"
}
