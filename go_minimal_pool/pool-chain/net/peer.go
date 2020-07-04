package net

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
	"github.com/google/uuid"
)

type Peer struct {
	Id uuid.UUID
	receiver P2PReceiver
}

func NewPeer() *Peer {
	return &Peer{
		Id:       uuid.New(),
	}
}

func (peer *Peer) RegisterReceiver(r P2PReceiver) {
	peer.receiver = r
}

func (peer *Peer) ReceiveShare(share *pb.ShareDistribution) {
	peer.receiver.ReceiveShare(share)
}