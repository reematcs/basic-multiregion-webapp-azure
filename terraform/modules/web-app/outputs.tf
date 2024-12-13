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
  value = azurerm_linux_web_app.app.identity[0].principal_id
}