package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Expense struct {
	Date        string  `json:"date"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

const dataFile = "expenses.json"

func saveExpenses(expenses []Expense, filename string) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadExpenses(filename string) ([]Expense, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Expense{}, nil
		}
		return nil, err
	}
	var expenses []Expense
	err = json.Unmarshal(data, &expenses)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func handleAdd() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	date := addCmd.String("d", "", "Date of the expense (YYYY-MM-DD)")
	desc := addCmd.String("desc", "", "Description of the expense")
	category := addCmd.String("c", "general", "Category of the expense")

	addCmd.Parse(os.Args[2:])

	if addCmd.NArg() < 1 {
		fmt.Println("Error: amount is required")
		fmt.Println("Usage: poopy add <amount> [options]")
		return
	}

	var amount float64
	_, err := fmt.Sscanf(addCmd.Arg(0), "%f", &amount)
	if err != nil {
		fmt.Println("Error: invalid amount")
		return
	}

	expenseDate := *date
	if expenseDate == "" {
		expenseDate = time.Now().Format("02/01/2006")
	}

	expense := Expense{
		Date:        expenseDate,
		Amount:      amount,
		Description: *desc,
		Category:    *category,
	}

	expenses, err := loadExpenses(dataFile)
	if err != nil {
		fmt.Println("Error loading expenses:", err)
		return
	}

	expenses = append(expenses, expense)

	err = saveExpenses(expenses, dataFile)
	if err != nil {
		fmt.Println("Error saving expenses:", err)
		return
	}

	fmt.Printf("✓ Added expense: %.2f on %s\n", amount, expenseDate)
	if *desc != "" {
		fmt.Printf("	Description: %s\n", *desc)
	}
}

func handleList() {
	fmt.Println("Coming soon")
}

func handleSummary() {
	fmt.Println("Coming soon")
}

func handleDeleteLast() {
	fmt.Println("Coming soon")
}

func printUsage() {
	fmt.Println("Poopy Expense Tracker")
	fmt.Println()
	fmt.Println("Usage: ")
	fmt.Println("	poopy add <amount> [options]	Add an expense")
	fmt.Println("	poopy list [options]			List expenses")
	fmt.Println("	poopy summary [options]			Show summary")
	fmt.Println()
	fmt.Println("Run 'poopy <command> -help' for specific command options")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd()
	case "list":
		handleList()
	case "summary":
		handleSummary()
	case "delete-last":
		handleDeleteLast()
	default:
		fmt.Println("Invalid command")
		return
	}
}
