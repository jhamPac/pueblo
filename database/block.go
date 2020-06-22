package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Hash replaces the Snapshot type
type Hash [32]byte

// MarshalText encodes a has into a hex value
func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

// UnmarshalText decodes hex into a string
func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

// BlockFS is ...
type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

// Block for a batch of transactions
type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

// BlockHeader is meta data for Blocks
type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Number uint64 `json:"number"`
	Time   uint64 `json:"time"`
}

// NewBlock creates and returns a Block
func NewBlock(parent Hash, number uint64, time uint64, txs []Tx) Block {
	return Block{BlockHeader{parent, number, time}, txs}
}

// Hash creates the Hash for each individual Block
func (b Block) Hash() (Hash, error) {
	blockJSON, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJSON), nil
}
