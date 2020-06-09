package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jhampac/pueblo/database"
)

const httpPort = 9000

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

// Run initializes the node on specified port
func Run(dataDir string) error {
	state, err := database.NewStateFromDisk(dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, state)
	})
	http.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		txAddHandler(w, r, state)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
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

func writeErrRes(w http.ResponseWriter, err error) {
	jsonErrRes, _ := json.Marshal(ErrRes{err.Error()})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonErrRes)
}

func writeRes(w http.ResponseWriter, content interface{}) {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(contentJSON)
}

func readReq(r *http.Request, reqBody interface{}) error {
	reqBodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read request body: %s", err.Error())
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJSON, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body: %s", err.Error())
	}

	return nil
}
