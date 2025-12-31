package main

import (
	"fmt"
)

type Expense struct {
	Date        string
	Amount      float64
	Description string
	Category    string
}

func main() {
	expenses := []Expense{}

	expenses = append(expenses, Expense{
		Date:        "2025-01-15",
		Amount:      49.99,
		Description: "Groceries",
		Category:    "food",
	})

	expenses = append(expenses, Expense{
		Date:        "2025-01-16",
		Amount:      25.50,
		Description: "Gas",
		Category:    "transport",
	})

	expenses = append(expenses, Expense{
		Date:        "2025-12-31",
		Amount:      30.22,
		Description: "Meal",
		Category:    "food",
	})

	expenses = append(expenses, Expense{
		Date:        "2025-12-31",
		Amount:      5.0,
		Description: "Meal",
		Category:    "food",
	})

	expenses = append(expenses, Expense{
		Date:        "2026-01-01",
		Amount:      122.0,
		Description: "Shopping",
		Category:    "shopping",
	})

	fmt.Println("========= FOOD EXPENSES ========")
	for _, expense := range expenses {
		if expense.Category == "food" {
			fmt.Printf("%s, $%.2f - %s\n",
				expense.Date,
				expense.Amount,
				expense.Description)
		}
	}

	total := 0.0
	for _, expense := range expenses {
		total += expense.Amount
	}
	average := total / float64(len(expenses))

	fmt.Println("\n======== EXPENSES ========")
	for _, expense := range expenses {
		fmt.Printf("%s, $%.2f - %s\n",
			expense.Date,
			expense.Amount,
			expense.Description)
	}
	fmt.Printf("\nTotal spent: $%.2f\nAverage spent: $%.2f\n", total, average)
}
