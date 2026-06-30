# poopy - a cli expense tracker

minimal personal expense tracker built in go

## requirements
- go 1.21+

## building

```bash
go build -o poopy .
```

## usage

```bash
# add an expense
poopy add 25.50
poopy add 25.50 -desc "lunch" -c food
poopy add 25.50 -d 2026-01-15 -desc "groceries" -c food

# list expenses
poopy list
poopy list -c food
poopy list -s 2026-01-01 -e 2026-01-31

# summary
poopy summary
poopy summary -y 2026
poopy summary -m 6 -y 2026

# delete last expense added
poopy delete-last
```

## flags

| command | flag | description |
|---|---|---|
| add | `<amount>` | amount spent (required) |
| add | `-d` | date in yyyy-mm-dd (default: today) |
| add | `-desc` | description |
| add | `-c` | category (default: general) |
| list | `-s` | start date (yyyy-mm-dd) |
| list | `-e` | end date (yyyy-mm-dd) |
| list | `-c` | filter by category |
| summary | `-m` | month (1-12) |
| summary | `-y` | year |

## data

expenses are stored at `~/.poopy/expenses.json`
