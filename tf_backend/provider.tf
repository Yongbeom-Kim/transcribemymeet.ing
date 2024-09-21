variable "project" {
  type        = string
  description = "The project ID to deploy resources"
}

variable "region" {
  type        = string
  description = "The region to deploy resources"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.3.0"
    }
  }
}

provider "google" {
  # Configuration options
  project = var.project
  region  = var.region

}