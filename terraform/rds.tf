resource "aws_db_subnet_group" "rds" {
  name       = "${var.cluster-name}-db-subnet"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_security_group" "rds" {
  name        = "${var.cluster-name}-rds-sg"
  vpc_id      = module.vpc.vpc_id
  description = "Alowing PostgresQL traffic  from EKS"

  ingress {
    description = "Allowing traffic within vpc"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }
}
module "db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~>6.0"

  identifier = "${var.cluster-name}-postgres"

  engine               = "postgres"
  engine_version       = "15"
  family               = "postgres15"
  major_engine_version = "15"

  instance_class    = "db.t3.micro"
  allocated_storage = 15

  db_name  = "pulse"
  port     = 5432
  username = "pulseadmin"

  manage_master_user_password = "true"
  db_subnet_group_name        = aws_db_subnet_group.rds.name
  subnet_ids                  = module.vpc.private_subnets
  vpc_security_group_ids      = [aws_security_group.rds.id]

  skip_final_snapshot = "true"
}
