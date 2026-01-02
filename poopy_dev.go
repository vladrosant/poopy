package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Expense struct {
	Date        string  `json:"date"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

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

func filterByCategory(expenses []Expense, category string) []Expense {
	filteredExpenses := []Expense{}
	for _, expense := range expenses {
		if expense.Category == category {
			filteredExpenses = append(filteredExpenses, expense)
		}
	}
	return filteredExpenses
}

func calculateAverage(expenses []Expense) float64 {
	if len(expenses) == 0 {
		return 0.0
	}

	total := 0.0
	for _, expense := range expenses {
		total += expense.Amount
	}
	return total / float64(len(expenses))
}

func main() {
	filename := "expenses.json"

	expenses, err := loadExpenses(filename)
	if err != nil {
		fmt.Println("Error loading expenses:", err)
		return
	}
	fmt.Printf("Loaded %d expenses\n", len(expenses))

	newExpense := Expense{
		Date:        "2026-01-01",
		Amount:      100.00,
		Description: "Groceries",
		Category:    "food",
	}
	expenses = append(expenses, newExpense)

	err = saveExpenses(expenses, filename)
	if err != nil {
		fmt.Println("Error saving expenses:", err)
		return
	}

	fmt.Printf("Saved %d expenses\n", len(expenses))

	average := calculateAverage(filterByCategory(expenses, "food"))
	fmt.Printf("Average food expense: $%.2f\n", average)
}
