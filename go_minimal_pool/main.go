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
	participants = runDKGForPools(poolData, config.PoolThreshold)

	// connect pools to each other
	for i, p1 := range participants {
		p1.Node.Net.AddPeer(p1.Node.Net.OwnPeer()) // add to self to receive shares
		for j := i+1 ; j < len(participants) ; j ++  {
			p2 := participants[j]
			net.BiDirectionalConnection(p1.Node.Net,p2.Node.Net)
		}
	}
	// start epoch processing
	for _, p := range participants {
		p.StartEpochProcessing()
	}

	for {
		select {
		case sig := <- participants[0].KillC():
			if sig == true {
				fmt.Printf("killed")
				return
			}
		}
	}
}

func runDKGForPools(poolData map[uint8][]uint32, threshold uint8) []*participant.Participant {
	ret := make([]*participant.Participant,0)
	for poolId, poolParticipants := range poolData {
		sks,err := runDKGForParticipants(threshold - 1, poolParticipants)
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

			ret = append(ret, p)
		}
	}

	return ret
}

func runDKGForParticipants(degree uint8, indexes []uint32) (map[uint32]*bls.Fr, error) {
	dkg,err := crypto.NewDKG(degree, indexes)
	if err != nil {
		return nil, err
	}

	return dkg.GroupSecrets(indexes)
}
