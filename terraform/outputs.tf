output "k8s_config" {
  value = module.aks.k8s_config
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
