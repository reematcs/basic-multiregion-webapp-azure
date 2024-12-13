
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

module "networking" {
  source = "../../modules/networking"

  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  vnet_name          = null  # Explicitly set to null to use auto-generation

  subnet_configurations = [
    {
      name                          = "subnet-1"
      address_prefix                = "10.0.1.0/24"
      create_network_security_group = true
      service_endpoints             = ["Microsoft.KeyVault", "Microsoft.Storage"]
    }
  ]

  tags = {
    environment = "dev"
    managed_by  = "terraform"
  }
}
module "container_registry" {
  source = "../../modules/container-registry"

  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  
  # Name is now optional and will be auto-generated if not provided
  # container_registry_name = "mycustomname" # Optional
  
  sku           = "Standard"
  admin_enabled = true
  
  tags = {
    environment = "dev"
    managed_by  = "terraform"
  }
}
module "traffic_manager" {
  source = "../../modules/traffic-manager"

  resource_group_name   = azurerm_resource_group.rg.name
  traffic_manager_name  = "tm-healthapp-dev"
  dns_name             = "healthapp-dev"
  primary_web_app_id   = module.web_app_east.web_app_id
  secondary_web_app_id = module.web_app_central.web_app_id
  
  tags = {
    environment = "dev"
    managed_by  = "terraform"
  }
}

module "web_app_east" {
  source = "../../modules/web-app"

  resource_group_name     = azurerm_resource_group.rg.name
  location               = "eastus"
  service_plan_name      = "asp-healthapp-east-dev"
  web_app_name           = "app-healthapp-east-dev"
  subnet_id              = module.networking.subnet_ids["subnet-1"]
  container_registry_url = module.container_registry.container_registry_login_server
  docker_image_name      = "health-dashboard"
  docker_image_tag       = "latest"

  additional_app_settings = {
    "ROLE" = "primary"
  }

  tags = {
    environment = "dev"
    managed_by  = "terraform"
    region      = "eastus"
  }
}

module "web_app_central" {
  source = "../../modules/web-app"

  resource_group_name     = azurerm_resource_group.rg.name
  location               = "centralus"
  service_plan_name      = "asp-healthapp-central-dev"
  web_app_name           = "app-healthapp-central-dev"
  subnet_id              = module.networking.subnet_ids["subnet-1"]
  container_registry_url = module.container_registry.container_registry_login_server
  docker_image_name      = "health-dashboard"
  docker_image_tag       = "latest"

  additional_app_settings = {
    "ROLE" = "secondary"
  }

  tags = {
    environment = "dev"
    managed_by  = "terraform"
    region      = "centralus"
  }
}