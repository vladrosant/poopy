#!/usr/bin/env python3
"""
Poopy - A simple expense tracker CLI
"""

import argparse
import json
import os
from datetime import datetime
from pathlib import Path

# Data file location
DATA_DIR = Path.home() / '.poopy'
DATA_FILE = DATA_DIR / 'expenses.json'

def init_data_file():
    """Initialize data directory and file if they don't exist"""
    DATA_DIR.mkdir(exist_ok=True)
    if not DATA_FILE.exists():
        with open(DATA_FILE, 'w') as f:
            json.dump([], f)

def load_expenses():
    """Load expenses from JSON file"""
    init_data_file()
    with open(DATA_FILE, 'r') as f:
        return json.load(f)

def save_expenses(expenses):
    """Save expenses to JSON file"""
    with open(DATA_FILE, 'w') as f:
        json.dump(expenses, f, indent=2)

def add_expense(amount, date=None, description=None, category=None):
    """Add a new expense"""
    if date is None:
        date = datetime.now().strftime('%Y-%m-%d')
    else:
        # Validate date format
        try:
            datetime.strptime(date, '%Y-%m-%d')
        except ValueError:
            print("Error: Date must be in YYYY-MM-DD format")
            return
    
    expenses = load_expenses()
    expense = {
        'date': date,
        'amount': float(amount),
        'description': description or '',
        'category': category or 'general'
    }
    expenses.append(expense)
    save_expenses(expenses)
    print(f"✓ Added expense: ${amount} on {date}")
    if description:
        print(f"  Description: {description}")

def list_expenses(start_date=None, end_date=None, category=None):
    """List expenses within a date range"""
    expenses = load_expenses()
    
    if not expenses:
        print("No expenses recorded yet.")
        return
    
    # Filter by date range
    if start_date:
        expenses = [e for e in expenses if e['date'] >= start_date]
    if end_date:
        expenses = [e for e in expenses if e['date'] <= end_date]
    
    # Filter by category
    if category:
        expenses = [e for e in expenses if e['category'] == category]
    
    if not expenses:
        print("No expenses found for the specified criteria.")
        return
    
    # Sort by date
    expenses.sort(key=lambda x: x['date'])
    
    # Display expenses
    total = 0
    print("\n" + "="*70)
    print(f"{'Date':<12} {'Amount':>10} {'Category':<15} {'Description'}")
    print("="*70)
    
    for expense in expenses:
        total += expense['amount']
        desc = expense['description'][:30] if expense['description'] else '-'
        print(f"{expense['date']:<12} ${expense['amount']:>9.2f} {expense['category']:<15} {desc}")
    
    print("="*70)
    print(f"{'TOTAL':>22} ${total:>9.2f}")
    print("="*70 + "\n")

def summary(month=None, year=None):
    """Show summary of expenses"""
    expenses = load_expenses()
    
    if not expenses:
        print("No expenses recorded yet.")
        return
    
    # Filter by month/year
    if year:
        expenses = [e for e in expenses if e['date'].startswith(str(year))]
    if month and year:
        month_str = f"{year}-{month:02d}"
        expenses = [e for e in expenses if e['date'].startswith(month_str)]
    
    if not expenses:
        print("No expenses found for the specified period.")
        return
    
    # Calculate totals by category
    categories = {}
    total = 0
    
    for expense in expenses:
        cat = expense['category']
        amount = expense['amount']
        categories[cat] = categories.get(cat, 0) + amount
        total += amount
    
    # Display summary
    period = f"{year}-{month:02d}" if month and year else str(year) if year else "All time"
    print(f"\n{'='*50}")
    print(f"Expense Summary - {period}")
    print(f"{'='*50}")
    print(f"{'Category':<25} {'Amount':>15}")
    print(f"{'-'*50}")
    
    for cat, amount in sorted(categories.items(), key=lambda x: x[1], reverse=True):
        print(f"{cat:<25} ${amount:>14.2f}")
    
    print(f"{'-'*50}")
    print(f"{'TOTAL':<25} ${total:>14.2f}")
    print(f"{'='*50}\n")

def delete_last():
    """Delete the last expense added"""
    expenses = load_expenses()
    if not expenses:
        print("No expenses to delete.")
        return
    
    deleted = expenses.pop()
    save_expenses(expenses)
    print(f"✓ Deleted expense: ${deleted['amount']} on {deleted['date']}")

def main():
    parser = argparse.ArgumentParser(
        description='Poopy - Track your expenses',
        epilog='Examples:\n'
               '  poopy add 50.00\n'
               '  poopy add 30.50 -d 2025-01-15 -desc "Groceries" -c food\n'
               '  poopy list -s 2025-01-01 -e 2025-01-15\n'
               '  poopy summary -m 1 -y 2025\n',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    subparsers = parser.add_subparsers(dest='command', help='Available commands')
    
    # Add expense
    add_parser = subparsers.add_parser('add', help='Add a new expense')
    add_parser.add_argument('amount', type=float, help='Amount spent')
    add_parser.add_argument('-d', '--date', help='Date (YYYY-MM-DD, default: today)')
    add_parser.add_argument('-desc', '--description', help='Description of expense')
    add_parser.add_argument('-c', '--category', help='Category (e.g., food, transport)')
    
    # List expenses
    list_parser = subparsers.add_parser('list', help='List expenses')
    list_parser.add_argument('-s', '--start', help='Start date (YYYY-MM-DD)')
    list_parser.add_argument('-e', '--end', help='End date (YYYY-MM-DD)')
    list_parser.add_argument('-c', '--category', help='Filter by category')
    
    # Summary
    summary_parser = subparsers.add_parser('summary', help='Show expense summary')
    summary_parser.add_argument('-m', '--month', type=int, help='Month (1-12)')
    summary_parser.add_argument('-y', '--year', type=int, help='Year (e.g., 2025)')
    
    # Delete last
    subparsers.add_parser('delete-last', help='Delete the last expense added')
    
    args = parser.parse_args()
    
    if args.command == 'add':
        add_expense(args.amount, args.date, args.description, args.category)
    elif args.command == 'list':
        list_expenses(args.start, args.end, args.category)
    elif args.command == 'summary':
        summary(args.month, args.year)
    elif args.command == 'delete-last':
        delete_last()
    else:
        parser.print_help()

if __name__ == '__main__':
    main()