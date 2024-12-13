variable "resource_group_name" {
  type        = string
  description = "Prefix of the resource group name that's combined with a random ID so name is unique in your Azure subscription."
  default     = "rg-dev-azure"
}

variable "resource_group_location" {
  type        = string
  description = "Location for all resources."
  default     = "eastus"
}

variable "storage_account_name" {
  type        = string
  description = "Name of the storage account (optional)"
  default     = null
}

variable "allowed_ip_addresses" {
  type        = list(string)
  description = "List of IP addresses allowed to access the storage account"
  default     = []
}

variable "key_vault_name" {
  type        = string
  description = "Name of the key vault"
  default     = null
}

variable "container_registry_name" {
  type        = string
  description = "Name of the container registry"
  default     = null
}

variable "traffic_manager_name" {
  type        = string
  description = "Name of the Traffic Manager profile"
  default     = "tm-healthapp-dev"
}

variable "dns_name" {
  type        = string
  description = "DNS name for the Traffic Manager profile"
  default     = "healthapp-dev"
}