package pool_chain

import (
	net2 "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/simple_net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
)

type PoolChainNode struct {
	State       *state.State
	net         net2.P2P
	epochTicker *EpochTicker
	Config      *NetworkConfig

	Killed 		chan bool
}

func NewTestChainNode() *PoolChainNode {
	config := NewTestNetworkConfig()
	state := state.NewInMemoryState()
	ticker := NewEpochTicker(config.EpochSpanSec)
	net := simple_net.NewSimpleP2P()

	return &PoolChainNode{
		State:       state,
		net:         net,
		epochTicker: ticker,
		Config:      config,
		Killed:      make(chan bool),
	}
}

func (p *PoolChainNode) EpochC () <- chan int {
	return p.epochTicker.C()
}

func (p *PoolChainNode) StartEpochProcessing() {
	p.epochTicker.Start()
}

//func (p *PoolChainNode) timeEpoch(epoch *state.Epoch) {
//	// start happens at 1/3 of the epoch
//	go func() {
//		d := time.Duration(p.Config.EpochSpanSec / 3)
//		<- time.After(d)
//
//		for _, p := range p.participants {
//			p.EpochStart(epoch)
//		}
//	}()
//
//	// start happens at 1/2 of the epoch
//	go func() {
//		d := time.Duration(p.Config.EpochSpanSec / 2)
//		<- time.After(d)
//
//		for _, p := range p.participants {
//			p.EpochMid(epoch)
//		}
//	}()
//
//	// start happens at 1/2 of the epoch
//	go func() {
//		d := time.Duration((p.Config.EpochSpanSec / 3) * 2)
//		<- time.After(d)
//
//		for _, p := range p.participants {
//			p.EpochEnd(epoch)
//		}
//	}()
//}