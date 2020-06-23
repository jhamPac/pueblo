package node

import (
	"context"
	"fmt"
	"time"
)

func (n *Node) sync(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Searching for new peers and blocks...")
			n.fetchAndSync()
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *Node) fetchAndSync() {
	for _, peer := range n.knownPeers {
		status, err := queryPeerStatus(peer)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		localBlockNumber := n.state.LatestBlock().Header.Number
		if localBlockNumber < status.Number {
			newBlocksCount := status.Number - localBlockNumber
			fmt.Printf("Found %d new blocks from Peer %s\n", newBlocksCount, peer.IP)
		}

		for _, statusPeer := range status.KnownPeers {
			newPeer, isKnownPeer := n.knownPeers[statusPeer.TCPAddress()]
			if !isKnownPeer {
				fmt.Printf("Found new Peer %s\n", peer.TCPAddress())
				n.knownPeers[statusPeer.TCPAddress()] = newPeer
			}
		}
	}
}
