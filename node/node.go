package node

import (
	"fmt"
	"net/http"

	"github.com/jhampac/pueblo/database"
)

// DefaultHTTPort can be configured
const DefaultHTTPort = 9000

//PeerNode is a Node with identifying properties for the calling Node
type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`
	IsActive    bool   `json:"is_active"`
}

// Node is a container on which services can be registered
type Node struct {
	dataDir string
	port    uint64

	state *database.State

	knownPeers []PeerNode
}

// New creates a new Node instance
func New(dataDir string, port uint64, bootstrap PeerNode) *Node {
	return &Node{
		dataDir:    dataDir,
		port:       port,
		knownPeers: []PeerNode{bootstrap},
	}
}

// NewPeerNode creates a new PeerNode instance
func NewPeerNode(ip string, port uint64, isBootstrap bool, isActive bool) PeerNode {
	return PeerNode{ip, port, isBootstrap, isActive}
}

// Run fires up a Node instance
func (n *Node) Run() error {
	fmt.Println(fmt.Sprintf("Listening on HTTP port: %d", n.port))

	state, err := database.NewStateFromDisk(n.dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	n.state = state

	http.HandleFunc("/node/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, state)
	})

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, state)
	})

	http.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		txAddHandler(w, r, state)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Popay your friends and family fast!")
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", n.port), nil)
}
