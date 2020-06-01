package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// State encapsulates all the business logic of the chain
type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile          *os.File
	latestBlockHash Hash
}

// AddTx adds a Tx during the AddBlock process
func (s *State) AddTx(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)
	return nil
}

// Persist the Mempool to the dbFile
func (s *State) Persist() (Hash, error) {
	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMempool)
	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, nil
	}

	blockFs := BlockFS{blockHash, block}

	blockFsJSON, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Println("Persisting new Block to disk:")
	fmt.Printf("\t%s\n", blockFsJSON)

	if _, err = s.dbFile.Write(append(blockFsJSON, '\n')); err != nil {
		return Hash{}, nil
	}

	s.latestBlockHash = blockHash

	s.txMempool = []Tx{}

	return s.latestBlockHash, nil
}

// Close the dbfile that State uses for mempool
func (s *State) Close() error {
	return s.dbFile.Close()
}

// LatestBlockHash return the most recent block hash
func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

// AddBlock adds a new Block to the db chain
func (s *State) AddBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.AddTx(tx); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient funds")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

// NewStateFromDisk creates State with a genesis file
func NewStateFromDisk() (*State, error) {
	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}

	// create the starting point or beginning state of balances
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	// retrieve all the transactions
	txDbFilePath := filepath.Join(cwd, "database", "blocks.json")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f, Hash{}}

	// replay all the transactions
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJSON := scanner.Bytes()
		var blockFs BlockFS
		err = json.Unmarshal(blockFsJSON, &blockFs)
		if err != nil {
			return nil, err
		}

		err = state.AddBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}

		state.latestBlockHash = blockFs.Key
	}

	return state, nil
}
