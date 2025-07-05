package main

import (
	"fmt"
	"sync"
)

const (
	MaxTransactionAmount = 2000.0
)

type BankAccount struct {
	ID         string
	Owner      string
	Balance    float64
	MinBalance float64
	mu         sync.Mutex
}

type AccountError struct {
	Message string
}

type InsufficientFundsError struct {
	MinBalance     float64
	InitialBalance float64
}

type NegativeAmountError struct {
	Amount float64
}

type ExceedsLimitError struct {
	MaxBalance   float64
	FinalBalance float64
}

func (e *AccountError) Error() string {
	return fmt.Sprintf("Account error: %s", e.Message)
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("Expected MinBalance %.2f but got %.2f", e.InitialBalance, e.MinBalance)
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Negative amount not allowed: %.2f", e.Amount)
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf("Transaction would go below maxmum balance. Final: %.2f, min: %2.f", e.FinalBalance, e.MaxBalance)
}

func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" || owner == "" {
		return nil, &AccountError{Message: "ID and or owner cannot be empty"}
	}

	if initialBalance < 0 {
		return nil, &NegativeAmountError{initialBalance}
	}

	if minBalance < 0 {
		return nil, &NegativeAmountError{minBalance}
	}

	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{InitialBalance: initialBalance, MinBalance: minBalance}
	}

	return &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}, nil
}

func (a *BankAccount) Deposit(amount float64) error {

	if amount < 0 {
		return &NegativeAmountError{amount}
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance+amount > MaxTransactionAmount {
		return &ExceedsLimitError{MaxBalance: MaxTransactionAmount, FinalBalance: a.Balance + amount}
	}

	a.Balance += amount

	return nil
}
func (a *BankAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{amount}
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if amount > a.Balance {
		return &ExceedsLimitError{MaxBalance: a.MinBalance, FinalBalance: a.Balance - amount}
	}

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{InitialBalance: a.Balance - amount, MinBalance: a.MinBalance}
	}

	a.Balance -= amount

	return nil
}

func (a *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if amount < 0 {
		return &NegativeAmountError{amount}
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if target.Balance+amount > MaxTransactionAmount {
		return &ExceedsLimitError{MaxBalance: a.Balance, FinalBalance: a.Balance - amount}
	}

	if a.Balance-amount < a.MinBalance {
		return &InsufficientFundsError{InitialBalance: a.Balance - amount, MinBalance: a.MinBalance}
	}
	a.Balance -= amount
	target.Balance += amount
	return nil
}

func main() {
	result, err := NewBankAccount("ab", "1", 1000, 100)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}

	result1, err := NewBankAccount("abc", "2", 1000, 100)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}

	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// defer wg.Done()
	// err = result.Deposit(500.0)
	// if err != nil {
	// 	fmt.Println("Error deposit 1: ", err)
	// 	return
	// }
	// }()

	// wg.Add(1)
	// go func() {
	// defer wg.Done()
	// err = result.Withdraw(200.0)
	// if err != nil {
	// 	fmt.Println("Error withdraw: ", err)
	// 	return
	// }
	// }()

	result.Transfer(200, result1)

	// wg.Wait()
	fmt.Println("acc 1", result.Balance)
	fmt.Println("acc 2", result1.Balance)
}
