variable "username" {
  default = "admin"
}
variable "password" {
  sensitive = true
}
variable "domain" {}
variable "insecure" {
  default = false
}
variable "psk" {
  sensitive = true
}
variable "client_isolation" {
  default = true
}
variable "band" {
  default = "both"
}
variable "zone" {}
variable "ssid" {}