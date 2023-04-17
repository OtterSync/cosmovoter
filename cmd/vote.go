package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/OtterSync/cosmovoter/internal/config"
	"github.com/OtterSync/cosmovoter/internal/cosmos"
	"github.com/OtterSync/cosmovoter/internal/crypto"
)

func Vote(chain string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	if chain != "" {
		voteOnChain(cfg, chain)
	} else {
		for _, chain := range cfg.Chains {
			voteOnChain(cfg, chain.Name)
		}
	}
}

func voteOnChain(cfg *config.Config, chainName string) {
	chainCfg, err := cfg.GetChain(chainName)
	if err != nil {
		fmt.Printf("Error getting chain '%s': %v\n", chainName, err)
		return
	}

	chainInfo, err := cosmos.GetChainInfo(chainName, chainCfg.RPCNode)
	if err != nil {
		fmt.Printf("Error fetching chain info for '%s': %v\n", chainName, err)
		return
	}

	proposals, err := cosmos.GetOpenProposals(chainInfo.RPCNode)
	if err != nil {
		fmt.Printf("Error getting open proposals for '%s': %v\n", chainName, err)
		return
	}

	for _, proposal := range proposals {
		voted, err := cosmos.HasVoted(chainInfo.RPCNode, chainCfg.WalletName, proposal.ProposalId)
		if err != nil {
			fmt.Printf("Error checking if wallet has voted on proposal %d for '%s': %v\n", proposal.ProposalId, chainName, err)
			continue
		}

		if !voted {
			fmt.Printf("Proposal %d on chain '%s':\n", proposal.ProposalId, chainName)
			fmt.Println(proposal.Content.GetDescription())
			voteOption := getVoteOption()
			walletPassword, err := crypto.Decrypt(chainCfg.WalletPassword, cfg.PasswordSalt)
			if err != nil {
				fmt.Printf("Error decrypting wallet password for '%s': %v\n", chainName, err)
				continue
			}

			txHash, err := cosmos.SubmitVote(chainInfo.RPCNode, chainCfg.WalletName, walletPassword, proposal.ProposalId, voteOption, chainCfg.GasPrices)
			if err != nil {
				fmt.Printf("Error submitting vote on proposal %d for '%s': %v\n", proposal.ProposalId, chainName, err)
				continue
			}

			fmt.Printf("Successfully submitted vote on proposal %d for '%s'. Tx hash: %s\n", proposal.ProposalId, chainName, txHash)
		}
	}
}
