default: e2e-apply

FEATURE  				?= examples/provider
FLAGS 					?= -auto-approve
INIT_FLAGS			?= -upgrade
VARS 						?= -var="api_url=$(COSMO_API_URL)" -var="api_key=$(COSMO_API_KEY)"

e2e-init:
	terraform -chdir=$(FEATURE) init $(INIT_FLAGS)
	terraform -chdir=$(FEATURE) validate

e2e-clean:
	rm -rf $(FEATURE)/.terraform
	rm -f $(FEATURE)/.terraform.lock.hcl
	rm -f $(FEATURE)/terraform.tfstate
	rm -f $(FEATURE)/terraform.tfstate.backup

e2e-apply: e2e-init
	terraform -chdir=$(FEATURE) plan $(VARS)
	terraform -chdir=$(FEATURE) apply $(FLAGS) $(VARS)

e2e-destroy:
	terraform -chdir=$(FEATURE) destroy $(FLAGS) $(VARS)
