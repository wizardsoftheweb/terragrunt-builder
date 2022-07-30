variable "image_id" {
  type        = string
  description = "The id of the machine image (AMI) to use for the server."
  sensitive   = false
}

variable "other_id" {
  type        = string
  description = "The id of the machine image (AMI) to use for the server."
  sensitive   = true
}