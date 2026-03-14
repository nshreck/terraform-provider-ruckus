terraform {
  required_providers {
    ruckus = {
      source  = "shreckn/ruckus"
      version = "0.0.1"
    }
  }
}
provider "ruckus" {
  username             = var.username
  password             = var.password
  domain               = var.domain
  api_version          = "v13_1"
  insecure_skip_verify = var.insecure
}

data "ruckus_zone" "hq" {
  name = var.zone
}

resource "ruckus_wlan" "corp" {
  zone_id     = data.ruckus_zone.hq.id
  name        = var.ssid
  ssid        = "${var.ssid}WLAN"

  security {
    mode        = "wpa2_psk"
    passphrase  = var.psk
    encryption  = "ccmp"
  }

  vlan {
    access_vlan  = 120
    dynamic_vlan = false
  }

  radio {
    band             = var.band
    client_isolation = var.client_isolation
  }

  advanced {
    min_bss_rate = 6000
    ofdma        = true
  }
}