output "hosts" {
  value       = module.provision.hosts
  description = "list of host machines used for the cluster"
}

output "mke_lb" {
  value       = "https://${module.provision.mke_lb}"
  description = "LB path for the MKE ingress"
}
