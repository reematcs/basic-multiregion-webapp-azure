# Azure provider configuration is moved to providers.tf

# Configure resource group

resource "azurerm_resource_group" "rg" {
  name     = var.resource_group_name
  location = var.resource_group_location
}
module "key_vault" {
  source = "../../modules/key-vault"

  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  # Other module variables...
}

module "storage" {
  source = "../../modules/storage"

  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  # Other module variables...
}




# module "networking" {
#   source              = "../../modules/networking"
#   resource_group_name = azurerm_resource_group.rg.name

# }

# module "container_registry" {
#   source              = "../../modules/container-registry"
#   resource_group_name = azurerm_resource_group.rg.name

# }