package main

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func waitForTransaction(hash common.Hash, client *ethclient.Client) *types.Receipt {
	log.Printf("Waiting for transaction %v\n", hash)
	for {
		_, isPending, err := client.TransactionByHash(context.Background(), hash)
		if err != nil {
			log.Fatalf("Error waiting for transaction: %v", err)
		}
		if !isPending {
			break
		}
		time.Sleep(1 * time.Second)
	}

	txReceipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		log.Fatalf("Error fetching transaction receipt: %v", err)
	}

	log.Printf("Transaction included in block number %v", txReceipt.BlockNumber)
	return txReceipt
}
