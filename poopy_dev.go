package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
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
		expenseDate = time.Now().Format("2006-01-02")
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

	fmt.Printf("✓ Added expense: $%.2f on %s\n", amount, expenseDate)
	if *desc != "" {
		fmt.Printf("	Description: %s\n", *desc)
	}
}

func handleList() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	startDate := listCmd.String("s", "", "Start date (YYYY-MM-DD)")
	endDate := listCmd.String("e", "", "End date (YYYY-MM-DD")
	category := listCmd.String("c", "", "Filter by category")

	listCmd.Parse(os.Args[2:])

	expenses, err := loadExpenses(dataFile)
	if err != nil {
		fmt.Println("Error loading expenses:", err)
		return
	}

	if len(expenses) == 0 {
		fmt.Println("No expenses to list")
		return
	}

	if *startDate != "" {
		expenses = filterByStartDate(expenses, *startDate)
	}

	if *endDate != "" {
		expenses = filterByEndDate(expenses, *endDate)
	}

	if *category != "" {
		expenses = filterByCategory(expenses, *category)
	}

	if len(expenses) == 0 {
		fmt.Println("No expenses with the selected filters")
		return
	}

	displayExpenses(expenses)
}

func filterByStartDate(expenses []Expense, startDate string) []Expense {
	filtered := []Expense{}

	for _, expense := range expenses {
		if expense.Date >= startDate {
			filtered = append(filtered, expense)
		}
	}
	return filtered
}

func filterByEndDate(expenses []Expense, endDate string) []Expense {
	filtered := []Expense{}

	for _, expense := range expenses {
		if expense.Date <= endDate {
			filtered = append(filtered, expense)
		}
	}
	return filtered
}

func filterByCategory(expenses []Expense, category string) []Expense {
	filtered := []Expense{}

	for _, expense := range expenses {
		if expense.Category == category {
			filtered = append(filtered, expense)
		}
	}

	return filtered
}

func displayExpenses(expenses []Expense) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("%-12s  %-10s  %-15s  %s\n", "Date", "Amount", "Category", "Description")
	fmt.Println(strings.Repeat("=", 80))

	total := 0.0
	for _, expense := range expenses {
		desc := expense.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..." // Truncate long descriptions
		}
		if desc == "" {
			desc = "-" // Show dash for empty descriptions
		}

		fmt.Printf("%-12s  $%-9.2f  %-15s  %s\n",
			expense.Date,
			expense.Amount,
			expense.Category,
			desc)

		total += expense.Amount
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-24s $%-9.2f\n", "TOTAL", total)
	fmt.Println(strings.Repeat("=", 80) + "\n")
}

func handleSummary() {
	fmt.Println("Coming soon")
}

func handleDeleteLast() {
	expenses, err := loadExpenses(dataFile)
	if err != nil {
		fmt.Println("Error loading expenses:", err)
		return
	}

	if len(expenses) == 0 {
		fmt.Println("No expenses to delete")
		return
	}

	deleted := expenses[len(expenses)-1]
	expenses = expenses[:len(expenses)-1]

	err = saveExpenses(expenses, dataFile)
	if err != nil {
		fmt.Println("Error saving expenses", err)
		return
	}

	fmt.Printf("✓ Deleted expense: %.2f on %s", deleted.Amount, deleted.Date)
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
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		return
	}
}
