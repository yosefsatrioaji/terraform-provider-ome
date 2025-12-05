TEST?=$$(go list ./... | grep -v 'vendor')
INSTALL_ROOT?=~/.terraform.d/plugins
HOSTNAME=registry.terraform.io
NAMESPACE=dell
NAME=ome
BINARY=terraform-provider-${NAME}
VERSION?=1.2.8
OSNAME=linux
OSARCH=amd64

default: install

build:
	mkdir -p build
	go mod tidy
	CGO_ENABLED=0 GOOS=${OSNAME} GOARCH=${OSARCH} go build -o build/${BINARY}_${VERSION} -trimpath

release:
	goreleaser release --clean

install: build
	mkdir -p ./terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OSNAME}_${OSARCH}
	mv build/${BINARY}_${VERSION} ./terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OSNAME}_${OSARCH}

test:
	go test $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=120m -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

package: clean build
	mkdir -p package/${VERSION}
	mv build/${BINARY}_${VERSION} package/${VERSION}/${BINARY}_${VERSION}
	cp README.md package/${VERSION}/README.md
	cp LICENSE package/${VERSION}/LICENSE
	chmod -R 777 package/${VERSION}/
	cd package/${VERSION}/ && zip -r terraform-provider-${NAME}_${VERSION}_${OSNAME}_${OSARCH}.zip ${BINARY}_${VERSION} README.md LICENSE
	cd package/${VERSION}/ && shasum -a 256 terraform-provider-${NAME}_${VERSION}_${OSNAME}_${OSARCH}.zip > terraform-provider-${NAME}_${VERSION}_SHA256SUMS

clean:
	rm -rf build
	rm -rf package
	rm -rf terraform.d