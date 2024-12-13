output "storage_account_id" {
  value       = azurerm_storage_account.storage.id
  description = "ID of the storage account"
}

output "storage_account_name" {
  value       = azurerm_storage_account.storage.name
  description = "Name of the storage account"
}

output "primary_access_key" {
  value       = azurerm_storage_account.storage.primary_access_key
  description = "Primary access key for the storage account"
  sensitive   = true
}

output "primary_connection_string" {
  value       = azurerm_storage_account.storage.primary_connection_string
  description = "Primary connection string for the storage account"
  sensitive   = true
}

output "container_names" {
  value       = values(azurerm_storage_container.containers)[*].name
  description = "Names of the created containers"
}

output "storage_account_uri" {
  value       = azurerm_storage_account.storage.primary_blob_endpoint
  description = "Primary blob endpoint URI"
}