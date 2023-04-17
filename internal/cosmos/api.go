package cosmos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func FetchProposals(chainConfig ChainConfig) ([]Proposal, error) {
	url := fmt.Sprintf("https://cosmos.directory/%s/proposals", chainConfig.ChainID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var proposals []Proposal
	if err := json.Unmarshal(body, &proposals); err != nil {
		return nil, err
	}

	return proposals, nil
}

func CheckWalletVote(chainConfig ChainConfig, walletAddress string, proposalID string) (Vote, error) {
	url := fmt.Sprintf("https://cosmos.directory/%s/proposals/%s/votes/%s", chainConfig.ChainID, proposalID, walletAddress)

	resp, err := http.Get(url)
	if err != nil {
		return Vote{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Vote{}, err
	}

	var vote Vote
	if err := json.Unmarshal(body, &vote); err != nil {
		return Vote{}, err
	}

	return vote, nil
}

func SubmitVote(chainConfig ChainConfig, walletAddress string, proposalID string, voteOption string) error {
	// Implement the function to submit a vote using the wallet address, proposal ID, and vote option
}
