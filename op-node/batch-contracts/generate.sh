#!/bin/sh

set -euo

forge build

cd out/BatchTicketAccounting.sol

cat BatchTicketAccounting.json | jq -r  '.bytecode.object' > BatchTicketAccounting.bin
cat BatchTicketAccounting.json | jq -r  '.abi' > BatchTicketAccounting.abi
cd ../..


abigen --abi ./out/BatchTicketAccounting.sol/BatchTicketAccounting.abi --bin ./out/BatchTicketAccounting.sol/BatchTicketAccounting.bin --pkg BatchTicketAccounting --out ./bindings/batch_ticket_accounting.go