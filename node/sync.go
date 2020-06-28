package node

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jhampac/pueblo/database"
)

func (n *Node) sync(ctx context.Context) error {
	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Searching for new peers and blocks...")
			n.doSync()
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *Node) doSync() {
	for _, peer := range n.knownPeers {
		// if the node is itself
		if n.ip == peer.IP && n.port == peer.Port {
			continue
		}

		fmt.Printf("Searching for new Peers and their Blocks and Peers: %q\n", peer.TCPAddress())

		status, err := queryPeerStatus(peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			fmt.Printf("Peer %q was removed from KnownPeers\n", peer.TCPAddress())
			n.RemovePeer(peer)
			continue
		}

		err = n.joinKnowPeers(peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		err = n.syncBlocks(peer, status)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		err = n.syncKnownPeers(peer, status)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
	}
}

func (n *Node) joinKnowPeers(peer PeerNode) error {
	if peer.connected {
		return nil
	}

	url := fmt.Sprintf(
		"http://%s%s?%s=%s&%s=%d",
		peer.TCPAddress(),
		endpointAddPeer,
		endpointAddPeerQueryKeyIP,
		n.ip,
		endpointAddPeerQueryKeyPort,
		n.port,
	)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	addPeerRes := AddPeerRes{}
	err = readRes(res, &addPeerRes)
	if err != nil {
		return err
	}
	if addPeerRes.Error != "" {
		return fmt.Errorf(addPeerRes.Error)
	}

	knownPeer := n.knownPeers[peer.TCPAddress()]
	knownPeer.connected = addPeerRes.Success

	n.AddPeer(knownPeer)

	if !addPeerRes.Success {
		return fmt.Errorf("unable to join KnowPeers of %q", peer.TCPAddress())
	}
	return nil
}

func (n *Node) syncBlocks(peer PeerNode, status StatusRes) error {
	localBlocksNumber := n.state.LatestBlock().Header.Number

	if status.Hash.IsEmpty() {
		return nil
	}

	if status.Number < localBlocksNumber {
		return nil
	}

	if status.Number == 0 && !n.state.LatestBlockHash().IsEmpty() {
		return nil
	}

	newBlocksCount := status.Number - localBlocksNumber
	if localBlocksNumber == 0 && status.Number == 0 {
		newBlocksCount = 1
	}
	fmt.Printf("Found %d new blocks from Peer %s\n", newBlocksCount, peer.TCPAddress())

	blocks, err := fetchBlocksFromPeer(peer, n.state.LatestBlockHash())
	if err != nil {
		return err
	}
	return n.state.AddBlocks(blocks)
}

func (n *Node) syncKnownPeers(peer PeerNode, status StatusRes) error {
	for _, statusPeer := range status.KnownPeers {
		if !n.IsKnownPeer(statusPeer) {
			fmt.Printf("Found new peer %s\n", statusPeer.TCPAddress())
			n.AddPeer(statusPeer)
		}
	}
	return nil
}

func fetchBlocksFromPeer(peer PeerNode, fromBlock database.Hash) ([]database.Block, error) {
	fmt.Printf("Importing blocks from Peer %s...\n", peer.TCPAddress())

	url := fmt.Sprintf(
		"http://%s%s?%s=%s",
		peer.TCPAddress(),
		endpointSync,
		endpointSyncQueryKeyFromBlock,
		fromBlock.Hex(),
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	syncRes := SyncRes{}
	err = readRes(res, &syncRes)
	if err != nil {
		return nil, err
	}
	return syncRes.Blocks, nil
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
