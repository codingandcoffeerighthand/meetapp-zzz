go-tools:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
.PHONY: go-tools # this will then modify your `go.mod`
	
go-generate:
	go generate ./...
.PHONY: go-generate
	
get-abi:
	pnpx solc --combined-json abi,bin ../smartcontract/src/meeting.sol > gencode/smartcontracts/meeting.abi
.PHONY: get-abi