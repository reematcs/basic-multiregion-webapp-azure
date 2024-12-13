output "traffic_manager_profile_id" {
  value = azurerm_traffic_manager_profile.tm.id
}

output "traffic_manager_fqdn" {
  value = azurerm_traffic_manager_profile.tm.fqdn
}