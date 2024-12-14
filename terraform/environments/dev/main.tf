locals {
  primary_tags = merge(var.common_tags, {
    region = var.primary_location
  })
  secondary_tags = merge(var.common_tags, {
    region = var.secondary_location
  })
}

resource "azurerm_resource_group" "rg" {
  name     = var.resource_group_name
  location = var.primary_location
  tags     = var.common_tags
}

module "key_vault" {
  source              = "../../modules/key-vault"
  resource_group_name = azurerm_resource_group.rg.name
  location           = azurerm_resource_group.rg.location
  tags               = var.common_tags
}

module "storage" {
  source                  = "../../modules/storage"
  resource_group_name     = azurerm_resource_group.rg.name
  location               = azurerm_resource_group.rg.location
  storage_account_name    = var.storage_account_name
  account_tier           = "Standard"
  account_replication_type = "LRS"
  container_names         = ["tfstate"]
  default_network_rule    = "Allow"
  tags                    = var.common_tags
}

module "networking_primary" {
  source              = "../../modules/networking"
  resource_group_name = azurerm_resource_group.rg.name
  location           = var.primary_location
  vnet_name          = "vnet-${var.primary_location}"
  address_space      = ["10.1.0.0/16"]
  tags               = local.primary_tags

  subnet_configurations = [
    {
      name                          = "subnet-app"
      address_prefix                = "10.1.1.0/24"
      create_network_security_group = true
      service_endpoints             = ["Microsoft.KeyVault", "Microsoft.Storage"]
      delegation = {
        name                       = "appservice"
        service_delegation_name    = "Microsoft.Web/serverFarms"
        service_delegation_actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
      }
    },
    {
      name                          = "subnet-endpoints"
      address_prefix                = "10.1.2.0/24"
      create_network_security_group = true
      service_endpoints             = ["Microsoft.KeyVault", "Microsoft.Storage"]
      delegation                    = null
    }
  ]
}

module "networking_secondary" {
  source              = "../../modules/networking"
  resource_group_name = azurerm_resource_group.rg.name
  location           = var.secondary_location
  vnet_name          = "vnet-${var.secondary_location}"
  address_space      = ["10.2.0.0/16"]
  tags               = local.secondary_tags

  subnet_configurations = [
    {
      name                          = "subnet-app"
      address_prefix                = "10.2.1.0/24"
      create_network_security_group = true
      service_endpoints             = ["Microsoft.KeyVault", "Microsoft.Storage"]
      delegation = {
        name                       = "appservice"
        service_delegation_name    = "Microsoft.Web/serverFarms"
        service_delegation_actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
      }
    },
    {
      name                          = "subnet-endpoints"
      address_prefix                = "10.2.2.0/24"
      create_network_security_group = true
      service_endpoints             = ["Microsoft.KeyVault", "Microsoft.Storage"]
      delegation                    = null
    }
  ]
}

module "container_registry" {
  source              = "../../modules/container-registry"
  resource_group_name = azurerm_resource_group.rg.name
  location           = azurerm_resource_group.rg.location
  sku                = "Standard"
  admin_enabled      = true
  tags               = var.common_tags
}

module "traffic_manager" {
  source               = "../../modules/traffic-manager"
  resource_group_name  = azurerm_resource_group.rg.name
  traffic_manager_name = "tm-${var.app_name}-${var.environment}"
  dns_name            = "${var.app_name}-${var.environment}-20241331"
  primary_web_app_id   = module.web_app_primary.web_app_id
  secondary_web_app_id = module.web_app_secondary.web_app_id
  tags                 = var.common_tags
}

module "web_app_primary" {
  source                    = "../../modules/web-app"
  resource_group_name       = azurerm_resource_group.rg.name
  location                  = var.primary_location
  service_plan_name         = "asp-${var.app_name}-${var.primary_location}-${var.environment}"
  web_app_name             = "app-${var.app_name}-${var.primary_location}-${var.environment}"
  subnet_id                = module.networking_primary.subnet_ids["subnet-app"]
  private_endpoint_subnet_id = module.networking_primary.subnet_ids["subnet-endpoints"]
  container_registry_url    = module.container_registry.container_registry_login_server
  docker_image_name         = "health-dashboard"
  docker_image_tag          = "latest"
  tags                      = local.primary_tags
  additional_app_settings   = {
    "ROLE" = "primary"
  }
}

module "web_app_secondary" {
  source                    = "../../modules/web-app"
  resource_group_name       = azurerm_resource_group.rg.name
  location                  = var.secondary_location
  service_plan_name         = "asp-${var.app_name}-${var.secondary_location}-${var.environment}"
  web_app_name             = "app-${var.app_name}-${var.secondary_location}-${var.environment}"
  subnet_id                = module.networking_secondary.subnet_ids["subnet-app"]
  private_endpoint_subnet_id = module.networking_secondary.subnet_ids["subnet-endpoints"]
  container_registry_url    = module.container_registry.container_registry_login_server
  docker_image_name         = "health-dashboard"
  docker_image_tag          = "latest"
  tags                      = local.secondary_tags
  additional_app_settings   = {
    "ROLE" = "secondary"
  }
}