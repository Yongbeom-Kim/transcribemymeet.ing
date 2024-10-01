#!/bin/bash

set -a

. ./.env

if [[ -f ../source.sh ]]; then
	. ../source.sh
fi

GOOGLE_APPLICATION_CREDENTIALS=$(realpath ../secrets/gcloud_service_key.json)
TF_VAR_service=main-backend
TF_VAR_resource_name=${TF_VAR_project}-$(echo ${TERRAFORM_WORKSPACE} | cut -c 1-10)-${TF_VAR_service}
TF_VAR_backend_identity_key=./gcloud_backend_identity.json

PORT=8080

set +a