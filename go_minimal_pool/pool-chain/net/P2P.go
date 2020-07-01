package net

import "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"

type P2P interface {
	AddPeer(peer *Peer) error
	RemovePeer(peer *Peer) error
	BroadcastShare(distribution pb.ShareDistribution) error
}
