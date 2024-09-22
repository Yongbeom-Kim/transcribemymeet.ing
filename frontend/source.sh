#!/bin/bash
. ../source.sh

export TF_VAR_service=frontend

export TF_VAR_frontend_bucket_name=$TERRAFORM_WORKSPACE.$DOMAIN