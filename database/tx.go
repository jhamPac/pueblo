package database

// Account represents a user address or customer
type Account string

// NewAccount takes in a value and returns a new Account type
func NewAccount(value string) Account {
	return Account(value)
}

// Tx is a transaction
type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

// NewTx creates a new transaction
func NewTx(from Account, to Account, value uint, data string) Tx {
	return Tx{from, to, value, data}
}

// IsReward checks if the transaction has an awarded bonus
func (t Tx) IsReward() bool {
	return t.Data == "reward"
}
