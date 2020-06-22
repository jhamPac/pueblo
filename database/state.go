package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// State encapsulates all the business logic of the chain
type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile *os.File

	latestBlock     Block
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

func (s *State) applyBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.apply(tx); err != nil {
			return err
		}
	}

	return nil
}

// Persist the Mempool to the dbFile
func (s *State) Persist() (Hash, error) {
	latestBlockHash, err := s.latestBlock.Hash()
	if err != nil {
		return Hash{}, err
	}

	block := NewBlock(
		s.latestBlockHash,
		s.latestBlock.Header.Number+1,
		uint64(time.Now().Unix()),
		s.txMempool,
	)

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

	s.latestBlockHash = latestBlockHash
	s.latestBlock = block
	s.txMempool = []Tx{}

	return latestBlockHash, nil
}

// Close the dbfile that State uses for mempool
func (s *State) Close() error {
	return s.dbFile.Close()
}

// LatestBlock returns the most recent block created
func (s *State) LatestBlock() Block {
	return s.latestBlock
}

// LatestBlockHash return the most recent block hash
func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

// NewStateFromDisk creates State with a genesis file
func NewStateFromDisk(dataDir string) (*State, error) {
	dataDir = ExpandPath(dataDir)

	err := initDataDirIfNotExists(dataDir)
	if err != nil {
		return nil, err
	}

	gen, err := loadGenesis(getGenesisJSONFilePath(dataDir))
	if err != nil {
		return nil, err
	}

	// create the starting point or beginning state of balances
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	f, err := os.OpenFile(getBlocksDbFilePath(dataDir), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f, Block{}, Hash{}}

	// replay all the transactions
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJSON := scanner.Bytes()

		if len(blockFsJSON) == 0 {
			break
		}

		var blockFs BlockFS
		err = json.Unmarshal(blockFsJSON, &blockFs)
		if err != nil {
			return nil, err
		}

		err = state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}

		state.latestBlock = blockFs.Value
		state.latestBlockHash = blockFs.Key
	}

	return state, nil
}
