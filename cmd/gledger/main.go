package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hirosato/gledger/adapters/inbound/cli/commands"
	"github.com/hirosato/gledger/adapters/outbound/filesystem"
	"github.com/hirosato/gledger/application"
)

const version = "0.1.0-alpha"

func main() {
	// Parse command-line flags
	var (
		fileFlag    = flag.String("f", "", "Read journal from FILE")
		fileFlagAlt = flag.String("file", "", "Read journal from FILE")
		helpFlag    = flag.Bool("h", false, "Display help")
		helpFlagAlt = flag.Bool("help", false, "Display help")
		versionFlag = flag.Bool("v", false, "Display version")
		versionFlagAlt = flag.Bool("version", false, "Display version")
	)

	flag.Parse()

	// Handle help and version flags
	if *helpFlag || *helpFlagAlt {
		printHelp()
		os.Exit(0)
	}

	if *versionFlag || *versionFlagAlt {
		fmt.Printf("gledger %s\n", version)
		fmt.Println("A Go implementation of ledger-cli")
		fmt.Println("Copyright (c) 2024")
		os.Exit(0)
	}

	// Get remaining args after flags
	args := flag.Args()
	
	if len(args) == 0 {
		fmt.Println("gledger: No command specified")
		fmt.Println("Run 'gledger --help' for usage information")
		os.Exit(1)
	}

	// Determine the journal file
	journalFile := *fileFlag
	if journalFile == "" {
		journalFile = *fileFlagAlt
	}
	
	// If no file specified, try to read from stdin or look for default
	var inputFile *os.File
	var err error
	
	if journalFile != "" {
		inputFile, err = os.Open(journalFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", journalFile, err)
			os.Exit(1)
		}
		defer inputFile.Close()
	} else {
		// For now, require a file to be specified
		// In the future, we could look for .ledger or read from stdin
		fmt.Fprintf(os.Stderr, "Error: No journal file specified. Use -f or --file option.\n")
		os.Exit(1)
	}

	// Create dependencies
	parser := filesystem.NewParserAdapter()
	
	// Create and load journal with injected dependencies
	journal := application.NewJournal(parser)
	if err := journal.LoadFromReader(inputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing journal: %v\n", err)
		os.Exit(1)
	}

	// Get command and remaining arguments
	command := strings.ToLower(args[0])
	commandArgs := args[1:]

	// Execute command
	switch command {
	case "accounts":
		cmd := commands.NewAccountsCommand(journal)
		if err := cmd.Execute(commandArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	
	case "balance", "bal":
		cmd := commands.NewBalanceCommand(journal)
		if err := cmd.Execute(commandArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	
	case "register", "reg", "r":
		cmd := commands.NewRegisterCommand(journal)
		if err := cmd.Execute(commandArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	
	case "print":
		fmt.Println("Print command not yet implemented")
		os.Exit(1)
	
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Println("Run 'gledger --help' for usage information")
		os.Exit(1)
	}
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