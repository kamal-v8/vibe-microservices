module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name    = var.cluster-name
  cluster_version = "1.30"

  cluster_endpoint_public_access = true

  vpc_id = module.vpc.vpc_id

  subnet_ids = module.vpc.private_subnets

  control_plane_subnet_ids = module.vpc.private_subnets

  eks_managed_node_groups = {
    pulse_nodes = {
      instance_types = ["c7i-flex.large"]

      min_size     = 1
      max_size     = 3
      desired_size = 2
    }
  }
  enable_cluster_creator_admin_permissions = true
}

