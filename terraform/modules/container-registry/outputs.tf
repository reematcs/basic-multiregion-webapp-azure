output "container_registry_id" {
  value       = azurerm_container_registry.acr.id
  description = "The Container Registry ID"
}

output "container_registry_name" {
  value       = azurerm_container_registry.acr.name
  description = "The Container Registry name"
}

output "container_registry_login_server" {
  value       = azurerm_container_registry.acr.login_server
  description = "The Container Registry login server"
}

output "container_registry_admin_username" {
  value       = azurerm_container_registry.acr.admin_username
  description = "The Container Registry admin username"
  sensitive   = true
}

output "container_registry_admin_password" {
  value       = azurerm_container_registry.acr.admin_password
  description = "The Container Registry admin password"
  sensitive   = true
}