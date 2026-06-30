package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Expense struct {
	Date        string  `json:"date"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

func getDataFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "expenses.json"
	}
	dir := filepath.Join(home, ".poopy")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "expenses.json")
}

func saveExpenses(expenses []Expense, filename string) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
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
	if err := json.Unmarshal(data, &expenses); err != nil {
		return nil, err
	}
	return expenses, nil
}

func handleAdd() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	date := addCmd.String("d", "", "date of the expense (yyyy-mm-dd)")
	desc := addCmd.String("desc", "", "description of the expense")
	category := addCmd.String("c", "general", "category of the expense")

	// separate the amount (first non-flag arg) from flag args so both orderings work:
	// poopy add 25.50 -c food  AND  poopy add -c food 25.50
	var amountStr string
	var flagArgs []string
	for _, arg := range os.Args[2:] {
		if !strings.HasPrefix(arg, "-") && amountStr == "" {
			amountStr = arg
		} else {
			flagArgs = append(flagArgs, arg)
		}
	}
	addCmd.Parse(flagArgs)

	if amountStr == "" {
		fmt.Println("error: amount is required")
		fmt.Println("usage: poopy add <amount> [options]")
		return
	}

	var amount float64
	if _, err := fmt.Sscanf(amountStr, "%f", &amount); err != nil {
		fmt.Println("error: invalid amount")
		return
	}

	expenseDate := *date
	if expenseDate == "" {
		expenseDate = time.Now().Format("2006-01-02")
	}

	dataFile := getDataFile()
	expenses, err := loadExpenses(dataFile)
	if err != nil {
		fmt.Println("error loading expenses:", err)
		return
	}

	expenses = append(expenses, Expense{
		Date:        expenseDate,
		Amount:      amount,
		Description: *desc,
		Category:    *category,
	})

	if err := saveExpenses(expenses, dataFile); err != nil {
		fmt.Println("error saving expenses:", err)
		return
	}

	fmt.Printf("added expense: $%.2f on %s\n", amount, expenseDate)
	if *desc != "" {
		fmt.Printf("  description: %s\n", *desc)
	}
}

func handleList() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	startDate := listCmd.String("s", "", "start date (yyyy-mm-dd)")
	endDate := listCmd.String("e", "", "end date (yyyy-mm-dd)")
	category := listCmd.String("c", "", "filter by category")
	listCmd.Parse(os.Args[2:])

	dataFile := getDataFile()
	expenses, err := loadExpenses(dataFile)
	if err != nil {
		fmt.Println("error loading expenses:", err)
		return
	}

	if len(expenses) == 0 {
		fmt.Println("no expenses to list")
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
		fmt.Println("no expenses with the selected filters")
		return
	}

	displayExpenses(expenses)
}

func filterByStartDate(expenses []Expense, startDate string) []Expense {
	filtered := []Expense{}
	for _, e := range expenses {
		if e.Date >= startDate {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func filterByEndDate(expenses []Expense, endDate string) []Expense {
	filtered := []Expense{}
	for _, e := range expenses {
		if e.Date <= endDate {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func filterByCategory(expenses []Expense, category string) []Expense {
	filtered := []Expense{}
	for _, e := range expenses {
		if e.Category == category {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func displayExpenses(expenses []Expense) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("%-12s  %-10s  %-15s  %s\n", "date", "amount", "category", "description")
	fmt.Println(strings.Repeat("=", 80))

	total := 0.0
	for _, e := range expenses {
		desc := e.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		if desc == "" {
			desc = "-"
		}
		fmt.Printf("%-12s  $%-9.2f  %-15s  %s\n", e.Date, e.Amount, e.Category, desc)
		total += e.Amount
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-13s $%-9.2f\n", "total", total)
	fmt.Println(strings.Repeat("=", 80) + "\n")
}

func handleSummary() {
	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	month := summaryCmd.Int("m", 0, "month (1-12)")
	year := summaryCmd.Int("y", 0, "year (e.g., 2026)")
	summaryCmd.Parse(os.Args[2:])

	dataFile := getDataFile()
	expenses, err := loadExpenses(dataFile)
	if err != nil {
		fmt.Println("error loading expenses:", err)
		return
	}

	if len(expenses) == 0 {
		fmt.Println("no expenses to show")
		return
	}

	if *year > 0 {
		expenses = filterByYear(expenses, *year)
	}
	if *month > 0 {
		expenses = filterByMonth(expenses, *year, *month)
	}

	if len(expenses) == 0 {
		fmt.Println("no expenses to show")
		return
	}

	displaySummary(expenses, *month, *year)
}

func filterByYear(expenses []Expense, year int) []Expense {
	filtered := []Expense{}
	yearStr := fmt.Sprintf("%d", year)
	for _, e := range expenses {
		if len(e.Date) >= 4 && e.Date[:4] == yearStr {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func filterByMonth(expenses []Expense, year, month int) []Expense {
	filtered := []Expense{}
	monthStr := fmt.Sprintf("%d-%02d", year, month)
	for _, e := range expenses {
		if len(e.Date) >= 7 && e.Date[:7] == monthStr {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func displaySummary(expenses []Expense, month, year int) {
	categories := map[string]float64{}
	total := 0.0
	for _, e := range expenses {
		categories[e.Category] += e.Amount
		total += e.Amount
	}

	type catAmount struct {
		name   string
		amount float64
	}
	cats := make([]catAmount, 0, len(categories))
	for name, amount := range categories {
		cats = append(cats, catAmount{name, amount})
	}
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].amount > cats[j].amount
	})

	var period string
	switch {
	case month > 0 && year > 0:
		period = fmt.Sprintf("%d-%02d", year, month)
	case year > 0:
		period = fmt.Sprintf("%d", year)
	default:
		period = "all time"
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("expense summary - %s\n", period)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("%-25s  %s\n", "category", "amount")
	fmt.Println(strings.Repeat("-", 50))
	for _, cat := range cats {
		fmt.Printf("%-25s  $%.2f\n", cat.name, cat.amount)
	}
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("%-25s  $%.2f\n", "total", total)
	fmt.Println(strings.Repeat("=", 50) + "\n")
}

func handleDeleteLast() {
	dataFile := getDataFile()
	expenses, err := loadExpenses(dataFile)
	if err != nil {
		fmt.Println("error loading expenses:", err)
		return
	}

	if len(expenses) == 0 {
		fmt.Println("no expenses to delete")
		return
	}

	deleted := expenses[len(expenses)-1]
	expenses = expenses[:len(expenses)-1]

	if err := saveExpenses(expenses, dataFile); err != nil {
		fmt.Println("error saving expenses:", err)
		return
	}

	fmt.Printf("deleted expense: $%.2f on %s\n", deleted.Amount, deleted.Date)
}

func printUsage() {
	fmt.Println("poopy - expense tracker")
	fmt.Println()
	fmt.Println("usage:")
	fmt.Println("  poopy add <amount> [options]   add an expense")
	fmt.Println("  poopy list [options]            list expenses")
	fmt.Println("  poopy summary [options]         show summary")
	fmt.Println("  poopy delete-last               delete last expense")
	fmt.Println()
	fmt.Println("run 'poopy <command> -help' for command options")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "add":
		handleAdd()
	case "list":
		handleList()
	case "summary":
		handleSummary()
	case "delete-last":
		handleDeleteLast()
	case "-v", "--version":
		fmt.Println("Poopy 0.0.1")
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		printUsage()
	}
}
