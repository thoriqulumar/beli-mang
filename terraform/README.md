### prerequisite
install terraform cli by following this guide: https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli

### init
terraform init

### define plan
terraform plan -var="teamname=" -var="db_security_group=" -var="service_security_group=" -var="db_name=" -var="db_user=" -var="db_password=" -var="task_execution_role_name=" -var="task_role_name=" -var="ecs_cluster_name="

### apply plan
terraform apply -var="teamname=teamname" -var="db_security_group=db_security_group" -var="service_security_group=service_security_group" -var="db_name=db_name" -var="db_user= db_user" -var="db_password=db_password" -var="container_name=container_name"