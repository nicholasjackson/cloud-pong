output "k8s_config" {
  value = module.aks.k8s_config
}

output "aks_consul_addr" {
  value = module.aks.consul_public_ip
}

output "aks_pong_addr" {
  value = module.aks.pong_public_ip
}

output "vms_consul_server_addr" {
  value = module.vms.consul_server_addr
}

output "vms_consul_gateway_addr" {
  value = module.vms.consul_gateway_addr
}

output "vms_pong_addr" {
  value = module.vms.pong_addr
}

output "vms_private_key" { 
  value = module.vms.private_key
}
