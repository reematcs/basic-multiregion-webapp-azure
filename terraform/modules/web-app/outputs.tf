output "web_app_id" {
  value = azurerm_linux_web_app.app.id
}

output "web_app_name" {
  value = azurerm_linux_web_app.app.name
}

output "web_app_hostname" {
  value = azurerm_linux_web_app.app.default_hostname
}

output "identity_principal_id" {
  value       = azurerm_linux_web_app.app.identity.0.principal_id
  description = "The Principal ID for the System Assigned Identity of the web app"
}
output "private_endpoint_id" {
  value       = azurerm_private_endpoint.webapp_pe.id
  description = "The ID of the private endpoint"
}