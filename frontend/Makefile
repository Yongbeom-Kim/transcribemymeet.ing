help: ## Show this help
	@echo "Usage: make [target]"
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deploy: infra_deploy code_deploy ## Deploy the frontend infrastructure + code
destroy: infra_destroy ## Destroy the frontend infrastructure

tofu_run_cmd: ## Run a tofu command. Usage: make tofu_run_cmd CMD="command"
	. ./source.sh && \
		cd tf && \
		tofu init && \
		tofu workspace select -or-create $${TERRAFORM_WORKSPACE} && \
		tofu $(CMD)

code_deploy: ## Deploy the frontend code to the bucket
	. ./source.sh && \
		yarn && \
		yarn run build && \
		gsutil -m rm -r gs://$${TF_VAR_frontend_bucket_name}/* || true && \
		gsutil -m cp -r dist/* gs://$${TF_VAR_frontend_bucket_name}

infra_deploy: ## Deploy the infrastructure
	$(MAKE) tofu_run_cmd CMD="apply -auto-approve"

infra_destroy: ## Destroy the infrastructure
	$(MAKE) tofu_run_cmd CMD="destroy -auto-approve"