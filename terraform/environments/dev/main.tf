# Azure provider configuration is moved to providers.tf

# Configure resource group

resource "random_pet" "rg_name" {
  prefix    = var.resource_group_name_prefix
  length    = 2
  separator = "-"
}

resource "azurerm_resource_group" "rg" {
  name     = random_pet.rg_name.id
  location = var.resource_group_location
}

module "key_vault" {
  source = "../../modules/key-vault"

  resource_group_name = azurerm_resource_group.rg.name
  location           = azurerm_resource_group.rg.location
  vault_name         = var.key_vault_name
  sku_name           = "standard"
}

module "storage" {
  source = "../../modules/storage"

  resource_group_name      = azurerm_resource_group.rg.name
  location                = azurerm_resource_group.rg.location
  storage_account_name    = var.storage_account_name

}