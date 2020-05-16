package database

import (
	"os"
	"path/filepath"
)

// State encapsulates all the business logic of the chain
type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile *os.File
}

// NewStateFromDisk creates State with a genesis file
func NewStateFromDisk(*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}
}
