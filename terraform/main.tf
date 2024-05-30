module "ecr" {
    source = "./modules/ecr"
    teamname = var.teamname
}

data "aws_subnets" "subnets" {
  filter {
    name   = "availability-zone"
    values = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]  # Replace with your availability zones
  }
}

module "db" {
    source = "./modules/db"
    teamname = var.teamname
    execution_role_arn = aws_iam_role.task_execution_role.arn
    task_role_arn = aws_iam_role.task_role.arn
    cluster_arn = aws_ecs_cluster.projectsprint.arn
    subnets = data.aws_subnets.subnets.ids
    db_security_group   = var.db_security_group
    db_name = var.db_name
    db_user = var.db_user
    db_password = var.db_password
}