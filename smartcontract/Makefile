RPC := https://eth-sepolia.public.blastapi.io 
PRIVKEY := 70c40e1e7b43ead48cf4d66d04b420bae58f108d65d2a2b0efe7d8b820bf89ca

deploy-contract:
	forge script script/meeting.s.sol \
		--rpc-url $(RPC) --private-key $(PRIVKEY) --broadcast
.PHONY: deploy-contract
			
