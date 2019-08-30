resource "azurerm_public_ip" "consul_gateway" {
  name                = "consul_gateway_ip"
  location            = var.location
  resource_group_name = var.resource_group
  allocation_method   = "Static"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_network_interface" "consul_gateway" {
  name                = "${var.project}-consul-gateway"
  location            = var.location
  resource_group_name = var.resource_group

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.internal.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id = azurerm_public_ip.consul_gateway.id
  }
}

data "template_file" "consul_gateway" {
  template = file("${path.module}/templates/consul-gateway.tpl")
  vars = {
    consul_cluster_addr = azurerm_network_interface.consul_server.private_ip_address
    gateway_addr = azurerm_public_ip.consul_gateway.ip_address
    advertise_addr = azurerm_network_interface.consul_gateway.private_ip_address
  }
}

resource "azurerm_virtual_machine" "consul_gateway" {
  name                  = "${var.project}-consul-gateway"
  location            = var.location
  resource_group_name = var.resource_group
  network_interface_ids = ["${azurerm_network_interface.consul_gateway.id}"]
  vm_size               = "Standard_DS1_v2"

  delete_os_disk_on_termination = true


  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk2"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "consul-gateway"
    admin_username = "ubuntu"
    admin_password = "Password1234!"
    custom_data = data.template_file.consul_gateway.rendered 
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      key_data = tls_private_key.pong.public_key_openssh
      path = "/home/ubuntu/.ssh/authorized_keys"
    }
  }

  tags = {
    environment = "staging"
  }
}
