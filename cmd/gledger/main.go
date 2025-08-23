package main

import (
	"fmt"
	"os"
)

const version = "0.1.0-alpha"

func main() {
	// TODO: Initialize application
	// - Set up dependency injection
	// - Parse command-line arguments
	// - Route to appropriate command handler
	// - Execute command
	
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("gledger %s\n", version)
		fmt.Println("A Go implementation of ledger-cli")
		fmt.Println("Copyright (c) 2024")
		os.Exit(0)
	}
	
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printHelp()
		os.Exit(0)
	}
	
	fmt.Println("gledger: Work in progress")
	fmt.Println("Run 'gledger --help' for usage information")
}

func printHelp() {
	fmt.Println("gledger - A Go implementation of ledger-cli")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gledger [command] [options] [file...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  balance, bal      Show account balances")
	fmt.Println("  register, reg     Show transaction register")
	fmt.Println("  print             Print transactions")
	fmt.Println("  accounts          List all accounts")
	fmt.Println("  payees            List all payees")
	fmt.Println("  commodities       List all commodities")
	fmt.Println("  stats             Show journal statistics")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -f, --file FILE   Read journal from FILE")
	fmt.Println("  -h, --help        Display this help")
	fmt.Println("  -v, --version     Display version information")
	fmt.Println()
	fmt.Println("For more information, see: https://github.com/hirosato/gledger")
}