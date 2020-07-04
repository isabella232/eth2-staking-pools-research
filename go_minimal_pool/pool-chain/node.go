package pool_chain

import (
	net2 "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/pb"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/simple_net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
)

type PoolChainNode struct {
	State       *state.State
	Net         net2.P2P
	epochTicker *EpochTicker
	Config      *NetworkConfig

	// just holds all messages for convenience
	SharesPerEpoch map[uint32][]*pb.ShareDistribution

	Killed 		chan bool
}

func NewTestChainNode() *PoolChainNode {
	config := NewTestNetworkConfig()
	state := state.NewInMemoryState()
	ticker := NewEpochTicker(config.EpochSpanSec)
	net := simple_net.NewSimpleP2P()

	ret := &PoolChainNode{
		State:          state,
		Net:            net,
		epochTicker:    ticker,
		Config:         config,
		Killed:         make(chan bool),
		SharesPerEpoch: make(map[uint32][]*pb.ShareDistribution),
	}

	net.RegisterReceiver(ret)

	return ret
}

func (p *PoolChainNode) EpochC () <- chan int {
	return p.epochTicker.C()
}

func (p *PoolChainNode) StartEpochProcessing() {
	p.epochTicker.Start()
}

func (p *PoolChainNode) ReceiveShare(share *pb.ShareDistribution) {
	if p.SharesPerEpoch[share.Epoch] == nil {
		p.SharesPerEpoch[share.Epoch] = make([]*pb.ShareDistribution, 0)
	}

	p.SharesPerEpoch[share.Epoch] = append(p.SharesPerEpoch[share.Epoch], share)
}