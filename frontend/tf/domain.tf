data "google_dns_managed_zone" "frontend" {
  name = "transcribemymeet-ing"
}

resource "google_dns_record_set" "frontend" {
  for_each     = toset(["${var.domain}.", "www.${var.domain}."])
  managed_zone = data.google_dns_managed_zone.frontend.name
  name         = each.value
  type         = "A"
  ttl          = 300
  rrdatas      = [google_compute_global_forwarding_rule.lb.ip_address]
}