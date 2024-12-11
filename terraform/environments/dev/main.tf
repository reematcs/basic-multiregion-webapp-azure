# Azure provider configuration is moved to providers.tf

# Configure resource group
resource "azurerm_resource_group" "rg" {
  name     = "multiregion-webapp-rg"
  location = "eastus"
}