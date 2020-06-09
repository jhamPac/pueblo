package node

import (
	"fmt"
	"net/http"

	"github.com/jhampac/pueblo/database"
)

const httpPort = 9000

// BalancesRes is a repsonse with State information JSON encoded
type BalancesRes struct {
	Hash     database.Hash             `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
}

// TxAddReq struct to add transactions to state
type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Data  string `json:"data"`
}

// Run initializes the node on specified port
func Run(dataDir string) error {
	state, err := database.NewStateFromDisk(dataDir)
	if err != nil {
		return
	}
	defer state.Close()

	http.HandleFunc("/balances/list", listBalancesHandler)

	http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func listBalancesHandler(w http.ResponseWriter, r *http.Request) {
	writeRes(w, BalancesRes{state.LatestBlockHash(), state.Balances})
}
