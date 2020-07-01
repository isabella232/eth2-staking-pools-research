package simple_net

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
)

type SimpleP2PNetwork struct {

}

func NewSimpleP2P() *SimpleP2PNetwork {
	return &SimpleP2PNetwork{}
}

func (p *SimpleP2PNetwork) AddPeer(peer *net.Peer) error {
	return nil
}

func (p *SimpleP2PNetwork) RemovePeer(peer *net.Peer) error {
	return nil
}

func (p *SimpleP2PNetwork) BroadcastShare(distribution pb.ShareDistribution) error {
	return nil
}