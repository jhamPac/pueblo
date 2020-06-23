package node

import (
	"context"
	"fmt"
	"net/http"
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
			continue
		}

		// just counts new blocks for right now
		localBlockNumber := n.state.LatestBlock().Header.Number
		if localBlockNumber < status.Number {
			newBlocksCount := status.Number - localBlockNumber
			fmt.Printf("Found %d new blocks from Peer %s\n", newBlocksCount, peer.IP)
		}

		// looks for new peers
		for _, uPeer := range status.KnownPeers {
			_, isKnownPeer := n.knownPeers[uPeer.TCPAddress()]
			if !isKnownPeer {
				fmt.Printf("Found new Peer %s\n", uPeer.TCPAddress())
				n.knownPeers[uPeer.TCPAddress()] = uPeer
			}
		}
	}
}

func queryPeerStatus(peer PeerNode) (StatusRes, error) {
	url := fmt.Sprintf("http://%s:%s", peer.TCPAddress(), endpointStatus)
	res, err := http.Get(url)
	if err != nil {
		return StatusRes{}, err
	}

	statusRes := StatusRes{}
	err = readRes(res, &statusRes)
	if err != nil {
		return StatusRes{}, err
	}
	return statusRes, nil
}
