package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	BaseAPIUrlFormat = "https://cosmos.directory/%s"
)

type ChainInfo struct {
	ChainID string `json:"chain_id"`
	RPCNode string `json:"rpc_node"`
}

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

func GetOpenProposals(rpcNode string) ([]types.Proposal, error) {
	client, err := rpcclient.New(rpcNode, "/websocket")
	if err != nil {
		return nil, err
	}

	proposals, err := client.GovProposals(context.Background(), &types.QueryProposalsRequest{
		ProposalStatus: types.StatusVotingPeriod,
	})

	if err != nil {
		return nil, err
	}

	return proposals.Proposals, nil
}

func HasVoted(rpcNode, walletName string, proposalID uint64) (bool, error) {
	client, err := rpcclient.New(rpcNode, "/websocket")
	if err != nil {
		return false, err
	}

	_, err = client.GovVote(context.Background(), &types.QueryVoteRequest{
		ProposalId: proposalID,
		Voter:      walletName,
	})

	if err != nil {
		if err.Error() == "rpc error: code = NotFound desc = vote not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func SubmitVote(rpcNode, walletName, password string, proposalID uint64, voteOption, gasPrices string) (string, error) {
	client, err := rpcclient.New(rpcNode, "/websocket")
	if err != nil {
		return "", err
	}

	voteOptionParsed, err := types.VoteOptionFromString(voteOption)
	if err != nil {
		return "", err
	}

	msg := types.NewMsgVote(walletName, proposalID, voteOptionParsed)
	txBuilder, err := client.BuildTx(context.Background(), walletName, password, []sdk.Msg{msg}, gasPrices)
	if err != nil {
		return "", err
	}

	res, err := client.BroadcastTx(context.Background(), txBuilder)
	if err != nil {
		return "", err
	}

	return res.TxHash, nil
}
