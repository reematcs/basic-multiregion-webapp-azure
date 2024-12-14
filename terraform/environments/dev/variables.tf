variable "resource_group_name" {
  type        = string
  description = "Name of the resource group"
  default     = "rg-dev-azure"
}

variable "primary_location" {
  type        = string
  description = "Location for primary region resources"
  default     = "westus"
}

variable "secondary_location" {
  type        = string
  description = "Location for secondary region resources"
  default     = "centralus"
}

variable "storage_account_name" {
  type        = string
  description = "Name of the storage account (optional)"
  default     = null
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

variable "app_name" {
  type        = string
  description = "Base name for the application resources"
  default     = "healthapp"
}

variable "environment" {
  type        = string
  description = "Environment name for resource naming"
  default     = "dev"
}

variable "common_tags" {
  type        = map(string)
  description = "Common tags for all resources"
  default = {
    environment = "dev"
    managed_by  = "terraform"
  }
}