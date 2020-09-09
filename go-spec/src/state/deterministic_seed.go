package state

import (
	"crypto/sha256"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
)

func (state *State) PoolCommittee(poolId uint64, epoch uint64) ([]uint64, error) {
	return shuffle(state.blockProducers, 0, epoch, state.GetSeed(epoch), []byte(fmt.Sprintf("pool %d committee", poolId)))
}

// Deterministic method to randomly pick DKG participants based on state seed and req number
func (state *State) DKGCommittee(reqId uint64, epoch uint64)([]uint64, error) {
	return shuffle(state.blockProducers, 0, epoch, state.GetSeed(epoch), []byte("dkg committee"))
}

func (state *State) BlockVotingCommittee(epoch uint64)([]uint64, error) {
	return shuffle(state.blockProducers, 0, epoch, state.GetSeed(epoch), []byte("block voting committee"))
}

func (state *State) GetBlockProposer(epoch uint64) (uint64, error) {
	lst, err := shuffle(state.blockProducers, 0, epoch, state.GetSeed(epoch), []byte("block proposer"))
	if err != nil {
		return 0, err
	}
	return lst[0], nil
}

// TODO - find out if secure
func shuffle(allBPs []*BlockProducer, committeeId uint64, epoch uint64, seed [32]byte, nonce []byte) ([]uint64, error) {
	// get active BPs
	var activeBps []uint64
	for _, bp := range allBPs {
		if bp.IsActive() || bp.ExitEpoch() > epoch {
			activeBps = append(activeBps, bp.GetId())
		}
	}

	// nonce is used as different categories for the seed
	seedToUse := seed
	if nonce != nil {
		h := sha256.New() // TODO - secure enough?
		_, err := h.Write(append(seed[:], nonce...))
		if err != nil {
			return []uint64{}, err
		}
		seedToUse = shared.SliceToByte32(h.Sum(nil))
	}

	// shuffle
	shuffled,err := shared.ShuffleList(activeBps, seedToUse, 60)
	if err != nil {
		return nil, err
	}

	//
	ret := make([]uint64, core.TestConfig().PoolExecutorsNumber)
	for i := uint64(0) ; i < core.TestConfig().PoolExecutorsNumber ; i++ {
		ret[i] = shuffled[committeeId* core.TestConfig().PoolExecutorsNumber + i]
	}

	return ret, nil
}