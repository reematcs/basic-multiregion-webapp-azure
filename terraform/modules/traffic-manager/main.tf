
resource "random_string" "tm_name" {
  length  = 8
  special = false
  upper   = false
}
resource "azurerm_traffic_manager_profile" "tm" {
  name                   = "tm-${var.traffic_manager_name}-${random_string.tm_name.result}"
  resource_group_name    = var.resource_group_name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "${var.dns_name}-${random_string.tm_name.result}"
    ttl           = 60
  }

  monitor_config {
    protocol                     = "HTTPS"
    port                         = 443
    path                         = "/api/health/live"
    interval_in_seconds          = 30
    timeout_in_seconds           = 10
    tolerated_number_of_failures = 3
  }

  tags = var.tags
}

resource "azurerm_traffic_manager_azure_endpoint" "primary" {
  name               = "primary-endpoint"
  profile_id         = azurerm_traffic_manager_profile.tm.id
  priority           = 1
  target_resource_id = var.primary_web_app_id
}

resource "azurerm_traffic_manager_azure_endpoint" "secondary" {
  name               = "secondary-endpoint"
  profile_id         = azurerm_traffic_manager_profile.tm.id
  priority           = 2
  target_resource_id = var.secondary_web_app_id
}