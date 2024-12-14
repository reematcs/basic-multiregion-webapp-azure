variable "resource_group_name" {
  type        = string
  description = "Name of the resource group"
}

variable "location" {
  type        = string
  description = "Location for the networking resources"
}

variable "vnet_name" {
  type        = string
  description = "Name of the virtual network"
}

variable "address_space" {
  type        = list(string)
  description = "Address space for the virtual network"
  default     = ["10.0.0.0/16"]
}

variable "dns_servers" {
  type        = list(string)
  description = "Custom DNS servers"
  default     = []
}
variable "subnet_configurations" {
  type = list(object({
    name                          = string
    address_prefix                = string
    create_network_security_group = bool
    service_endpoints             = list(string)
    delegation = optional(object({
      name                       = string
      service_delegation_name    = string
      service_delegation_actions = list(string)
    }))
    nsg_rules = optional(list(object({
      name                       = string
      priority                   = number
      direction                  = string
      access                     = string
      protocol                   = string
      source_port_range          = string
      destination_port_range     = string
      source_address_prefix      = string
      destination_address_prefix = string
    })))
  }))
  description = "List of subnet configurations"
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to the networking resources"
  default     = {}
}