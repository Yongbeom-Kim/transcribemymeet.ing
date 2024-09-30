variable "resource_name" {
  type        = string
  description = "The name of the resource"
}

variable "backend_identity_key" {
  type        = string
  description = "The path to the backend identity key"
}

resource "google_storage_bucket" "bucket" {
  name          = "${var.resource_name}"
  location      = "US"
  force_destroy = true

  uniform_bucket_level_access = true
}

# Identity to create presigned upload and download urls to the bucket.
resource "google_service_account" "bucket_identity" {
  account_id = "${substr(terraform.workspace, 0, 10)}-${var.service}"
  display_name = "Bucket Identity for ${var.resource_name}"
}

resource "google_storage_bucket_iam_member" "bucket_identity_member" {
  bucket = google_storage_bucket.bucket.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.bucket_identity.email}"
}

resource "random_id" "key_rotator" {
  byte_length = 8
  keepers = {
    # This will change every time, forcing a new key to be created
    rotation_time = timestamp()
  }
}

resource "google_service_account_key" "bucket_identity_key" {
  service_account_id = google_service_account.bucket_identity.name

  # Always delete the old key and create a new one
  # This is done so that we always create a new key
  keepers = {
    rotation_time = random_id.key_rotator.keepers.rotation_time
  }
}

resource "local_file" "bucket_identity_key" {
  content  = base64decode(google_service_account_key.bucket_identity_key.private_key)
  filename = var.backend_identity_key
}

