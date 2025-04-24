package smc_gen

//go:generate go tool abigen --v2 --abi abi.json --pkg smc_gen --type Meeeting --out meeting_gen.go
