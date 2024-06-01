module "ecr" {
    source = "./modules/ecr"
    teamname = var.teamname
}

data "aws_subnets" "subnets" {
  filter {
    name   = "availability-zone"
    values = ["ap-southeast-1a"]  # Replace with your availability zones
  }
}

data "aws_iam_role" "task_execution_role" {
  name = var.task_execution_role_name
}

data "aws_iam_role" "task_role" {
  name = var.task_role_name
}

data "aws_ecs_cluster" "ecs_cluster" {
  cluster_name = var.ecs_cluster_name
}

module "db" {
    source = "./modules/db"
    teamname = var.teamname
    execution_role_arn = data.aws_iam_role.task_execution_role.arn
    task_role_arn = data.aws_iam_role.task_role.arn
    cluster_arn = data.aws_ecs_cluster.ecs_cluster.arn
    subnets = data.aws_subnets.subnets.ids
    db_security_group   = var.db_security_group
    db_name = var.db_name
    db_user = var.db_user
    db_password = var.db_password
}