
tofu_run_cmd:
	. ./source.sh && \
	tofu init && \
	tofu $(CMD)

infra_deploy:
	$(MAKE) tofu_run_cmd CMD="apply -auto-approve"

infra_destroy:
	$(MAKE) tofu_run_cmd CMD="destroy -auto-approve"