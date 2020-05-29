package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Hash replaces the Snapshot type
type Hash [32]byte

// MarshelText encodes a has into a hex value
func (h Hash) MarshelText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

// UnmarshalText decodes hex into a string
func (h Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return
}

// Block for a batch of transactions
type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

// BlockHeader is meta data for a Block
type BlockHeader struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

// NewBlock creates and returns a Block
func NewBlock(parent Hash, time uint64, txs []Tx) Block {
	return Block{BlockHeader{parent, time}, txs}
}

// Hash creates the Hash for each individual Block
func (b Block) Hash() (Hash, error) {
	blockJSON, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJSON), nil
}
