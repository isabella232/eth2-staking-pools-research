package participant

import (
	pool_chain "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"log"
	"time"
)

type Participant struct {
	Id uint32
	node *pool_chain.PoolChainNode
}

func NewParticipant(id uint32) *Participant {
	return &Participant{
		Id:   id,
	}
}

func (p *Participant) SetNode(node *pool_chain.PoolChainNode) {
	p.node = node
}

func (p *Participant) StartEpochProcessing() {
	go func() {
		for {
			select {
			case epoch := <- p.node.EpochC():
				if epoch == -1 {
					return
				}

				e := p.node.State.GetEpoch(uint32(epoch))
				go p.timeEpoch(e)
			}
		}
	}()

	p.node.StartEpochProcessing()

	log.Printf("Participant %d started", p.Id)
}

func (p *Participant) KillC() <- chan bool {
	return p.node.Killed
}

func (p *Participant) timeEpoch(epoch *state.Epoch) {
	// start happens at 1/3 of the epoch
	go func() {
		d := time.Duration(p.node.Config.EpochSpanSec / 3)
		<- time.After(d)

		p.EpochStart(epoch)
	}()

	// start happens at 1/2 of the epoch
	go func() {
		d := time.Duration(p.node.Config.EpochSpanSec / 2)
		<- time.After(d)

		p.EpochMid(epoch)
	}()

	// start happens at 1/2 of the epoch
	go func() {
		d := time.Duration((p.node.Config.EpochSpanSec / 3) * 2)
		<- time.After(d)

		p.EpochEnd(epoch)
	}()
}

// start happens at 1/3 of the epoch
func (p *Participant) EpochStart(epoch *state.Epoch) {
	log.Println("epoch ", epoch.Number, "start")

	//share := &pb.ShareDistribution{
	//	Type:            pb.ShareType_EPOCH,
	//	Id:              uuid.New().String(),
	//	FromParticipant: &pb.Participant{Id:p.Id},
	//	ToParticipant:   &pb.Participant{Id:4},
	//	Share:           []byte(""),
	//	Commitments:     nil,
	//	PoolId:          1,
	//	Epoch:           epoch.Number,
	//}
	//
	//err := p.node.Net.BroadcastShare(share)
	//if err != nil {
	//	log.Printf("broadcasting error: %s", err.Error())
	//}
	//
	//
	//
	//
	//
	//
	//share2 := &pb.ShareDistribution{
	//	Type:            pb.ShareType_EPOCH,
	//	Id:              uuid.New().String(),
	//	FromParticipant: &pb.Participant{Id:p.Id},
	//	ToParticipant:   &pb.Participant{Id:5},
	//	Share:           []byte(""),
	//	Commitments:     nil,
	//	PoolId:          1,
	//	Epoch:           epoch.Number,
	//}
	//
	//err = p.node.Net.BroadcastShare(share2)
	//if err != nil {
	//	log.Printf("broadcasting error: %s", err.Error())
	//}
}

// start happens at 1/2 of the epoch
func (p *Participant) EpochMid(epoch *state.Epoch) {
	log.Printf("epoch %d mid with %d shares", epoch.Number, len(p.node.SharesPerEpoch[epoch.Number]))
}

// start happens at 2/3 of the epoch
func (p *Participant) EpochEnd(epoch *state.Epoch) {
	log.Printf("epoch %d end", epoch.Number)
}