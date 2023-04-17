package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/OtterSync/cosmovisor/internal/cosmos"
	"github.com/OtterSync/cosmovisor/internal/crypto"
)

func Vote() {
	configs, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	for _, config := range configs {
		proposals, err := cosmos.GetOpenProposals(config.RPCNode)
		if err != nil {
			fmt.Printf("Error getting open proposals for chain %s: %v\n", config.ChainName, err)
			continue
		}

		for _, proposal := range proposals {
			voted, err := cosmos.HasVoted(config.RPCNode, config.WalletName, proposal.ID)
			if err != nil {
				fmt.Printf("Error checking vote status for proposal %d on chain %s: %v\n", proposal.ID, config.ChainName, err)
				continue
			}

			if !voted {
				fmt.Printf("Proposal %d on chain %s:\n", proposal.ID, config.ChainName)
				fmt.Printf("Title: %s\n", proposal.Title)
				fmt.Printf("Description: %s\n", proposal.Description)

				var voteOption int
				for {
					fmt.Print("Enter your vote (1 - Yes, 2 - No, 3 - NoWithVeto, 4 - Abstain): ")
					voteOptionStr, _ := reader.ReadString('\n')
					voteOption, err = strconv.Atoi(strings.TrimSpace(voteOptionStr))
					if err == nil && voteOption >= 1 && voteOption <= 4 {
						break
					}
					fmt.Println("Invalid option. Please enter a valid vote option.")
				}

				voteOptionText := cosmos.VoteOptionText(voteOption)
				password, err := crypto.DecryptWithPassword(config.Password, config.Password, config.Salt)
				if err != nil {
					fmt.Printf("Error decrypting wallet password for chain %s: %v\n", config.ChainName, err)
					continue
				}

				txHash, err := cosmos.SubmitVote(config.RPCNode, config.WalletName, password, proposal.ID, voteOptionText, config.GasPrices)
				if err != nil {
					fmt.Printf("Error submitting vote for proposal %d on chain %s: %v\n", proposal.ID, config.ChainName, err)
					continue
				}

				fmt.Printf("Successfully submitted vote for proposal %d on chain %s. Transaction hash: %s\n", proposal.ID, config.ChainName, txHash)
			}
		}
	}
}
