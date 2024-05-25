variable "ecr_name" {
  description = "The name of the team for the ECR repository naming convention"
  type        = string
}

variable "image_mutability" {
  description = "Provide image mutability"
  type = string
  default = "IMMUTABLE"
}

resource "aws_ecr_repository" "ecr" {
    name = "${var.ecr_name}-repository"
    image_tag_mutability = var.image_mutability
    image_scanning_configuration {
        scan_on_push = true
    }
    encryption_configuration {
        encryption_type = "AES256"
    }
}