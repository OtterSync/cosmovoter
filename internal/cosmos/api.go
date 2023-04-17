package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	"github.com/cosmos/cosmos-sdk/client/rpc"
)

const (
	BaseAPIUrlFormat = "https://cosmos.directory/%s"
)

type ChainInfo struct {
	ChainID string `json:"chain_id"`
	RPCNode string
}

func GetChainInfo(chain, rpcNode string) (ChainInfo, error) {
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

	chainInfo.RPCNode = rpcNode
	return chainInfo, nil
}

func GetOpenProposals(rpcNode string) ([]types.Proposal, error) {
	client, err := rpc.New(rpcNode, "/websocket")
	if err != nil {
		return nil, err
	}

	res, _, err := client.QueryWithData(context.Background(), "/custom/gov/proposals", query.NewPageRequest())
	if err != nil {
		return nil, err
	}

	var proposalsResponse types.QueryProposalsResponse
	if err := client.Codec.UnmarshalJSON(res, &proposalsResponse); err != nil {
		return nil, err
	}

	return proposalsResponse.GetProposals(), nil
}

func HasVoted(rpcNode, walletName string, proposalID uint64) (bool, error) {
	client, err := rpc.New(rpcNode, "/websocket")
	if err != nil {
		return false, err
	}

	params := types.NewQueryVoteRequest(walletName, proposalID)
	req, err := client.Codec.MarshalJSON(params)
	if err != nil {
		return false, err
	}

	_, _, err = client.QueryWithData(context.Background(), "/custom/gov/vote", req)
	if err != nil {
		if err.Error() == "rpc error: code = NotFound desc = vote not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
func SubmitVote(rpcNode, walletName, password string, proposalID uint64, voteOption, gasPrices string) (string, error) {
	client, err := rpc.New(rpcNode, "/websocket")
	if err != nil {
		return "", err
	}

	voteOptionParsed, err := types.VoteOptionFromString(voteOption)
	if err != nil {
		return "", err
	}

	msg := types.NewMsgVote(walletName, proposalID, voteOptionParsed)
	tx, err := cli.BuildUnsignedTx(msg, walletName, password, gasPrices)
	if err != nil {
		return "", err
	}

	res, err := client.BroadcastTxCommit(context.Background(), tx)
	if err != nil {
		return "", err
	}

	return res.Hash.String(), nil
}

