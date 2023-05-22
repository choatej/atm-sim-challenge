package main

//goland:noinspection SpellCheckingInspection
import (
	"agile-coder.com/atm-sim/cmd"
	"agile-coder.com/atm-sim/internal"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/c-bata/go-prompt"
	cobraprompt "github.com/stromland/cobra-prompt"
	"os"
	"strconv"
	"strings"
	"time"
)

var appPrompt = &cobraprompt.CobraPrompt{
	RootCmd:                  cmd.RootCmd,
	PersistFlagValues:        false,
	ShowHelpCommandAndFlags:  false,
	DisableCompletionCommand: false,
	AddDefaultExitCommand:    false,
	GoPromptOptions: []prompt.Option{
		prompt.OptionPrefix(">> "),
		prompt.OptionMaxSuggestion(0),
	},
	OnErrorFunc: func(err error) {
		if strings.Contains(err.Error(), "unknown command") {
			return
		}
	},
}

func main() {

	// initialize the application
	initLogger()

	// this could be injected from a config file or somewhere external
	startingCashInMachine := 10000.00
	initData(startingCashInMachine)

	// monitor session timeouts
	go func() {
		ticker := time.NewTicker(1 * time.Minute) // Adjust the ticker interval as needed
		for range ticker.C {
			internal.Logger.Println("tick...")
			session := internal.GetSession()
			// Check for session expiration if the user is authenticated
			if session.IsAuthenticated && time.Since(session.LastActivityTime) > 2*time.Minute {
				internal.Logger.Printf("session expired for %s", session.AccountId)
				session.IsAuthenticated = false
				session.AccountId = ""
				fmt.Println("Session expired due to inactivity.")
			}
		}
	}()

	cmd.RootCmd.SetHelpTemplate(`Available Commands:
{{- range $index, $command := .Commands}}
	{{printf "%-15s" $command.Name}}{{.Short}}{{end}}
`)

	// start the prompt
	fmt.Println("Welcome to the ATM simulator. Enter 'help' for available commands.")
	appPrompt.Run()
}

func initLogger() {
	// Initialize the logger
	internal.InitLogger("logfile.log", false)
	internal.Logger.Println("logging started")
}

func initData(startingCash float64) {
	internal.Logger.Println("reading in account data")
	filePath := "data/accounts.csv"

	// Open the CSV file
	file, err := Asset(filePath)
	if err != nil {
		internal.Logger.Printf("Error opening file: %+v", err)
		fmt.Println("Error opening file:", err)
		os.Exit(-1)
	}

	// Create a new CSV reader
	reader := csv.NewReader(bytes.NewReader(file))
	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		internal.Logger.Printf("Error reading CSV: %+v\n", err)
		os.Exit(-1)
	}

	internal.Logger.Printf("found %d records\n", len(records)-1)

	// map te field names to their locations in the array
	fieldIndexes := make(map[string]int)
	for i, field := range records[0] {
		fieldIndexes[field] = i
	}

	// data structures to hold the csv data
	authAccounts := map[string]internal.EncryptedPin{}
	ledgerAccounts := map[string]float64{}

	// Parse and process each CSV record starting from the second row (data rows)
	for i, record := range records[1:] {
		internal.Logger.Printf("reading record %d\n", i+1)
		// Ensure the record has the expected number of fields
		if len(record) != len(fieldIndexes) {
			internal.Logger.Println("Invalid record:", record)
			continue
		}

		// Parse the values from the record using the field indexes
		val, ok := fieldIndexes["ACCOUNT_ID"]
		if !ok {
			internal.Logger.Println("column index missing for ACCOUNT_ID")
			os.Exit(-1)
		}
		accountNumber := record[val]

		val, ok = fieldIndexes["PIN"]
		if !ok {
			internal.Logger.Println("column index missing for PIN")
			os.Exit(-1)
		}
		pin := record[val]

		val, ok = fieldIndexes["BALANCE"]
		if !ok {
			internal.Logger.Println("column index missing for BALANCE")
			os.Exit(-1)
		}
		balance, err := strconv.ParseFloat(record[val], 64)
		if err != nil {
			internal.Logger.Println("Error parsing balance:", err)
			continue
		}
		enc, err := internal.EncryptPin(pin)
		if err != nil {
			fmt.Printf("error reading input data\n")
			os.Exit(-1)
		}
		internal.Logger.Printf("read record %s, %s, %.2f\n", accountNumber, pin, balance)
		authAccounts[accountNumber] = enc
		ledgerAccounts[accountNumber] = balance
	}
	auth := internal.GetAuthorizationService()
	auth.SetAuthData(authAccounts)

	ledger := internal.GetLedgerService()
	ledger.SetInitialBalances(startingCash, ledgerAccounts)
}
