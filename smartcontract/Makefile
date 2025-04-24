RPC = ws://localhost:8545
PRIVATE_KEY := ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
CHAINID = 31337

deploy-contract:
	forge build
	forge script script/meeting.s.sol \
		--rpc-url $(RPC) --private-key $(PRIVATE_KEY) --broadcast
.PHONY: deploy-contract

WEB_ABI_PATH = ../web/src/abi/abi.json
BE_ABI_PATH = ../server/pkg/gencode/smc_gen/abi.json
gen-abi:
	forge inspect --json DAppMeeting abi > meeting.json
	cp meeting.json $(BE_ABI_PATH)
	cp meeting.json $(WEB_ABI_PATH)
.PHONY: gen-abi


##################################################################
#### GET CONTRACT ADDRESS FROM JSON FILE AND UPDATE YAML FILE ####
##################################################################

# Tên file JSON chứa contractAddress
JSON_FILE = broadcast/meeting.s.sol/${CHAINID}/run-latest.json

# Tên file YAML cần cập nhật
YAML_FILE = ../server/config/config.yaml

# Field trong file JSON chứa contractAddress
JSON_FIELD = .transactions[0].contractAddress

# Field trong file YAML cần cập nhật
YAML_FIELD = .web3.contract_address

# Biến để lưu trữ giá trị contractAddress
CONTRACT_ADDRESS = $(shell jq -r '$(JSON_FIELD)' '$(JSON_FILE)')


update_yaml:
	yq -i "$(YAML_FIELD) = \"$(CONTRACT_ADDRESS)\"" "$(YAML_FILE)" -Y

.PHONY: update_yaml

anvil-up:
	anvil --state ./cache/anvil_state.json
.PHONY: anvil-up

deploy: deploy-contract gen-abi update_yaml
.PHONY: deploy