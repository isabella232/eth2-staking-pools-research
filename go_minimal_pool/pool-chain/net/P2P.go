package net

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
)

type P2PReceiver interface {
	ReceiveShare(share *pb.ShareDistribution)
}

type P2P interface {
	// returns the p2p own pre for connections
	OwnPeer() *Peer
	RegisterReceiver(r P2PReceiver)
	AddPeer(peer *Peer) error
	RemovePeer(peer *Peer) error
	BroadcastShare(share *pb.ShareDistribution) error
}


func BiDirectionalConnection(p1 P2P, p2 P2P) {
	p1.AddPeer(p2.OwnPeer())
	p2.AddPeer(p1.OwnPeer())
}
