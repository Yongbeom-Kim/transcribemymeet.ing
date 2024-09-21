variable "tf_backend_bucket" {
    type        = string
    description = "The name of the bucket to store the Terraform state file"
}

resource "google_storage_bucket" "backend" {
    name = "${var.tf_backend_bucket}"
    location = "US"

    storage_class = "STANDARD"

    versioning {
      enabled = true
    }
}