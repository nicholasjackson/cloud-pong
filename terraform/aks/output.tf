output "k8s_config" {
  value = azurerm_kubernetes_cluster.aks.kube_config_raw
}

output "consul_public_ip" {
  value = kubernetes_service.consul.load_balancer_ingress.0.ip
}

output "pong_public_ip" {
  value = kubernetes_service.pong_service.load_balancer_ingress.0.ip
}
