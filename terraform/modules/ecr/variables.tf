variable "ecr_name" {
  description = "The name of the team for the ECR repository naming convention"
  type        = string
}

variable "image_mutability" {
  description = "Provide image mutability"
  type = string
  default = "IMMUTABLE"
}