package core

import (
	"crypto/sha256"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
)

func DecreaseBPBalance(bp *BlockProducer, change uint64) error {
	if bp.Balance < change {
		return fmt.Errorf("BP %d dosen't have enonugh Balance (%d) to decrease (%d)", bp.Id, bp.Balance, change)
	}

	bp.Balance -= change
	return nil
}

func IncreaseBPBalance(bp *BlockProducer, change uint64) error {
	bp.Balance += change
	return nil
}

// will return nil if not found
func GetBlockProducer(state *State, id uint64) *BlockProducer {
	for _, p := range state.BlockProducers {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

// will return nil if not found
func GetPool(state *State, id uint64) *Pool {
	for _, p := range state.Pools {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

func PoolCommittee(state *State, poolId uint64, epoch uint64) ([]uint64,error) {
	// TODO validate state.Seeds[epoch] returns not nil?
	return shuffle(
		state.BlockProducers,
		0,
		epoch,
		shared.SliceToByte32(state.Seeds[epoch]),
		[]byte(fmt.Sprintf("pool %d committee", poolId)),
	)
}

func DKGCommittee(state *State, reqId uint64, epoch uint64)([]uint64, error) {
	return shuffle(
		state.BlockProducers,
		0,
		epoch,
		shared.SliceToByte32(state.Seeds[epoch]),
		[]byte("dkg committee"),
	)
}

func BlockVotingCommittee(state *State, epoch uint64)([]uint64, error) {
	return shuffle(
		state.BlockProducers,
		0,
		epoch,
		shared.SliceToByte32(state.Seeds[epoch]),
		[]byte("block voting committee"),
	)
}

func GetBlockProposer(state *State, epoch uint64) (uint64, error) {
	lst, err := shuffle(
		state.BlockProducers,
		0,
		epoch,
		shared.SliceToByte32(state.Seeds[epoch]),
		[]byte("block proposer"),
	)
	if err != nil {
		return 0, err
	}
	return lst[0], nil
}


// TODO - find out if secure
func shuffle(allBPs []*BlockProducer, committeeId uint64, epoch uint64, seed [32]byte, nonce []byte) ([]uint64, error) {
	// get Active BPs
	var activeBps []uint64
	for _, bp := range allBPs {
		if bp.Active || bp.GetExitEpoch() > epoch {
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
	ret := make([]uint64, TestConfig().PoolExecutorsNumber)
	for i := uint64(0) ; i < TestConfig().PoolExecutorsNumber ; i++ {
		ret[i] = shuffled[committeeId* TestConfig().PoolExecutorsNumber + i]
	}

	return ret, nil
}