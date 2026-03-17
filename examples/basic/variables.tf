variable "username" {
  default = "admin"
  type = string
}
variable "password" {
  sensitive = true
  type = string
}
variable "domain" {
  default = ""
  type = string
}
variable "insecure" {
  default = false
  type = bool
}
variable "psk" {
  sensitive = true
  type = string
}
variable "client_isolation" {
  default = true
  type = bool
}
variable "band" {
  default = "both"
  type = string
}
variable "controller" {
  type = string
}
variable "zone" {
  type = string
}
variable "username" {
  default = "admin"
  type = string
}
variable "password" {
  sensitive = true
  type = string
}
variable "domain" {
  default = ""
  type = string
}
variable "insecure" {
  default = false
  type = bool
}
variable "psk" {
  sensitive = true
  type = string
}
variable "client_isolation" {
  default = true
  type = bool
}
variable "band" {
  default = "both"
  type = string
}
variable "controller" {
  type = string
}
variable "zone" {
  type = string
}
variable "ssid" {
  type = string
}
variable "vlan" {
  type = number
}