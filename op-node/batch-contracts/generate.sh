#!/bin/sh

set -euo

forge build

# BatchTicketAccounting
cd out/BatchTicketAccounting.sol

cat BatchTicketAccounting.json | jq -r  '.bytecode.object' > BatchTicketAccounting.bin
cat BatchTicketAccounting.json | jq -r  '.abi' > BatchTicketAccounting.abi
cd ../..


abigen --abi ./out/BatchTicketAccounting.sol/BatchTicketAccounting.abi --bin ./out/BatchTicketAccounting.sol/BatchTicketAccounting.bin --pkg BatchTicketAccounting --out ./bindings/BatchTicketAccounting/batch_ticket_accounting.go


# BatchRandomTicketInstruction
cd out/BatchRandomTicketInstruction.sol

cat BatchRandomTicketInstruction.json | jq -r  '.bytecode.object' > BatchRandomTicketInstruction.bin
cat BatchRandomTicketInstruction.json | jq -r  '.abi' > BatchRandomTicketInstruction.abi

cd ../..

abigen --abi ./out/BatchRandomTicketInstruction.sol/BatchRandomTicketInstruction.abi --bin ./out/BatchRandomTicketInstruction.sol/BatchRandomTicketInstruction.bin --pkg BatchRandomTicketInstruction --out ./bindings/BatchRandomTicketInstruction/batch_random_ticket_instruction.go