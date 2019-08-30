output "private_key" {
  value = tls_private_key.pong.private_key_pem
}

output "consul_server_addr" {
  value = azurerm_public_ip.consul_server.ip_address
}

output "consul_gateway_addr" {
  value = azurerm_public_ip.consul_gateway.ip_address
}

output "pong_addr" {
  value = azurerm_public_ip.pong.ip_address
}
