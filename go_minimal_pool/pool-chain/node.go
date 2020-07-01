package pool_chain

import (
	net2 "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net/simple_net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"log"
	"time"
)

type PoolChainNode struct {
	state       *state.State
	net         net2.P2P
	epochTicker *EpochTicker
	config      *NetworkConfig
	Killed 		chan bool
}

func NewTestChainNode() *PoolChainNode {
	config := NewTestNetworkConfig()
	state := state.NewInMemoryState()
	ticker := NewEpochTicker(config.EpochSpanSec)
	net := simple_net.NewSimpleP2P()

	return &PoolChainNode{
		state:       state,
		net:         net,
		epochTicker: ticker,
		config:      config,
		Killed:		 make(chan bool),
	}
}

func (p *PoolChainNode) StartEpochProcessing() {
	go func() {
		for {
			select {
			case epoch := <- p.epochTicker.C():
				if epoch == -1 {
					return
				}

				e := p.state.GetEpoch(uint32(epoch))
				go p.timeEpoch(e)
			}
		}
	}()

	p.epochTicker.Start()
}

func (p *PoolChainNode) timeEpoch(epoch *state.Epoch) {
	// start
	go func() {
		d := time.Duration(p.config.EpochSpanSec / 3)
		<- time.After(d)
		p.executeEpochStart(epoch)
	}()

	// mid
	go func() {
		d := time.Duration(p.config.EpochSpanSec / 2)
		<- time.After(d)
		p.executeEpochMid(epoch)
	}()

	// end
	go func() {
		d := time.Duration((p.config.EpochSpanSec / 3) * 2)
		<- time.After(d)
		p.executeEpochEnd(epoch)
	}()
}

// start happens at 1/3 of the epoch
func (p *PoolChainNode) executeEpochStart(epoch *state.Epoch) {
	log.Println("epoch %d start", epoch.Number)
}

// start happens at 1/2 of the epoch
func (p *PoolChainNode) executeEpochMid(epoch *state.Epoch) {
	log.Println("epoch %d mid", epoch.Number)
}

// start happens at 2/3 of the epoch
func (p *PoolChainNode) executeEpochEnd(epoch *state.Epoch) {
	log.Println("epoch %d end", epoch.Number)
}