variable "frontend_bucket_name" {
  type        = string
  description = "The name of the bucket to store the frontend files"
}

resource "google_storage_bucket" "frontend" {
  name          = var.frontend_bucket_name
  location      = "EU"
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
  }
  cors {
    origin          = ["*"]
    method          = ["*"]
    response_header = ["*"]
    max_age_seconds = 3600
  }
}


resource "google_storage_bucket_iam_member" "public_access" {
  bucket = google_storage_bucket.frontend.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
