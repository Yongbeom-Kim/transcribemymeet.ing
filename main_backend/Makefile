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
		go test ./...

tofu_run_cmd: ## Run a tofu command. Usage: make tofu_run_cmd CMD="command"
	. ./source.sh && \
		cd tf && \
		tofu init && \
		tofu workspace select -or-create $${TERRAFORM_WORKSPACE} && \
		tofu $(CMD)

code_deploy: ## Deploy the backend code

infra_deploy: ## Deploy the infrastructure
	$(MAKE) tofu_run_cmd CMD="apply -auto-approve"

infra_destroy: ## Destroy the infrastructure
	$(MAKE) tofu_run_cmd CMD="destroy -auto-approve"