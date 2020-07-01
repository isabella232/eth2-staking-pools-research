package main

import (
	"fmt"
	pool_chain "github.com/bloxapp/eth2-staking-pools-research/minimal_pool/pool-chain"
	"log"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	node := pool_chain.NewTestChainNode()
	node.StartEpochProcessing()

	for {
		select {
		case sig := <- node.Killed:
			if sig == true {
				fmt.Printf("killed")
				return
			}

		}
	}
}