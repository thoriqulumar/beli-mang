module "ecr" {
    source = "./modules/ecr"
    ecr_name = var.ecr_name
}