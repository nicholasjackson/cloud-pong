resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.project
  location            = var.location
  resource_group_name = var.resource_group
  dns_prefix          = var.project

  agent_pool_profile {
    name            = "default"
    count           = var.client_nodes
    vm_size         = "Standard_D1_v2"
    os_type         = "Linux"
    os_disk_size_gb = 30
  }

  service_principal {
    client_id     = var.client_id
    client_secret = var.client_secret
  }

  tags = {
    Environment = "Production"
  }
}
