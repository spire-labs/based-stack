package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

const (
	l1Url                    = "http://localhost:8545"
	l1Beacon                 = "http://localhost:5052"
	l1ChainId                = 900
	l2Url                    = "http://localhost:9545"
	l2ChainId                = 901
	deployerAddr             = "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
	deployerPrivKey          = "0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e"
	randomAddr               = "0xdead000000000000000000000000000000000000"
	l1BlockL2ContractAddress = "0x4200000000000000000000000000000000000015"
	optimismPortalAddress    = "0x6509f2a854BA7441039fCE3b959d5bAdd2fFCFCD"
	l2BlocksPerEpoch         = 3
)

func main() {
	app := cli.NewApp()
	app.Name = "Op devnet scripts"
	app.Description = "Collection of scripts to help analyse the op-stack devnet"

	app.Commands = []*cli.Command{
		{
			Name:  "l2-tx",
			Usage: "Sends a simple L2 Transaction",
			Action: func(_ *cli.Context) error {
				sendL2Tx()
				return nil
			},
		},
		{
			Name:  "l1-block",
			Usage: "Fetches data from L1Block contract",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:     "start",
					Usage:    "Start block, inclusive",
					Required: true,
				},
				&cli.IntFlag{
					Name:     "end",
					Usage:    "End block, exclusive",
					Required: true,
				},
			},
			Action: func(clx *cli.Context) error {
				start := clx.Int64("start")
				end := clx.Int64("end")
				field := clx.String("field")
				l1Block(field, start, end)
				return nil
			},
		},
		{
			Name:  "deposit",
			Usage: "Deposits a transaction to OptimismPortal",
			Action: func(_ *cli.Context) error {
				deposit()
				return nil
			},
		},
		{
			Name:  "decode",
			Usage: "Decodes batches from L2",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:     "start",
					Required: true,
					Usage:    "First block (inclusive) to fetch",
				},
				&cli.IntFlag{
					Name:     "end",
					Required: true,
					Usage:    "Last block (exclusive) to fetch",
				},
			},
			Action: func(clx *cli.Context) error {
				start := clx.Int("start")
				end := clx.Int("end")
				decodeBatch(start, end)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func l1Block(field string, start, end int64) {
	l2Client, err := ethclient.Dial(l2Url)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	maxConcurrent := 20
	sem := make(chan struct{}, maxConcurrent)
	txs := make([][]byte, end-start)

	ctx := context.Background()

	for i := int64(0); i < end-start; i++ {
		wg.Add(1)
		go func(index int64) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				// err?
				return
			}

			txs[index] = downloadSystemTx(ctx, start+index, *l2Client)
		}(i)
	}

	wg.Wait()

	for i, tx := range txs {
		blockNumber := start + int64(i)
		electionWinner := common.Bytes2Hex(tx[164:])
		timestamp := new(big.Int)
		timestamp.SetBytes(tx[20:28])
		l1Origin := new(big.Int)
		l1Origin.SetBytes(tx[28:36])
		fmt.Printf("L2 block: %3d, timestamp: %d, L1 origin: %d, Election Winner: %s\n", blockNumber, timestamp, l1Origin, electionWinner)
	}
}

func downloadSystemTx(ctx context.Context, blockNumber int64, l2Client ethclient.Client) []byte {
	log.Println("Downloading system tx from block", blockNumber)
	block, err := l2Client.BlockByNumber(ctx, big.NewInt(blockNumber))
	if err != nil {
		log.Fatal(err)
	}
	txs := block.Transactions()
	if len(txs) == 0 {
		log.Fatal("Empty transactions", blockNumber)
	}
	tx := txs[0]
	// if !tx.IsSystemTx() {
	// 	log.Fatalf("Not a system tx. block=%d, hash=%s", blockNumber, tx.Hash())
	// }

	return tx.Data()
}

