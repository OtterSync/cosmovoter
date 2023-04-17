package cosmos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	BaseAPIUrlFormat = "https://cosmos.directory/%s"
)

func GetChainInfo(chain string) (ChainInfo, error) {
	var chainInfo ChainInfo
	url := fmt.Sprintf(BaseAPIUrlFormat, chain)

	resp, err := http.Get(url)
	if err != nil {
		return chainInfo, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return chainInfo, fmt.Errorf("error fetching chain info: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return chainInfo, err
	}

	err = json.Unmarshal(body, &chainInfo)
	if err != nil {
		return chainInfo, err
	}

	return chainInfo, nil
}

func GetOpenProposals(rpcNode string) ([]Proposal, error) {
	// Fetch open proposals using the Cosmos SDK API
	// Implement this function based on the specific API version and requirements
	return []Proposal{}, nil
}

func HasVoted(rpcNode, walletName string, proposalID int) (bool, error) {
	// Check if the wallet has voted on the proposal using the Cosmos SDK API
	// Implement this function based on the specific API version and requirements
	return false, nil
}

func SubmitVote(rpcNode, walletName, password string, proposalID int, voteOption, gasPrices string) (string, error) {
	// Submit the vote using the Cosmos SDK API
	// Implement this function based on the specific API version and requirements
	return "", nil
}
