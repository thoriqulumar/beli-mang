variable "teamname" {
  description = "The name of the team for the ECR repository naming convention"
  type        = string
}

variable "db_security_group" {
  description = "Security group ID for the database service"
  type        = string
}

variable "service_security_group" {
  description = "Security group ID for the service"
  type        = string
}

variable "db_name" {
  description = "Name of the Postgres database"
  type        = string
}

variable "db_user" {
  description = "Username for the Postgres database"
  type        = string
}

variable "db_password" {
  description = "Password for the Postgres database"
  type        = string
}

variable "container_name" {
  description = "Name of the container for the service"
  type        = string
}