resource "azurerm_virtual_network" "pong" {
  name                = "${var.project}-network"
  address_space       = ["10.0.0.0/16"]
  location            = var.location
  resource_group_name = var.resource_group
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name = var.resource_group
  virtual_network_name = azurerm_virtual_network.pong.name
  address_prefix       = "10.0.2.0/24"
}
