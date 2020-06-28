package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

// State encapsulates all the business logic of the chain
type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile *os.File

	latestBlock     Block
	latestBlockHash Hash
	hasGenesisBlock bool
}

// NewStateFromDisk creates State with a genesis file
func NewStateFromDisk(dataDir string) (*State, error) {
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
	state := &State{balances, make([]Tx, 0), f, Block{}, Hash{}, false}

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

		err = applyTXs(blockFs.Value.TXs, state)
		if err != nil {
			return nil, err
		}

		state.latestBlock = blockFs.Value
		state.latestBlockHash = blockFs.Key
		state.hasGenesisBlock = true
	}

	return state, nil
}

// AddBlocks iterates over a slice of Block types and adds them to state
func (s *State) AddBlocks(blocks []Block) error {
	for _, b := range blocks {
		_, err := s.AddBlock(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddBlock adds a new Block to the db chain
func (s *State) AddBlock(b Block) (Hash, error) {
	pendingState := s.copy()

	err := applyBlock(b, pendingState)
	if err != nil {
		return Hash{}, err
	}

	blockHash, err := b.Hash()
	if err != nil {
		return Hash{}, err
	}

	bFS := BlockFS{blockHash, b}
	bFSJSON, err := json.Marshal(bFS)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", bFSJSON)

	_, err = s.dbFile.Write(append(bFSJSON, '\n'))
	if err != nil {
		return Hash{}, err
	}

	s.Balances = pendingState.Balances
	s.latestBlockHash = blockHash
	s.latestBlock = b

	return blockHash, nil
}

func (s *State) copy() State {
	c := State{}
	c.hasGenesisBlock = s.hasGenesisBlock
	c.latestBlock = s.latestBlock
	c.latestBlockHash = s.latestBlockHash
	c.txMempool = make([]Tx, len(s.txMempool))
	c.Balances = make(map[Account]uint)

	for acc, balance := range s.Balances {
		c.Balances[acc] = balance
	}

	for _, tx := range s.txMempool {
		c.txMempool = append(c.txMempool, tx)
	}

	return c
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

// NextBlockNumber return the block number from the state
func (s *State) NextBlockNumber() uint64 {
	// when would it ever not have the genesis block?
	if !s.hasGenesisBlock {
		return uint64(0)
	}
	return s.LatestBlock().Header.Number + 1
}

// applyBlock is a validation layer for incoming blocks
func applyBlock(b Block, s State) error {
	nextExpectedBlockNumber := s.latestBlock.Header.Number + 1

	if s.hasGenesisBlock && b.Header.Number != nextExpectedBlockNumber {
		return fmt.Errorf("next expected block must be %d not %d", nextExpectedBlockNumber, b.Header.Number)
	}

	// validate the incoming block parent hash equals the current hash; if this doesn't match then the chain is broken
	if s.hasGenesisBlock && s.latestBlock.Header.Number > 0 && !reflect.DeepEqual(b.Header.Parent, s.latestBlockHash) {
		return fmt.Errorf("next block parent hash must be '%x' not '%x'", s.latestBlockHash, b.Header.Parent)
	}

	return applyTXs(b.TXs, &s)
}

// applyTXs iterates through each txs and passes it to applyTx to be applied
func applyTXs(txs []Tx, s *State) error {
	for _, tx := range txs {
		err := updateBalances(tx, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateBalances(tx Tx, s *State) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("transaction error: Sender %q balance is %d and the request sent amount is %d", tx.From, s.Balances[tx.From], tx.Value)
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}
