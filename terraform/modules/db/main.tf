data "aws_subnets" "subnets" {
  filter {
    name   = "availability-zone"
    values = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]  # Replace with your availability zones
  }
}

resource "aws_ecs_task_definition" "db_task" {
    family = "${var.teamname}-db-task-definition"
    requires_compatibilities = ["FARGATE"]
    cpu = "256" # 0.25 vCPU
    memory = "512" # 0.5 GB
    network_mode = "awsvpc"
    execution_role_arn = var.execution_role_arn
    task_role_arn = var.task_role_arn

  container_definitions = jsonencode([
    {
      name      = "${var.teamname}-db"
      image     = "postgres:16.3-alpine3.19"
      essential = true
      portMappings = [
        {
          containerPort = 5432
          hostPort      = 5432
        }
      ]
      environment = [
        { name = "DB_USERNAME", value = var.db_user },
        { name = "DB_PASSWORD", value = var.db_password }
      ]
    }
  ])
}

resource "aws_ecs_service" "db_service" {
    name = "${var.teamname}-db-service"
    cluster = var.cluster_arn
    task_definition = aws_ecs_task_definition.db_task.arn
    launch_type = "FARGATE"
    desired_count = 1

    network_configuration {
        subnets         = var.subnets
        security_groups = [var.db_security_group]
        assign_public_ip = true
    }
}