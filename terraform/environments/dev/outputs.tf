output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "key_vault_name" {
  value = module.key_vault.key_vault_name
}

output "storage_account_name" {
  value = module.storage.storage_account_name
}

output "container_registry_name" {
  value = module.container_registry.container_registry_name
}

output "container_registry_login_server" {
  value = module.container_registry.container_registry_login_server
}