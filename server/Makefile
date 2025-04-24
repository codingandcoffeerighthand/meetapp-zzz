go-tools:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	
go-generate:
	go generate ./...
	
RSA_FILE := config/rsa_key
gen-rsa-key:
	go run cmd/generate-rsa/main.go ${RSA_FILE}
.PHONY: gen-rsa-key