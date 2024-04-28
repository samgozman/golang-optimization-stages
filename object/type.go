package object

import "time"

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type BankDetails struct {
	BankID     string `json:"bank_id"`
	BankBranch string `json:"bank_branch"`
	AccountID  string `json:"account_id"`
}

type ExternalIDs struct {
	UserID    string `json:"user_id"`
	BankID    string `json:"bank_id"`
	AccountID string `json:"account_id"`
}

type Balance struct {
	Current string `json:"current"`
	Pending string `json:"pending"`
}

type Transaction struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type Card struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Number     string `json:"number"`
	ExpiryDate string `json:"expiry_date"`
	Status     string `json:"status"`
}

type Loan struct {
	ID           string    `json:"id"`
	Amount       string    `json:"amount"`
	Status       string    `json:"status"`
	Installments int       `json:"installments"`
	InterestRate string    `json:"interest_rate"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

type Notifications struct {
	Email bool `json:"email"`
	SMS   bool `json:"sms"`
	Push  bool `json:"push"`
}

type Preferences struct {
	Language      string        `json:"language"`
	Notifications Notifications `json:"notifications"`
}

// User a dummy structure to represent a user for the sake of this example
type User struct {
	ID           string        `json:"id"`
	FirstName    string        `json:"firstname"`
	LastName     string        `json:"lastname"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone"`
	Address      Address       `json:"address"`
	BankDetails  BankDetails   `json:"bank_details"`
	ExternalIDs  ExternalIDs   `json:"external_ids"`
	Balance      Balance       `json:"balance"`
	Transactions []Transaction `json:"transactions"`
	Cards        []Card        `json:"cards"`
	Loans        []Loan        `json:"loans"`
	Preferences  Preferences   `json:"preferences"`
}
