resource "azurerm_service_plan" "plan" {
  name                = var.service_plan_name
  resource_group_name = var.resource_group_name
  location            = var.location
  os_type            = "Linux"
  sku_name           = var.sku_name

  tags = var.tags
}

resource "azurerm_linux_web_app" "app" {
  name                = var.web_app_name
  resource_group_name = var.resource_group_name
  location            = var.location
  service_plan_id     = azurerm_service_plan.plan.id

  virtual_network_subnet_id = var.subnet_id

  site_config {
    container_registry_use_managed_identity = true
    
    application_stack {
      docker_image     = "${var.container_registry_url}/${var.docker_image_name}"
      docker_image_tag = var.docker_image_tag
    }

    health_check_path = "/api/health/live"
  }

  identity {
    type = "SystemAssigned"
  }

  app_settings = merge({
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
    "DOCKER_REGISTRY_SERVER_URL"          = var.container_registry_url
    "AZURE_REGION"                        = var.location
    "CONTAINER_VERSION"                   = var.docker_image_tag
  }, var.additional_app_settings)

  tags = var.tags
}

# Add Private Endpoint for Web App
resource "azurerm_private_endpoint" "webapp_pe" {
  name                = "${var.web_app_name}-pe"
  location            = var.location
  resource_group_name = var.resource_group_name
  subnet_id           = var.subnet_id

  private_service_connection {
    name                           = "${var.web_app_name}-privateserviceconnection"
    private_connection_resource_id = azurerm_linux_web_app.app.id
    subresource_names             = ["sites"]
    is_manual_connection          = false
  }

  tags = var.tags
}