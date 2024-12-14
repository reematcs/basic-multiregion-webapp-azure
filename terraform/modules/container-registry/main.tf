# modules/container-registry/main.tf
resource "random_string" "acr_name" {
  length      = 8
  special     = false
  upper       = false
  min_lower   = 1
  min_numeric = 1
}

locals {
  container_registry_name = var.container_registry_name != null ? var.container_registry_name : "acr${random_string.acr_name.result}"
}

resource "azurerm_container_registry" "acr" {
  name                = local.container_registry_name
  resource_group_name = var.resource_group_name
  location            = var.location
  sku                 = var.sku
  admin_enabled       = var.admin_enabled

  dynamic "georeplications" {
    for_each = var.sku == "Premium" ? var.georeplication_locations : []
    content {
      location = georeplications.value
      tags     = var.tags
    }
  }

  dynamic "network_rule_set" {
    for_each = var.sku == "Premium" ? [1] : []
    content {
      default_action = var.network_rule_default_action
    }
  }

  tags = var.tags
}

