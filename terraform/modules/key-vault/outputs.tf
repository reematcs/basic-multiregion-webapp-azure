output "resource_group_name" {
  value = var.resource_group_name
}
output "key_vault_name" {
  value = azurerm_key_vault.vault.name
}