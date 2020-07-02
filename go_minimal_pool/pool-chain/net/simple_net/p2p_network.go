package simple_net

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
	"sync"
)

type SimpleP2PNetwork struct {
	peers []*net.Peer
	peersLock sync.Mutex
}

func NewSimpleP2P() *SimpleP2PNetwork {
	return &SimpleP2PNetwork{}
}

func (p *SimpleP2PNetwork) AddPeer(peer *net.Peer) error {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	p.peers = append(p.peers, peer)
	return nil
}

func (p *SimpleP2PNetwork) RemovePeer(peer *net.Peer) error {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	if len(p.peers) == 1 {
		p.peers = make([]*net.Peer, 0)
	} else {
		newPeers := make([]*net.Peer, len(p.peers)-1)
		for _, p := range p.peers {
			if p.Id != peer.Id {
				newPeers = append(newPeers,p)
			}
		}
		p.peers = newPeers
	}
	return nil
}

func (p *SimpleP2PNetwork) BroadcastShare(share *pb.ShareDistribution) error {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	for _, p := range p.peers {
		p.ReceiveShare(share)
	}

	return nil
}