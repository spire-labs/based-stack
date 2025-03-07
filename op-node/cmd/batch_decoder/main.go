package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/fetch"
	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/reassemble"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "batch-decoder"
	app.Usage = "Optimism Batch Decoding Utility"
	app.Commands = []*cli.Command{
		{
			Name:  "fetch",
			Usage: "Fetches batches in the specified range",
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
				&cli.StringFlag{
					Name:     "inbox",
					Required: true,
					Usage:    "Batch Inbox Address",
				},
				&cli.StringFlag{
					Name:  "out",
					Value: "/tmp/batch_decoder/transactions_cache",
					Usage: "Cache directory for the found transactions",
				},
				&cli.StringFlag{
					Name:     "l1",
					Required: true,
					Usage:    "L1 RPC URL",
					EnvVars:  []string{"L1_RPC"},
				},
				&cli.StringFlag{
					Name:     "l1.beacon",
					Required: false,
					Usage:    "Address of L1 Beacon-node HTTP endpoint to use",
					EnvVars:  []string{"L1_BEACON"},
				},
				&cli.IntFlag{
					Name:  "concurrent-requests",
					Value: 10,
					Usage: "Concurrency level when fetching L1",
				},
			},
			Action: func(cliCtx *cli.Context) error {
				l1Client, err := ethclient.Dial(cliCtx.String("l1"))
				if err != nil {
					log.Fatal(err)
				}
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				chainID, err := l1Client.ChainID(ctx)
				if err != nil {
					log.Fatal(err)
				}
				beaconAddr := cliCtx.String("l1.beacon")
				var beacon *sources.L1BeaconClient
				if beaconAddr != "" {
					beaconClient := sources.NewBeaconHTTPClient(client.NewBasicHTTPClient(beaconAddr, nil))
					beaconCfg := sources.L1BeaconClientConfig{FetchAllSidecars: false}
					beacon = sources.NewL1BeaconClient(beaconClient, beaconCfg)
					_, err := beacon.GetVersion(ctx)
					if err != nil {
						log.Fatal(fmt.Errorf("failed to check L1 Beacon API version: %w", err))
					}
				} else {
					fmt.Println("L1 Beacon endpoint not set. Unable to fetch post-ecotone channel frames")
				}
				config := fetch.Config{
					Start:              uint64(cliCtx.Int("start")),
					End:                uint64(cliCtx.Int("end")),
					ChainID:            chainID,
					BatchInbox:         common.HexToAddress(cliCtx.String("inbox")),
					OutDirectory:       cliCtx.String("out"),
					ConcurrentRequests: uint64(cliCtx.Int("concurrent-requests")),
				}
				totalValid, totalInvalid := fetch.Batches(l1Client, beacon, config)
				fmt.Printf("Fetched batches in range [%v,%v). Found %v valid & %v invalid batches\n", config.Start, config.End, totalValid, totalInvalid)
				fmt.Printf("Fetch Config: Chain ID: %v. Inbox Address: %v. Valid Senders: %v.\n", config.ChainID, config.BatchInbox, config.BatchSenders)
				fmt.Printf("Wrote transactions with batches to %v\n", config.OutDirectory)
				return nil
			},
		},
		{
			Name:  "reassemble",
			Usage: "Reassembles channels from fetched batch transactions and decode batches",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "in",
					Value: "/tmp/batch_decoder/transactions_cache",
					Usage: "Cache directory for the found transactions",
				},
				&cli.StringFlag{
					Name:  "out",
					Value: "/tmp/batch_decoder/channel_cache",
					Usage: "Cache directory for the found channels",
				},
				&cli.Uint64Flag{
					Name:  "l2-chain-id",
					Value: 10,
					Usage: "L2 chain id for span batch derivation. Default value from op-mainnet.",
				},
				&cli.Uint64Flag{
					Name:  "l2-genesis-timestamp",
					Value: 1686068903,
					Usage: "L2 genesis time for span batch derivation. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
				&cli.Uint64Flag{
					Name:  "l2-block-time",
					Value: 2,
					Usage: "L2 block time for span batch derivation. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
				&cli.StringFlag{
					Name:  "inbox",
					Value: "0xFF00000000000000000000000000000000000010",
					Usage: "Batch Inbox Address. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
			},
			Action: func(cliCtx *cli.Context) error {
				var (
					L2GenesisTime     uint64         = cliCtx.Uint64("l2-genesis-timestamp")
					L2BlockTime       uint64         = cliCtx.Uint64("l2-block-time")
					BatchInboxAddress common.Address = common.HexToAddress(cliCtx.String("inbox"))
				)
				L2ChainID := new(big.Int).SetUint64(cliCtx.Uint64("l2-chain-id"))
				rollupCfg, err := rollup.LoadOPStackRollupConfig(L2ChainID.Uint64())
				if err == nil {
					// prioritize superchain config
					if L2GenesisTime != rollupCfg.Genesis.L2Time {
						L2GenesisTime = rollupCfg.Genesis.L2Time
						fmt.Printf("L2GenesisTime overridden: %v\n", L2GenesisTime)
					}
					if L2BlockTime != rollupCfg.BlockTime {
						L2BlockTime = rollupCfg.BlockTime
						fmt.Printf("L2BlockTime overridden: %v\n", L2BlockTime)
					}
					if BatchInboxAddress != rollupCfg.BatchInboxContractAddress {
						BatchInboxAddress = rollupCfg.BatchInboxContractAddress
						fmt.Printf("BatchInboxAddress overridden: %v\n", BatchInboxAddress)
					}
				}
				config := reassemble.Config{
					BatchInbox:    BatchInboxAddress,
					InDirectory:   cliCtx.String("in"),
					OutDirectory:  cliCtx.String("out"),
					L2ChainID:     L2ChainID,
					L2GenesisTime: L2GenesisTime,
					L2BlockTime:   L2BlockTime,
				}
				reassemble.Channels(config, rollupCfg)
				return nil
			},
		},
		{
			Name:  "reassemble-devnet",
			Usage: "Reassembles channels from fetched batch transactions and decode batches for devnet use",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "in",
					Value: "/tmp/batch_decoder/transactions_cache",
					Usage: "Cache directory for the found transactions",
				},
				&cli.StringFlag{
					Name:  "out",
					Value: "/tmp/batch_decoder/channel_cache",
					Usage: "Cache directory for the found channels",
				},
			},
			Action: func(cliCtx *cli.Context) error {
				rollupCfg, err := readDevnetConfig()
				if err != nil {
					log.Fatal(err)
				}

				var (
					L2GenesisTime     uint64         = rollupCfg.Genesis.L2Time
					L2BlockTime       uint64         = rollupCfg.BlockTime
					L2ChainID         *big.Int       = rollupCfg.L2ChainID
					BatchInboxAddress common.Address = rollupCfg.BatchInboxContractAddress
				)

				config := reassemble.Config{
					BatchInbox:    BatchInboxAddress,
					InDirectory:   cliCtx.String("in"),
					OutDirectory:  cliCtx.String("out"),
					L2ChainID:     L2ChainID,
					L2GenesisTime: L2GenesisTime,
					L2BlockTime:   L2BlockTime,
				}

				reassemble.Channels(config, rollupCfg)
				fmt.Printf("Reassembled channels saved to %v\n", cliCtx.String("out"))
				return nil
			},
		},
		{
			Name:  "force-close",
			Usage: "Create the tx data which will force close a channel",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Required: true,
					Usage:    "ID of the channel to close",
				},
				&cli.StringFlag{
					Name:  "inbox",
					Value: "0x0000000000000000000000000000000000000000",
					Usage: "(Optional) Batch Inbox Address",
				},
				&cli.StringFlag{
					Name:  "in",
					Value: "/tmp/batch_decoder/transactions_cache",
					Usage: "Cache directory for the found transactions",
				},
			},
			Action: func(cliCtx *cli.Context) error {
				var id derive.ChannelID
				if err := (&id).UnmarshalText([]byte(cliCtx.String("id"))); err != nil {
					log.Fatal(err)
				}
				frames := reassemble.LoadFrames(cliCtx.String("in"), common.HexToAddress(cliCtx.String("inbox")))
				var filteredFrames []derive.Frame
				for _, frame := range frames {
					if frame.Frame.ID == id {
						filteredFrames = append(filteredFrames, frame.Frame)
					}
				}
				data, err := derive.ForceCloseTxData(filteredFrames)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%x\n", data)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

const devnetCfgPath = "../../../packages/contracts-bedrock/deploy-config/devnetL1.json"

type DevnetConfig struct {
	L2ChainID               *big.Int       `json:"l2ChainID"`
	L2BlockTime             uint64         `json:"l2BlockTime"`
	BatchInboxAddress       common.Address `json:"batchInboxAddress"`
	L1GenesisBlockTimestamp string         `json:"l1GenesisBlockTimestamp"`
}

func readDevnetConfig() (cfg *rollup.Config, err error) {
	content, err := os.ReadFile(devnetCfgPath)
	if err != nil {
		return nil, fmt.Errorf("error when opening file: %w", err)
	}

	var config DevnetConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("error during Unmarshal(): %w", err)
	}

	L2Time, _ := strconv.ParseUint(config.L1GenesisBlockTimestamp[2:], 16, 64)

	rollupCfg := rollup.Config{
		Genesis: rollup.Genesis{
			L2Time: L2Time,
		},
		L2ChainID:                 config.L2ChainID,
		BatchInboxContractAddress: config.BatchInboxAddress,
		BlockTime:                 config.L2BlockTime,
	}

	return &rollupCfg, nil
}
