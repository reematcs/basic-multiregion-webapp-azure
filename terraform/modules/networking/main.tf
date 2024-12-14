resource "random_string" "vnet_name" {
  length      = 8
  special     = false
  upper       = false
  min_lower   = 1
  min_numeric = 1
}

locals {
  # Add prefix to ensure name starts with letter (Azure requirement)
  vnet_name = var.vnet_name != null ? var.vnet_name : "vnet-${random_string.vnet_name.result}"

  # Generate NSG names based on subnet names
  nsg_names = {
    for subnet in var.subnet_configurations :
    subnet.name => "${subnet.name}-nsg-${random_string.vnet_name.result}"
    if subnet.create_network_security_group
  }
}

resource "azurerm_virtual_network" "vnet" {
  name                = local.vnet_name
  location            = var.location
  resource_group_name = var.resource_group_name
  address_space       = var.address_space
  dns_servers         = var.dns_servers

  tags = var.tags
}

resource "azurerm_subnet" "subnet" {
  for_each = { for subnet in var.subnet_configurations : subnet.name => subnet }

  name                 = each.value.name
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = [each.value.address_prefix]
  service_endpoints    = each.value.service_endpoints

  dynamic "delegation" {
    for_each = each.value.delegation != null ? [each.value.delegation] : []
    content {
      name = delegation.value.name
      service_delegation {
        name    = delegation.value.service_delegation_name
        actions = delegation.value.service_delegation_actions
      }
    }
  }
}
resource "azurerm_network_security_group" "nsg" {
  for_each = { for subnet in var.subnet_configurations : subnet.name => subnet if subnet.create_network_security_group }

  name                = "${local.vnet_name}-${each.value.name}-nsg" # Use local.vnet_name
  location            = var.location
  resource_group_name = var.resource_group_name

  dynamic "security_rule" {
    for_each = each.value.nsg_rules != null ? each.value.nsg_rules : []
    content {
      name                       = security_rule.value.name
      priority                   = security_rule.value.priority
      direction                  = security_rule.value.direction
      access                     = security_rule.value.access
      protocol                   = security_rule.value.protocol
      source_port_range          = security_rule.value.source_port_range
      destination_port_range     = security_rule.value.destination_port_range
      source_address_prefix      = security_rule.value.source_address_prefix
      destination_address_prefix = security_rule.value.destination_address_prefix
    }
  }

  tags = var.tags
}

resource "azurerm_subnet_network_security_group_association" "nsg_association" {
  for_each = { for subnet in var.subnet_configurations : subnet.name => subnet if subnet.create_network_security_group }

  subnet_id                 = azurerm_subnet.subnet[each.key].id
  network_security_group_id = azurerm_network_security_group.nsg[each.key].id
}