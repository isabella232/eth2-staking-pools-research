package participant

import (
	pool_chain "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"log"
	"sync"
	"time"
)

type Participant struct {
	Id   uint32
	Node *pool_chain.PoolChainNode

	epochProcessingLock sync.Mutex
}

func NewParticipant(id uint32) *Participant {
	return &Participant{
		Id:   id,
	}
}

func (p *Participant) SetNode(node *pool_chain.PoolChainNode) {
	p.Node = node
	p.Node.FilterId = p.Id
}

func (p *Participant) StartEpochProcessing() {
	go func() {
		for {
			select {
			case epoch := <- p.Node.EpochC():
				if epoch == -1 {
					return
				}

				e := p.Node.State.GetEpoch(uint32(epoch))
				go p.timeEpoch(e)
			}
		}
	}()

	p.Node.StartEpochProcessing()

	log.Printf("Participant %d started", p.Id)
}

func (p *Participant) KillC() <- chan bool {
	return p.Node.Killed
}

func (p *Participant) timeEpoch(epoch *state.Epoch) {
	// start happens at 1/4 of the epoch
	go func() {
		d := time.Duration(p.Node.Config.EpochSpanSec / 4)
		<- time.After(d)

		p.epochInit(epoch)
	}()

	// mid happens at 1/2 of the epoch
	go func() {
		d := time.Duration(p.Node.Config.EpochSpanSec / 2)
		<- time.After(d)

		p.epochMid(epoch)
	}()

	// end happens at 3/4 of the epoch
	go func() {
		d := time.Duration((p.Node.Config.EpochSpanSec / 4) * 3)
		<- time.After(d)

		p.epochEnd(epoch)
	}()
}