package main

import (
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/crypto"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/net"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain/state"
	"github.com/herumi/bls-eth-go-binary/bls"
	"log"
)

func main() {
	crypto.InitBLS()
	log.SetFlags(log.Lmicroseconds)

	config := net.NewTestNetworkConfig()

	seed,err := crypto.MixSeed(config.GenesisSeed, 0)
	if err != nil {
		log.Fatalf(err.Error())
	}

	epoch := state.NewEpochInstance(0, seed)

	poolData,err := epoch.PoolsParticipantIds()
	if err != nil {
		log.Fatalf(err.Error())
	}
	// DKG for pools
	for _, participants := range poolData {
		sks,err := runDKGForParticipant(config.PoolThreshold - 1, participants)
		if err != nil {
			log.Fatalf(err.Error())
		}

		for k, v := range sks {
			log.Printf("p %d, sk: %s", k, v.GetString(10))
		}
	}

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


//func shufflePools(seed [32]byte) (map[uint8][]uint32, error) {
//	config := main.NewTestNetworkConfig()
//
//	shuffled, err := ShuffleList(config.ParticipantIndexesList(), seed)
//	if err != nil {
//		return nil, err
//	}
//
//	ret := make(map[uint8][]uint32)
//	for p_id := uint8(0) ; p_id < config.NumberOfPools ; p_id ++ {
//		start := p_id * config.PoolSize
//		end := start + config.PoolSize
//		ret[p_id] = shuffled[start: end]
//	}
//
//	return ret, nil
//}