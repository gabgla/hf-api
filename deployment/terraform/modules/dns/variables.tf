variable "hosted_zone_name" {
  description = "Name of the Hosted Zone, hosted in AWS"
  type        = string
}

variable "api_subdomain" {
  description = "Fully qualified subdomain of the hosted zone"
  type        = string
}

variable "tags" {
  description = "A map of tags to apply to resources"
  type        = map(string)
}
