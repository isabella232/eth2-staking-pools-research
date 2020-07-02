package net

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
	"github.com/google/uuid"
)

type Peer struct {
	Id uuid.UUID
}

func (peer *Peer) ReceiveShare(share *pb.ShareDistribution) {

}