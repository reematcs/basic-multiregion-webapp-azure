variable "resource_group_name" {
  type        = string
  description = "Name of the resource group"
}

variable "location" {
  type        = string
  description = "Location for the Container Registry"
  default     = "eastus"
}

variable "container_registry_name" {
  type        = string
  description = "Name of the Container Registry (optional)"
  default     = null
}


variable "sku" {
  type        = string
  description = "SKU of the Container Registry (Basic, Standard, or Premium)"
  default     = "Standard"
  validation {
    condition     = contains(["Basic", "Standard", "Premium"], var.sku)
    error_message = "SKU must be one of: Basic, Standard, Premium."
  }
}

variable "admin_enabled" {
  type        = bool
  description = "Enable admin user"
  default     = false
}

variable "georeplication_locations" {
  type        = list(string)
  description = "List of locations for geo-replication (only available for Premium SKU)"
  default     = []
}

variable "network_rule_default_action" {
  type        = string
  description = "Default action for network rules (only available for Premium SKU)"
  default     = "Allow"
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to the Container Registry"
  default     = {}
}