
resource "random_string" "storage_name" {
  length  = 12
  lower   = true
  numeric = true
  special = false
  upper   = false
}

resource "azurerm_storage_account" "storage" {
  name                     = coalesce(var.storage_account_name, "st${random_string.storage_name.result}")
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = var.account_tier
  account_replication_type = var.account_replication_type
  tags                     = var.tags

  blob_properties {
    versioning_enabled = var.enable_versioning
  }

  network_rules {
    default_action = var.default_network_rule
    ip_rules       = var.allowed_ips
    bypass         = ["AzureServices"]
  }
}

resource "azurerm_storage_container" "containers" {
  for_each = toset(var.container_names)

  name                  = each.value
  storage_account_name  = azurerm_storage_account.storage.name
  container_access_type = "private"
}
