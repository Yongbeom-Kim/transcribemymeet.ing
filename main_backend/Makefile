.PHONY: help deploy infra_deploy infra_destroy code_deploy test tofu_run_cmd

SHELL := /bin/bash

help: ## Show this help
	@echo "Usage: make [target]"
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deploy: infra_deploy code_deploy ## Deploy the backend infrastructure + code
destroy: infra_destroy ## Destroy the backend infrastructure

test: ## Run tests
	. ./source.sh && \
		TF_VAR_backend_identity_key=$$(realpath $${TF_VAR_backend_identity_key}) && \
		go test ./...

integration_test: ## Run integration tests
	. ./source.sh && \
		cd test && \
		go test ./...

tofu_run_cmd: ## Run a tofu command. Usage: make tofu_run_cmd CMD="command"
	. ./source.sh && \
		cd tf && \
		tofu init && \
		tofu workspace select -or-create $${TERRAFORM_WORKSPACE} && \
		tofu $(CMD)

local: ## Run the backend locally
	. ./source.sh && \
		docker run -p $${PORT}:$${PORT} transcribemymeet.ing-backend # TODO: this needs to be the image name for a gcloud docker repo

code_build: ## Build the backend code
	. ./source.sh && \
		docker build \
			-t transcribemymeet.ing-backend \
			--build-arg TF_VAR_backend_identity_key=$${TF_VAR_backend_identity_key} \
			--build-arg TF_VAR_resource_name=$${TF_VAR_resource_name} \
			.


code_deploy: ## Deploy the backend code

infra_deploy: ## Deploy the infrastructure
	$(MAKE) tofu_run_cmd CMD="apply -auto-approve"

infra_destroy: ## Destroy the infrastructure
	$(MAKE) tofu_run_cmd CMD="destroy -auto-approve"