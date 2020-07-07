package main

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/participant"
	pool_chain "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"github.com/herumi/bls-eth-go-binary/bls"
	"log"
)

var participants []*participant.Participant

func main() {
	crypto.InitBLS()
	log.SetFlags(log.Lmicroseconds)

	config := net.NewTestNetworkConfig()
	participants = make([]*participant.Participant, config.TotalNumberOfParticipants())

	// simulate epoch 0 and get initial pool assignments
	seed,err := crypto.MixSeed(config.GenesisSeed, 0)
	if err != nil {
		log.Fatalf(err.Error())
	}
	epoch := state.NewEpochInstance(0, seed)
	poolData,err := epoch.PoolsParticipantIds()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// initial DKG for pools
	for poolId, poolParticipants := range poolData {
		sks,err := runDKGForParticipant(config.PoolThreshold - 1, poolParticipants)
		if err != nil {
			log.Fatalf(err.Error())
		}

		log.Printf("pool %d:", poolId)
		for k, v := range sks {
			log.Printf("		p %d, sk: %s", k, v.GetString(10))

			// create the participant
			p := participant.NewParticipant(k)
			n := pool_chain.NewTestChainNode()
			p.SetNode(n)
			// set secret
			e := n.State.GetEpoch(0)
			e.ParticipantShare = v
			n.State.SaveEpoch(e)

			participants[k-1] = p
		}
	}

	fmt.Printf("")

	//n1 := pool_chain.NewTestChainNode()
	//n2 := pool_chain.NewTestChainNode()
	//net.BiDirectionalConnection(n1.Net, n2.Net)
	//
	//p1 := participant.NewParticipant(1)
	//p1.SetNode(n1)
	//p2 := participant.NewParticipant(2)
	//p2.SetNode(n2)
	//
	//p1.StartEpochProcessing()
	//p2.StartEpochProcessing()
	//
	//for {
	//	select {
	//	case sig := <- p1.KillC():
	//		if sig == true {
	//			fmt.Printf("killed")
	//			return
	//		}
	//	}
	//}
}

func runDKGForParticipant(degree uint8, indexes []uint32) (map[uint32]*bls.Fr, error) {
	dkg,err := crypto.NewDKG(degree, indexes)
	if err != nil {
		return nil, err
	}

	return dkg.GroupSecrets(indexes)
}
