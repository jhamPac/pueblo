package node

import (
	"net/http"

	"github.com/jhampac/pueblo/database"
)

// ErrRes custom error type to write to responses
type ErrRes struct {
	Error string `json:"error"`
}

// BalancesRes is a repsonse with State information JSON encoded
type BalancesRes struct {
	Hash     database.Hash             `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
}

// TxAddReq struct to add transactions to state
type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

// TxAddRes is a response type for Hashs
type TxAddRes struct {
	Hash database.Hash `json:"block_hash"`
}

// StatusRes displays the status of a Node
type StatusRes struct {
	Hash       database.Hash       `json:"block_hash"`
	Number     uint64              `json:"block_number"`
	KnownPeers map[string]PeerNode `json:"known_peers"`
}

func statusHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	res := StatusRes{
		Hash:       node.state.LatestBlockHash(),
		Number:     node.state.LatestBlock().Header.Number,
		KnownPeers: node.knownPeers,
	}
	writeRes(w, res)
}

func listBalancesHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.LatestBlockHash(), state.Balances})
}

func txAddHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	req := TxAddReq{}
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	tx := database.NewTx(database.NewAccount(req.From), database.NewAccount(req.To), req.Value, req.Data)

	err = state.AddTx(tx)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	hash, err := state.Persist()
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, TxAddRes{hash})
}
