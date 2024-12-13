variable "resource_group_name" {
  type        = string
  description = "Name of the resource group"
}

variable "traffic_manager_name" {
  type        = string
  description = "Name of the Traffic Manager profile"
}

variable "dns_name" {
  type        = string
  description = "DNS name for the Traffic Manager profile"
}

variable "primary_web_app_id" {
  type        = string
  description = "Resource ID of the primary web app"
}

variable "secondary_web_app_id" {
  type        = string
  description = "Resource ID of the secondary web app"
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to the Traffic Manager profile"
  default     = {}
}