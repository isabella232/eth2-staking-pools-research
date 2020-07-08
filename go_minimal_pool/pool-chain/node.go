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
	Config      *net2.NetworkConfig

	// just holds all messages for convenience
	SharesPerEpoch map[uint32]map[string]*pb.ShareDistribution
	// messages will be saved only for the specific Id
	FilterId uint32

	Killed 		chan bool
}

func NewTestChainNode() *PoolChainNode {
	config := net2.NewTestNetworkConfig()
	state := state.NewInMemoryState(config.GenesisSeed)
	ticker := NewEpochTicker(config.EpochSpanSec)
	net := simple_net.NewSimpleP2P()

	ret := &PoolChainNode{
		State:          state,
		Net:            net,
		epochTicker:    ticker,
		Config:         config,
		Killed:         make(chan bool),
		SharesPerEpoch: make(map[uint32]map[string]*pb.ShareDistribution),
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
		p.SharesPerEpoch[share.Epoch] = make(map[string]*pb.ShareDistribution)
	}

	// filter only relevant messages
	if share.ToParticipant.Id == p.FilterId {
		// do not insert duplicates
		if p.SharesPerEpoch[share.Epoch][share.Id] == nil {
			p.SharesPerEpoch[share.Epoch][share.Id] = share
		}
	}
}