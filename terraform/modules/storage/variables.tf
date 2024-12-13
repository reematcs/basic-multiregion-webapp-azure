variable "resource_group_name" {
  type        = string
  description = "Name of the resource group."
}

variable "location" {
  type        = string
  description = "Location of the resource group."
}
variable "storage_account_name" {
  type        = string
  description = "Name of the storage account (optional, will be generated if not provided)"
  default     = null
}

variable "account_tier" {
  type        = string
  description = "Tier of the storage account (Standard or Premium)"
  default     = "Standard"
  validation {
    condition     = contains(["Standard", "Premium"], var.account_tier)
    error_message = "Account tier must be either Standard or Premium."
  }
}

variable "account_replication_type" {
  type        = string
  description = "Type of replication for the storage account"
  default     = "LRS"
  validation {
    condition     = contains(["LRS", "GRS", "RAGRS", "ZRS", "GZRS", "RAGZRS"], var.account_replication_type)
    error_message = "Account replication type must be one of: LRS, GRS, RAGRS, ZRS, GZRS, RAGZRS."
  }
}

variable "container_names" {
  type        = list(string)
  description = "List of container names to create"
  default     = []
}

variable "enable_versioning" {
  type        = bool
  description = "Enable blob versioning"
  default     = true
}

variable "allowed_ips" {
  type        = list(string)
  description = "List of IP addresses allowed to access the storage account"
  default     = []
}

variable "default_network_rule" {
  type        = string
  description = "Default network rule (Allow or Deny)"
  default     = "Deny"
  validation {
    condition     = contains(["Allow", "Deny"], var.default_network_rule)
    error_message = "Default network rule must be either Allow or Deny."
  }
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to the storage account"
  default     = {}
}
