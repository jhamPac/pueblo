package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jhampac/pueblo/database"
)

// DefaultIP for a node
const DefaultIP = "127.0.0.1"

// DefaultHTTPort can be configured
const DefaultHTTPort = 9000
const endpointStatus = "/node/status"

const endpointSync = "/node/sync"
const endpointSyncQueryKeyFromBlock = "fromBlock"

const endpointAddPeer = "/node/peer"
const endpointAddPeerQueryKeyIP = "ip"
const endpointAddPeerQueryKeyPort = "port"

//PeerNode is a Node with identifying properties for the calling Node
type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`

	connected bool
}

// TCPAddress returns the IP address and port number for communication
func (pn PeerNode) TCPAddress() string {
	return fmt.Sprintf("%s:%d", pn.IP, pn.Port)
}

// Node is a container on which services can be registered
type Node struct {
	dataDir string
	ip      string
	port    uint64

	state *database.State

	knownPeers map[string]PeerNode
}

// New creates a new Node instance
func New(dataDir string, ip string, port uint64, bootstrap PeerNode) *Node {
	knownPeers := make(map[string]PeerNode)
	knownPeers[bootstrap.TCPAddress()] = bootstrap

	return &Node{
		dataDir:    dataDir,
		ip:         ip,
		port:       port,
		knownPeers: knownPeers,
	}
}

// NewPeerNode creates a new PeerNode instance
func NewPeerNode(ip string, port uint64, isBootstrap bool, connected bool) PeerNode {
	return PeerNode{ip, port, isBootstrap, connected}
}

// Run fires up a Node instance
func (n *Node) Run() error {
	ctx := context.Background()
	fmt.Println(fmt.Sprintf("Listening on: %s:%d", n.ip, n.port))

	state, err := database.NewStateFromDisk(n.dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	n.state = state

	go n.sync(ctx)

	http.HandleFunc(endpointSync, func(w http.ResponseWriter, r *http.Request) {
		syncHandler(w, r, n)
	})

	http.HandleFunc(endpointStatus, func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, n)
	})

	http.HandleFunc(endpointAddPeer, func(w http.ResponseWriter, r *http.Request) {
		addPeerHandler(w, r, n)
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

	return http.ListenAndServe(fmt.Sprintf("%s:%d", n.ip, n.port), nil)
}

// AddPeer if not in known peers
func (n *Node) AddPeer(peer PeerNode) {
	n.knownPeers[peer.TCPAddress()] = peer
}

// RemovePeer from knownpeers
func (n *Node) RemovePeer(peer PeerNode) {
	delete(n.knownPeers, peer.TCPAddress())
}

// IsKnownPeer checks weather the peer is in known peers
func (n *Node) IsKnownPeer(peer PeerNode) bool {
	if peer.IP == n.ip && peer.Port == n.port {
		return true
	}
	_, isKnownPeer := n.knownPeers[peer.TCPAddress()]
	return isKnownPeer
}
