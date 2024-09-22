variable "project" {
  type        = string
  description = "The project ID to deploy resources"
}

variable "region" {
  type        = string
  description = "The region to deploy resources"
}

variable "service" {
  type        = string
  description = "The service name of this microservice."
}

variable "tf_backend_bucket" {
  type        = string
  description = "The name of the bucket to store the Terraform state file"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.3.0"
    }
  }
  backend "gcs" {
    bucket  = var.tf_backend_bucket
    prefix  = "terraform/state/${var.service}"
  }
}

provider "google" {
  # Configuration options
  project = var.project
  region  = var.region

}