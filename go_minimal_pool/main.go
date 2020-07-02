package main

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/minimal_pool/participant"
	pool_chain "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain"
	"log"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	participant := participant.NewParticipant(1)
	participant.SetNode(pool_chain.NewTestChainNode())
	participant.StartEpochProcessing()

	for {
		select {
		case sig := <- participant.KillC():
			if sig == true {
				fmt.Printf("killed")
				return
			}
		}
	}
}