func deposit() {
	privateKey, err := crypto.HexToECDSA(deployerPrivKey[2:])
	deployerAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(l1ChainId))
	if err != nil {
		log.Fatal(err)
	}
	l1Client, err := ethclient.Dial(l1Url)
	if err != nil {
		log.Fatal(err)
	}
	l2Client, err := ethclient.Dial(l2Url)
	if err != nil {
		log.Fatal(err)
	}
	nonce, err := l1Client.PendingNonceAt(context.Background(), deployerAddr)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = new(big.Int).SetUint64(nonce)
	auth.Value = big.NewInt(0)

	address := common.HexToAddress(optimismPortalAddress)
	instance, err := bindings.NewOptimismPortalTransactor(address, l1Client)
	if err != nil {
		log.Fatal(err)
	}

	l2BlockNumber, _ := l2Client.BlockNumber(context.Background())

	tx, err := instance.DepositTransaction(auth, common.HexToAddress(randomAddr), big.NewInt(1), uint64(300000), false, []byte{})
	if err != nil {
		log.Fatal(err)
	}
	txReceipt := waitForTransaction(tx.Hash(), l1Client)

	log.Println("On L1:")
	printDeployerBalance(l1Client)

	log.Println("On L2:")
	printDeployerBalance(l2Client)

	confirmL2BlockWithDeposit(l2BlockNumber, txReceipt.BlockNumber.Uint64(), deployerAddr, l2Client)
}

func confirmL2BlockWithDeposit(startL2Block, depositL1Block uint64, deployerAddr common.Address, l2Client *ethclient.Client) {
	// get data from system contract
	address := common.HexToAddress(l1BlockL2ContractAddress)
	l1BlockContract, err := bindings.NewL1BlockCaller(address, l2Client)
	if err != nil {
		log.Fatal(err)
	}

	l1Origin, err := l1BlockContract.Number(&bind.CallOpts{BlockNumber: new(big.Int).SetUint64(startL2Block)})
	if err != nil {
		log.Fatal(err)
	}
	seqNumber, err := l1BlockContract.SequenceNumber(&bind.CallOpts{BlockNumber: new(big.Int).SetUint64(startL2Block)})
	if err != nil {
		log.Fatal(err)
	}

	// calculate where the deposit should land
	startOfEpoch := startL2Block - seqNumber
	l2BlockWithDeposit := startOfEpoch + (depositL1Block-l1Origin)*l2BlocksPerEpoch
	log.Println("Deposit should be included in block", l2BlockWithDeposit)
	waitForBlock(l2BlockWithDeposit, l2Client)

	l2BalanceBefore, _ := l2Client.BalanceAt(context.Background(), deployerAddr, new(big.Int).SetUint64(l2BlockWithDeposit-1))
	l2BalanceAfter, _ := l2Client.BalanceAt(context.Background(), deployerAddr, new(big.Int).SetUint64(l2BlockWithDeposit))

	log.Printf("L2 balance before: %v", l2BalanceBefore)
	log.Printf("L2 balance after: %v", l2BalanceAfter)

	if l2BalanceBefore.Cmp(l2BalanceAfter) != 1 {
		log.Fatal("L2 balance did not change in this block, somethings off")
	}

	log.Println("All good!")
}

func waitForBlock(blockNumber uint64, client *ethclient.Client) {
	log.Println("Waiting for block ", blockNumber)
	for {
		curr, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if curr >= blockNumber {
			break
		}
		time.Sleep(time.Second * (time.Duration(blockNumber - curr)))
	}
	log.Printf("Block %v included", blockNumber)
}

func sendL2Tx() {
	l2Client, err := ethclient.Dial(l2Url)
	if err != nil {
		log.Fatalf("Error dialing ethereum client: %v", err)
	}

	printDeployerBalance(l2Client)

	tx := signTransaction(randomAddr, "", deployerPrivKey, uint64(21400), uint64(1), l2Client)
	err = l2Client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalln(err)
	}
	receipt := waitForTransaction(tx.Hash(), l2Client)
	printDeployerBalance(l2Client)
	l1Origin := getL1Origin(receipt.BlockHash, l2Client)
	log.Println("L1Origin: ", l1Origin)
}

func getL1Origin(l2BlockHash common.Hash, l2Client *ethclient.Client) uint64 {
	address := common.HexToAddress(l1BlockL2ContractAddress)
	instance, err := bindings.NewL1BlockCaller(address, l2Client)
	if err != nil {
		log.Fatal(err)
	}

	l1Origin, err := instance.Number(&bind.CallOpts{BlockHash: l2BlockHash})
	if err != nil {
		log.Fatal(err)
	}

	return l1Origin
}

func printDeployerBalance(client *ethclient.Client) {
	ctx := context.Background()
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Block number: %v", blockNumber)

	balance, err := client.BalanceAt(ctx, common.HexToAddress(deployerAddr), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Deployer balance: %v", balance)
}
