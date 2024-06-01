variable "teamname" {
  description = "The name of the team for the ECR repository naming convention"
  type        = string
}

variable "task_execution_role_name" {
  description = "The name of the existing IAM role for task execution"
  type        = string
}

variable "task_role_name" {
  description = "The name of the existing IAM role for task"
  type        = string
}

variable "ecs_cluster_name" {
  description = "The name of the existing ECS cluster"
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