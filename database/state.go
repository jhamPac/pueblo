package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Snapshot of the latest state
type Snapshot [32]byte

// State encapsulates all the business logic of the chain
type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile   *os.File
	snapshot Snapshot
}

// Add a tx to the Mempool
func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)
	return nil
}

// Persist the Mempool to the dbFile
func (s *State) Persist() (Snapshot, error) {
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		txJSON, err := json.Marshal(s.txMempool[i])
		if err != nil {
			return Snapshot{}, err
		}

		fmt.Printf("Persisting new TX to disk:\n")
		fmt.Printf("\t%s\n", txJSON)
		if _, err = s.dbFile.Write(append(txJSON, '\n')); err != nil {
			return Snapshot{}, err
		}

		if err := s.doSnapshot(); err != nil {
			return Snapshot{}, err
		}
		fmt.Printf("\tNew DB Snapshot: %x\n", s.snapshot)

		s.txMempool = append(s.txMempool[:i], s.txMempool[i+1:]...)
	}

	return s.snapshot, nil
}

// Close the dbfile that State uses for mempool
func (s *State) Close() error {
	return s.dbFile.Close()
}

// LatestSnapshot returns the most recent hash of the db
func (s *State) LatestSnapshot() Snapshot {
	return s.snapshot
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

func (s *State) doSnapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}

	txsData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}

	s.snapshot = sha256.Sum256(txsData)

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
	txDbFilePath := filepath.Join(cwd, "database", "tx.json")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f, Snapshot{}}

	// replay all the transactions
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	err = state.doSnapshot()
	if err != nil {
		return nil, err
	}

	return state, nil
}
