TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=jubril.xyz
NAMESPACE=custom
NAME=payments
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_arm64

default: install

build:
	go build -o ${BINARY}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

clean:
	rm -rf ./examples/.terraform/ 
	rm ./examples/terraform.tfstate 
	rm ./examples/.terraform.lock.hcl

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m 
