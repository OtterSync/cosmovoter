package main

import (
	"fmt"
	"os"

	"github.com/OtterSync/cosmovoter/cmd"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println(version)
	case "config":
		handleConfigCommand()
	case "vote":
		cmd.Vote()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage:
	cosmovoter <command>

Available commands:
	version
	config
	vote`)
}

func handleConfigCommand() {
	if len(os.Args) < 3 {
		printConfigUsage()
		os.Exit(1)
	}

	switch os.Args[2] {
	case "list":
		if err := cmd.ListChains(); err != nil {
			fmt.Printf("Error listing chains: %v\n", err)
			os.Exit(1)
		}
	case "add":
		if err := cmd.AddChain(); err != nil {
			fmt.Printf("Error adding chain: %v\n", err)
			os.Exit(1)
		}
	case "remove":
		if err := cmd.RemoveChain(); err != nil {
			fmt.Printf("Error removing chain: %v\n", err)
			os.Exit(1)
		}
	default:
		printConfigUsage()
		os.Exit(1)
	}
}

func printConfigUsage() {
	fmt.Println(`Usage:
	cosmovoter config <command>

Available commands:
	list
	add
	remove`)
}
