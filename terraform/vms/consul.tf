resource "azurerm_public_ip" "consul_server" {
  name                = "consul_public_ip"
  location            = var.location
  resource_group_name = var.resource_group
  allocation_method   = "Static"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_network_interface" "consul_server" {
  name                = "${var.project}-consul-server"
  location            = var.location
  resource_group_name = var.resource_group

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.internal.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id = azurerm_public_ip.consul_server.id
  }
}

data "template_file" "consul_server" {
  template = file("${path.module}/templates/consul-server.tpl")
  vars = {
    primary_cluster_addr = var.consul_primary_addr
    wan_addr = azurerm_public_ip.consul_server.ip_address
    advertise_addr = azurerm_network_interface.consul_server.private_ip_address
  }
}

resource "azurerm_virtual_machine" "consul_server" {
  name                  = "${var.project}-consul-server"
  location            = var.location
  resource_group_name = var.resource_group
  network_interface_ids = ["${azurerm_network_interface.consul_server.id}"]
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
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "consul-server"
    admin_username = "ubuntu"
    admin_password = "Password1234!"
    custom_data = data.template_file.consul_server.rendered 
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
