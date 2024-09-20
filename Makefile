NAMESPACE						 = wundergraph
NAME							 = cosmo
BINARY							 = terraform-provider-${NAME}
VERSION 						 = 0.0.1
OS_ARCH  						 = linux_amd64
EXAMPLES   						 = examples

TEST							?= $$(go list ./... | grep -v 'vendor')
HOSTNAME						?= terraform.local

COSMO_API_URL					 ?= http://localhost:3001
COSMO_API_KEY					 ?= cosmo_669b576aaadc10ee1ae81d9193425705

default: testacc

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v -timeout 120m

.PHONY: test-go
test-go:
	go test $(TEST) -v 


.PHONY: test
test: clean build install testacc e2e

generate:
	go generate ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...
	terraform fmt -recursive 

build:
	go build -o bin/${BINARY}

install:
	rm -f examples/**/.terraform.lock.hcl
	rm -f ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv bin/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

clean-local:
	rm -rf bin
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

build-all-arches:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

release: generate build-all-arches

include examples/Makefile

.PHONY: e2e-apply-cd e2e-destroy-cd e2e-clean-cd
.PHONY: e2e-apply-cosmo e2e-destroy-cosmo e2e-clean-cosmo
.PHONY: e2e-apply-cosmo-monograph e2e-destroy-cosmo-monograph e2e-clean-cosmo-monograph
.PHONY: e2e-cd e2e-cosmo e2e-cosmo-monograph clean

e2e-apply-cd:
	rm -rf examples/provider/.terraform.lock.hcl
	FEATURE=examples/provider make e2e-init
	FEATURE=examples/provider make e2e-apply 

e2e-destroy-cd: 
	make e2e-destroy 

e2e-clean-cd: 
	make e2e-clean 

e2e-apply-cosmo: 
	rm -rf examples/guides/cosmo/.terraform.lock.hcl
	FEATURE=examples/guides/cosmo make e2e-init 
	FEATURE=examples/guides/cosmo make e2e-apply 

e2e-destroy-cosmo: 
	FEATURE=examples/guides/cosmo make e2e-destroy 

e2e-clean-cosmo: 
	FEATURE=examples/guides/cosmo make e2e-clean

e2e-apply-cosmo-monograph: 
	rm -rf examples/guides/cosmo-monograph/.terraform.lock.hcl
	FEATURE=examples/guides/cosmo-monograph make e2e-init 
	FEATURE=examples/guides/cosmo-monograph make e2e-apply 

e2e-destroy-cosmo-monograph: 
	FEATURE=examples/guides/cosmo-monograph make e2e-destroy 

e2e-clean-cosmo-monograph: 
	FEATURE=examples/guides/cosmo-monograph make e2e-clean

e2e-apply-cosmo-monograph-contract: 
	rm -rf examples/guides/cosmo-monograph-contract/.terraform.lock.hcl
	FEATURE=examples/guides/cosmo-monograph-contract make e2e-init 
	FEATURE=examples/guides/cosmo-monograph-contract make e2e-apply 

e2e-destroy-cosmo-monograph-contract: 
	FEATURE=examples/guides/cosmo-monograph-contract make e2e-destroy 

e2e-clean-cosmo-monograph-contract: 
	FEATURE=examples/guides/cosmo-monograph-contract make e2e-clean

## Cosmo Local
# Full example installing cosmo locally with a minikube kubernetes cluster 
# This will also deploy a router and configure it to use the generated router token
# Ensure to update your /etc/hosts file with
# output "hosts" generated after apply

e2e-apply-cosmo-local: 
	rm -rf examples/guides/cosmo-local/.terraform.lock.hcl
	FEATURE=examples/guides/cosmo-local make e2e-init 
	FEATURE=examples/guides/cosmo-local make e2e-apply 

e2e-destroy-cosmo-local: 
	FEATURE=examples/guides/cosmo-local make e2e-destroy 

e2e-clean-cosmo-local: 
	FEATURE=examples/guides/cosmo-local make e2e-clean

## Convenience targets to run specific e2e tests

e2e-cd: e2e-apply-cd e2e-destroy-cd
e2e-cosmo: e2e-apply-cosmo e2e-destroy-cosmo
e2e-cosmo-monograph: e2e-apply-cosmo-monograph e2e-destroy-cosmo-monograph
e2e-cosmo-monograph-contract: e2e-apply-cosmo-monograph-contract e2e-destroy-cosmo-monograph-contract
e2e-cosmo-local: e2e-apply-cosmo-local e2e-destroy-cosmo-local

e2e: e2e-cd e2e-cosmo e2e-cosmo-monograph e2e-cosmo-monograph-contract e2e-cosmo-local

clean: e2e-clean-cd e2e-clean-cosmo e2e-clean-cosmo-monograph e2e-clean-cosmo-monograph-contract e2e-clean-cosmo-local clean-local
destroy: e2e-destroy-cd e2e-destroy-cosmo e2e-destroy-cosmo-monograph e2e-destroy-cosmo-monograph-contract e2e-destroy-cosmo-local