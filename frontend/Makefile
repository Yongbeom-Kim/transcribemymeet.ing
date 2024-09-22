tofu_run_cmd:
	. ./source.sh && \
		cd tf && \
		tofu init && \
		tofu workspace select -or-create $${TERRAFORM_WORKSPACE} && \
		tofu $(CMD)

code_deploy:
	. ./source.sh && \
		yarn && \
		yarn run build && \
		gsutil -m rm -r gs://$${TF_VAR_frontend_bucket_name}/* || true && \
		gsutil -m cp -r dist/* gs://$${TF_VAR_frontend_bucket_name}

infra_deploy:
	$(MAKE) tofu_run_cmd CMD="apply -auto-approve"

infra_destroy:
	$(MAKE) tofu_run_cmd CMD="destroy -auto-approve"