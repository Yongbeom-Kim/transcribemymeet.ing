resource "google_compute_managed_ssl_certificate" "lb" {
  provider = google-beta
  name     = var.frontend_bucket_name

  managed {
    domains = [var.domain]
  }
}

resource "google_compute_global_address" "static" {
  provider = google-beta
  name     = var.frontend_bucket_name
}

## Backend bucket
resource "google_compute_backend_bucket" "frontend" {
  name        = var.frontend_bucket_name
  description = "The bucket to store the frontend files"
  bucket_name = google_storage_bucket.frontend.name
  enable_cdn  = true

  cdn_policy {
    cache_mode                   = "CACHE_ALL_STATIC"
    client_ttl                   = 3600
    default_ttl                  = 3600
    max_ttl                      = 86400
    negative_caching             = true
    request_coalescing           = true
    serve_while_stale            = 86400
    signed_url_cache_max_age_sec = 0
  }
}

## Redirect http to https
resource "google_compute_url_map" "http-https-redirect" {
  provider    = google-beta
  name        = "${var.frontend_bucket_name}-redirect"
  description = "Redirects HTTP to HTTPS"

  default_url_redirect {
    https_redirect         = true
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
    strip_query            = false
  }
}

resource "google_compute_target_http_proxy" "http-redirect-roxy" {
  provider = google-beta
  name     = "${var.frontend_bucket_name}-redirect-proxy"
  url_map  = google_compute_url_map.http-https-redirect.id
}

resource "google_compute_global_forwarding_rule" "http-https-redirect" {
  provider = google-beta
  name     = "${var.frontend_bucket_name}-redirect"

  ip_address            = google_compute_global_address.static.id
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL"
  port_range            = "80"
  target                = google_compute_target_https_proxy.lb.id

}

## Load balancer
resource "google_compute_url_map" "lb" {
  provider        = google-beta
  name            = var.frontend_bucket_name
  description     = "LB with backend bucket backend"
  default_service = google_compute_backend_bucket.frontend.id
}

resource "google_compute_target_https_proxy" "lb" {
  provider      = google-beta
  name          = var.frontend_bucket_name
  quic_override = "NONE"
  url_map       = google_compute_url_map.lb.id
  ssl_certificates = [
    google_compute_managed_ssl_certificate.lb.id
  ]
}

resource "google_compute_global_forwarding_rule" "lb" {
  provider = google-beta
  name     = var.frontend_bucket_name

  ip_protocol           = "TCP"
  ip_address            = google_compute_global_address.static.id
  load_balancing_scheme = "EXTERNAL"
  port_range            = "443"
  target                = google_compute_target_https_proxy.lb.id

}
