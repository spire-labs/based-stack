package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func fetchValidatorLookahead(epoch string) (ProposerDutiesResponse, error) {
	// Fetch validator duties for the given epoch
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

	body, err := ioutil.ReadAll(resp.Body)
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
