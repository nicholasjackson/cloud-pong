resource "azurerm_resource_group" "pong" {
  name     = var.project
  location = var.region
}

module "aks" {
  source = "./aks"

  project = var.project
  resource_group = azurerm_resource_group.pong.name
  location = azurerm_resource_group.pong.location

  client_nodes = 3

  # Azure client id and secret to allow K8s to create loadbalancers
  client_id = var.client_id
  client_secret = var.client_secret
}

module "vms" {
  source = "./vms"

  project = var.project
  resource_group = azurerm_resource_group.pong.name
  location = azurerm_resource_group.pong.location

  consul_primary_addr = module.aks.consul_public_ip
}
