package simple_net

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
	"sync"
)

type SimpleP2PNetwork struct {
	myPeer *net.Peer
	peers []*net.Peer
	peersLock sync.Mutex
	receiver net.P2PReceiver
}

func NewSimpleP2P() *SimpleP2PNetwork {
	return &SimpleP2PNetwork{
		myPeer: net.NewPeer(),
	}
}

func (p *SimpleP2PNetwork) RegisterReceiver(r net.P2PReceiver) {
	p.receiver = r
}

// returns this p2p network own peer
func (p *SimpleP2PNetwork) OwnPeer() *net.Peer {
	return p.myPeer
}

func (p *SimpleP2PNetwork) AddPeer(peer *net.Peer) error {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	peer.RegisterReceiver(p.receiver)
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