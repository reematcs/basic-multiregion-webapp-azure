variable "resource_group_name" {
  type        = string
  description = "Name of the resource group"
}

variable "location" {
  type        = string
  description = "Azure region for the web app"
}

variable "service_plan_name" {
  type        = string
  description = "Name of the App Service Plan"
}

variable "web_app_name" {
  type        = string
  description = "Name of the Web App"
}

variable "sku_name" {
  type        = string
  description = "SKU name for the App Service Plan"
  default     = "P1v2"
}

variable "subnet_id" {
  type        = string
  description = "Subnet ID for the web app's private endpoint"
}

variable "container_registry_url" {
  type        = string
  description = "URL of the container registry"
}

variable "docker_image_name" {
  type        = string
  description = "Name of the docker image"
}

variable "docker_image_tag" {
  type        = string
  description = "Tag of the docker image"
  default     = "latest"
}

variable "additional_app_settings" {
  type        = map(string)
  description = "Additional app settings for the web app"
  default     = {}
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to resources"
  default     = {}
}