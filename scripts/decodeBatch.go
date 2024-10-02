package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

type FileData struct {
	ID             string   `json:"id"`
	IsReady        bool     `json:"is_ready"`
	InvalidFrames  bool     `json:"invalid_frames"`
	InvalidBatches bool     `json:"invalid_batches"`
	Frames         []Frame  `json:"frames"`
	Batches        []Batch  `json:"batches"`
	BatchTypes     []int    `json:"batch_types"`
	ComprAlgos     []string `json:"compr_algos"`
}

type Frame struct {
	TransactionHash string    `json:"transaction_hash"`
	InclusionBlock  int       `json:"inclusion_block"`
	Timestamp       int64     `json:"timestamp"`
	BlockHash       string    `json:"block_hash"`
	FrameData       FrameData `json:"frame"`
}

type FrameData struct {
	ID          string `json:"id"`
	FrameNumber int    `json:"frame_number"`
	Data        string `json:"data"`
	IsLast      bool   `json:"is_last"`
}

type Batch struct {
	ParentCheck       []string           `json:"parent_check"`
	L1OriginCheck     []string           `json:"l1_origin_check"`
	SpanBatchElements []SpanBatchElement `json:"span_batch_elements"`
}

type SpanBatchElement struct {
	EpochNum     int           `json:"EpochNum"`
	Timestamp    int64         `json:"Timestamp"`
	Transactions []interface{} `json:"Transactions"`
}

type BatchWithDelay struct {
	inner          SpanBatchElement
	DelaySec       int64
	DelayBlock     int
	InclusionBlock int
}

const (
	execPath     = "../op-node/cmd/batch_decoder"
	batchInbox   = "0x7e804b214944c5eC552dfD199f665cC001ABA460"
	batchSender  = "0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc"
	cacheDir     = "/tmp/batch_decoder/"
	txCache      = cacheDir + "transactions_cache"
	channelCache = cacheDir + "channel_cache"
)

// decodeBatch uses batch_decoder as a command.
// This could be improved by just reading it in golang, but that's enough right now.
func decodeBatch(start, end int) {
	// clear cache
	_, err := os.Stat(cacheDir)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error reading directory ", cacheDir, err)
	} else {
		os.RemoveAll(cacheDir)
		log.Println("Cache removed!")
	}

	_, err = os.Stat(execPath)
	if os.IsNotExist(err) {
		log.Fatal("Directory does not exist: ", execPath)
	}
	err = os.Chdir(execPath)
	if err != nil {
		log.Fatal(err)
	}
	runCommand("go", "run", ".", "fetch", "--start", fmt.Sprintf("%d", start), "--end", fmt.Sprintf("%d", end), "--inbox", batchInbox, "--sender", batchSender, "--l1", l1Url, "--l1.beacon", l1Beacon, "--out", txCache)
	runCommand("go", "run", ".", "reassemble-devnet", "--in", txCache, "--out", channelCache)

	// read all reassembled channels
	_, err = os.Stat(channelCache)
	if err != nil {
		log.Fatal("Error reading channel cache", err)
	}
	files, err := os.ReadDir(channelCache)
	if err != nil {
		log.Fatal(err)
	}

	var batches []BatchWithDelay

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(channelCache, file.Name())
			data, err := readJSONFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", file.Name(), err)
				continue
			}
			fmt.Printf("File: %s\n", file.Name())
			fmt.Printf("ID: %s\n", data.ID)
			fmt.Printf("Is Ready: %v\n", data.IsReady)
			fmt.Printf("Number of Frames: %d\n", len(data.Frames))
			if len(data.Frames) != 1 {
				panic("More than one frame")
			}
			frameTimestamp := data.Frames[0].Timestamp
			frameBlock := data.Frames[0].InclusionBlock

			fmt.Printf("Number of Batches: %d\n", len(data.Batches))
			for i, batch := range data.Batches {
				fmt.Printf(" Batch %d\n", i)
				fmt.Printf(" Number of Span Batch Elements: %d\n", len(batch.SpanBatchElements))
				for i, elem := range batch.SpanBatchElements {
					batches = append(batches, BatchWithDelay{
						inner:          elem,
						DelaySec:       frameTimestamp - elem.Timestamp,
						DelayBlock:     frameBlock - elem.EpochNum,
						InclusionBlock: frameBlock,
					})
					fmt.Printf("  Span Batch Element %d\n", i)
					fmt.Printf("  Timestamp: %d\n", elem.Timestamp)
					fmt.Printf("  Number of Transactions: %d\n", len(elem.Transactions))
				}
			}
			fmt.Println("----------------------------")
		}
	}

	// sort batches by timestamp
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].inner.Timestamp < batches[j].inner.Timestamp
	})
	for i := 1; i < len(batches); i++ {
		curr := batches[i].inner
		prev := batches[i-1].inner
		fmt.Printf("%+v; deltaSec: %d, \n", batches[i], curr.Timestamp-prev.Timestamp)
	}
}

func runCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}
	log.Print(string(out[:]))
}

func readJSONFile(filePath string) (*FileData, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data FileData
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
