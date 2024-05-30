### init
terraform init

### define plan
terraform plan -var="teamname=teamname" -var="db_security_group=db_security_group" -var="service_security_group=service_security_group" -var="db_name=db_name" -var="db_user= db_user" -var="db_password=db_password" -var="container_name=container_name"

### apply plan
terraform apply -var="teamname=teamname" -var="db_security_group=db_security_group" -var="service_security_group=service_security_group" -var="db_name=db_name" -var="db_user= db_user" -var="db_password=db_password" -var="container_name=container_name"