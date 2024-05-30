variable "teamname" {
  description = "The name of the team for naming conventions"
  type        = string
}

variable "execution_role_arn" {
  description = "ARN of the task execution role"
  type        = string
}

variable "task_role_arn" {
  description = "ARN of the task role"
  type        = string
}

variable "cluster_arn" {
  description = "ARN of the ECS cluster"
  type        = string
}

variable "subnets" {
  description = "List of subnet IDs for the ECS service"
  type        = list(string)
}

variable "db_security_group" {
  description = "Security group ID for the database service"
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