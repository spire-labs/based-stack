package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type ProposerDuty struct {
	Pubkey         string `json:"pubkey"`
	ValidatorIndex string `json:"validator_index"`
	Slot           string `json:"slot"`
}

type ProposerDutiesResponse struct {
	DependentRoot       string         `json:"dependent_root"`
	ExecutionOptimistic bool           `json:"execution_optimistic"`
	Data                []ProposerDuty `json:"data"`
}

type BeaconHeadResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		Root      string `json:"root"`
		Canonical bool   `json:"canonical"`
		Header    struct {
			Message struct {
				Slot          string `json:"slot"`
				ProposerIndex string `json:"proposer_index"`
				ParentRoot    string `json:"parent_root"`
				StateRoot     string `json:"state_root"`
				BodyRoot      string `json:"body_root"`
			} `json:"message"`
			Signature string `json:"signature"`
		} `json:"header"`
	} `json:"data"`
}

// fetch the latest slot and convert to an epoch
func getLatestEpoch() (string, error) {
	apiKey := ""
	beaconUrl := "https://ethereum-mainnet.core.chainstack.com/beacon/" + apiKey + "/eth/v1/beacon/headers/head"

	resp, err := http.Get(beaconUrl)
	if err != nil {
		return "", fmt.Errorf("error making the request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var headResp BeaconHeadResponse
	err = json.Unmarshal(body, &headResp)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	slot, err := strconv.ParseUint(headResp.Data.Header.Message.Slot, 10, 64)
	if err != nil {
		return "", fmt.Errorf("error parsing slot: %w", err)
	}

	// Calculate the head epoch by dividing the slot by 32, and then get the next epoch
	epoch := slot/32 + 1

	return strconv.FormatUint(epoch, 10), nil
}

// Fetch validator duties for the given epoch
func fetchValidatorLookahead(epoch string) (ProposerDutiesResponse, error) {
	apiKey := ""
	beaconUrl := "https://ethereum-mainnet.core.chainstack.com/beacon/" + apiKey + "/eth/v1/validator/duties/proposer/" + epoch

	req, err := http.NewRequest("GET", beaconUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ProposerDutiesResponse{}, fmt.Errorf("error creating beacon request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ProposerDutiesResponse{}, fmt.Errorf("error sending beacon request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProposerDutiesResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ProposerDutiesResponse{}, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var duties ProposerDutiesResponse
	err = json.Unmarshal(body, &duties)
	if err != nil {
		return ProposerDutiesResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	return duties, nil
}

func fetchNextLookahead() (ProposerDutiesResponse, error) {
	epoch, err := getLatestEpoch()
	if err != nil {
		return ProposerDutiesResponse{}, fmt.Errorf("error getting latest epoch: %w", err)
	}

	return fetchValidatorLookahead(epoch)
}
