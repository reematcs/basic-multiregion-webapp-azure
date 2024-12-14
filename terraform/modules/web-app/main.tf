resource "azurerm_service_plan" "plan" {
  name                = var.service_plan_name
  resource_group_name = var.resource_group_name
  location            = var.location
  os_type             = "Linux"
  sku_name            = var.sku_name

  tags = var.tags
}

resource "azurerm_linux_web_app" "app" {
  name                = var.web_app_name
  resource_group_name = var.resource_group_name
  location            = var.location
  service_plan_id     = azurerm_service_plan.plan.id

  virtual_network_subnet_id = var.subnet_id


  identity {
    type = "SystemAssigned"
  }
  site_config {
    container_registry_use_managed_identity = true

    application_stack {
      docker_image_name   = "${var.docker_image_name}:${var.docker_image_tag}"
      docker_registry_url = "https://${var.container_registry_url}"
    }

    health_check_path = "/api/health/live"
  }

 # In modules/web-app/main.tf, modify the app_settings in azurerm_linux_web_app:

app_settings = merge({
  "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
  "DOCKER_REGISTRY_SERVER_URL"          = var.container_registry_url
  "AZURE_REGION"                        = var.location
  "CONTAINER_VERSION"                   = var.docker_image_tag
  "APPLICATIONINSIGHTS_CONNECTION_STRING" = azurerm_application_insights.app_insights.connection_string
  "ApplicationInsightsAgent_EXTENSION_VERSION" = "~3"
}, var.additional_app_settings)

  tags = var.tags
}
resource "azurerm_private_endpoint" "webapp_pe" {
  name                = "${var.web_app_name}-pe"
  location            = var.location
  resource_group_name = var.resource_group_name
  subnet_id           = var.private_endpoint_subnet_id # Use the new variable

  private_service_connection {
    name                           = "${var.web_app_name}-privateserviceconnection"
    private_connection_resource_id = azurerm_linux_web_app.app.id
    subresource_names              = ["sites"]
    is_manual_connection           = false
  }

  tags = var.tags
}

# In modules/web-app/main.tf

resource "azurerm_application_insights" "app_insights" {
  name                = "${var.web_app_name}-appinsights"  # Use web_app_name instead of name
  resource_group_name = var.resource_group_name
  location            = var.location
  application_type    = "web"
  tags                = var.tags  # Add tags to match your pattern
}