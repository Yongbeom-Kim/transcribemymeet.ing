#!/bin/bash
# General
export DOMAIN="transcribemymeet.ing"

# GCP
export GCP_PROJECT="transcribemymeet-ing"
export GCP_REGION="us-central1"
export GOOGLE_APPLICATION_CREDENTIALS="$(git rev-parse --show-toplevel)/secrets/gcloud_service_key.json"

# Terraform variables
export TF_VAR_project=$GCP_PROJECT
export TF_VAR_region=$GCP_REGION
export TF_VAR_tf_backend_bucket="${GCP_PROJECT}-tfstate"
export TF_VAR_domain=$DOMAIN
export TERRAFORM_WORKSPACE=$(git rev-parse --abbrev-ref HEAD | tr '/' '-')


# Setup
gcloud auth activate-service-account --key-file=$GOOGLE_APPLICATION_CREDENTIALS