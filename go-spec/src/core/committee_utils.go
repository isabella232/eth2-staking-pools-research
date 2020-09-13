package core

import (
	"crypto/sha256"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
)

// Pool committee is a randomly selected committee of BPs that are chosen to receive a pool's
// keys via key rotation.
// For new pools, the first committee is responsible for doing DKG
//
// Pool committee is chosen randomly by shuffling a seed + category (pool %d committee)
// The previous epoch's seed is used to choose the DKG committee as the current one (the block's epoch)
func PoolCommittee(state *State, poolId uint64, epoch uint64) ([]uint64,error) {
	// TODO - handle integer overflow
	seed, err := GetSeed(state, epoch - 1) // we always use the seed from previous epoch
	if err != nil {
		return []uint64{}, err
	}
	return shuffleActiveBPs(
		GetActiveBlockProducers(state, epoch),
		shared.SliceToByte32(seed),
		[]byte(fmt.Sprintf("pool %d committee", poolId)),
	)
}

// Block voting committee is chosen randomly by shuffling a seed + category (block voting committee)
// The previous epoch's seed is used to choose the block voting committee as the current one (the block's epoch)
func BlockVotingCommittee(state *State, epoch uint64)([]uint64, error) {
	// TODO - handle integer overflow
	seed, err := GetSeed(state, epoch - 1) // we always use the seed from previous epoch
	if err != nil {
		return []uint64{}, err
	}
	return shuffleActiveBPs(
		GetActiveBlockProducers(state, epoch),
		shared.SliceToByte32(seed),
		[]byte("block voting committee"),
	)
}

// Shuffle takes in a list of block producers Ids, a seed and a nonce to create a unique shuffle for that
// combination by hashing seed + nonce.
// Changing the nonce for different purposes can be used as "categories" from the same seed
// TODO - find out if secure
func shuffleActiveBPs(bps []uint64, seed [32]byte, nonce []byte) ([]uint64, error) {
	// nonce is used to randomly select multiple types of committees from the same seed
	seedToUse := seed
	if nonce != nil {
		h := sha256.New() // TODO - secure enough?
		_, err := h.Write(append(seed[:], nonce...))
		if err != nil {
			return []uint64{}, err
		}
		seedToUse = shared.SliceToByte32(h.Sum(nil))
	}

	// shuffleActiveBPs
	shuffled,err := shared.ShuffleList(bps, seedToUse, 60)
	if err != nil {
		return nil, err
	}

	//
	ret := make([]uint64, TestConfig().PoolExecutorsNumber)
	for i := uint64(0) ; i < TestConfig().PoolExecutorsNumber ; i++ {
		ret[i] = shuffled[TestConfig().PoolExecutorsNumber + i]
	}

	return ret, nil
}
