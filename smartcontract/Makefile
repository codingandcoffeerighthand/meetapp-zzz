RPC = ws://localhost:8545
PRIVATE_KEY := ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
CHAINID = 31337

deploy-contract:
	forge build --sizes
	forge script script/meet.s.sol \
		--rpc-url $(RPC) --private-key $(PRIVATE_KEY) --broadcast
.PHONY: deploy-contract

WEB_ABI_PATH = ../web/src/abi/abi.json
BE_ABI_PATH = ../server/pkg/gencode/smc_gen/meet.json
gen-abi:
	forge inspect --json Meet abi > meet.json
	cp meet.json $(BE_ABI_PATH)
	cp meet.json $(WEB_ABI_PATH)
.PHONY: gen-abi


##################################################################
#### GET CONTRACT ADDRESS FROM JSON FILE AND UPDATE YAML FILE ####
##################################################################

# Tên file JSON chứa contractAddress
JSON_FILE = broadcast/meet.s.sol/${CHAINID}/run-latest.json

# Tên file YAML cần cập nhật
YAML_FILE = ../server/config/config.yaml

# Field trong file JSON chứa contractAddress
JSON_FIELD = .transactions[0].contractAddress

# Field trong file YAML cần cập nhật
YAML_FIELD = .web3.contract_address

# Biến để lưu trữ giá trị contractAddress
CONTRACT_ADDRESS = $(shell jq -r '$(JSON_FIELD)' '$(JSON_FILE)')


update-yaml:
	yq -i "$(YAML_FIELD) = \"$(CONTRACT_ADDRESS)\"" "$(YAML_FILE)" -Y

.PHONY: update-yaml

anvil-up:
	anvil --state ./cache/anvil_state.json
.PHONY: anvil-up

deploy: deploy-contract gen-abi update-yaml
.PHONY: deploy