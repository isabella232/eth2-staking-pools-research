package shared

import (
	"github.com/prysmaticlabs/go-bitfield"
)


//// Slot committee is chosen randomly by shuffling a seed + category (block voting committee)
//// The previous epoch's seed is used to choose the block voting committee as the current one (the block's epoch)
//func SlotCommitteeByIndex(state *core.State, slot uint64, committeeIdx uint64)([]uint64, error) {
//	epoch := params.SlotToEpoch(slot)
//	slotInEpoch := slot - epoch * params.ChainConfig.SlotsInEpoch
//
//	// TODO - handle integer overflow
//	var seed []byte
//	var err error
//	if epoch == 0 {
//		seed, err = GetEpochSeed(state, epoch) // will return genesis seed
//	} else {
//		seed, err = GetEpochSeed(state, epoch - 1) // we always use the seed from previous epoch
//	}
//	if err != nil {
//		return []uint64{}, err
//	}
//	shuffled, err :=  shuffleActiveBPs(
//		GetActiveBlockProducers(state, epoch),
//		SliceToByte32(seed),
//		[]byte("slot committee"),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	retAll := CommitteeStructure(shuffled)[slotInEpoch]
//	if uint64(len(retAll)) < committeeIdx {
//		return nil, fmt.Errorf("committee index out of range")
//	}
//
//	return CommitteeStructure(shuffled)[slotInEpoch][committeeIdx], nil
//}

// returns the committee indexes for which bf[i] is set to 1
func AttestingIndices(bf bitfield.Bitfield, committee []uint64) []uint64 {
	indices := make([]uint64, 0, bf.Count())
	for _, idx := range bf.BitIndices() {
		if idx < len(committee) {
			indices = append(indices, committee[idx])
		}
	}
	return indices
}

//func CommitteeStructure(activeBps []uint64) map[uint64][][]uint64 /* slot -> []committee */ {
//	cntCommittees := uint64(len(activeBps)) / params.ChainConfig.MinAttestationCommitteeSize
//
//	// divide equally all BPs into committees
//	committees := make([][]uint64, cntCommittees)
//	cIdx := 0 // committee indx
//	for _, i := range activeBps {
//		if len(committees[cIdx]) == 0 {
//			committees[cIdx] = make([]uint64, 0)
//		}
//		committees[cIdx] = append(committees[cIdx], i)
//		cIdx++
//		if cIdx >= len(committees) { // reset to first committee
//			cIdx = 0
//		}
//	}
//
//	ret := make(map[uint64][][]uint64)
//	// structure committees
//	slot := uint64(0)
//	for _, c := range committees {
//		if len (ret[slot]) == 0 { // new
//			ret[slot] = make([][]uint64, 0)
//		}
//		ret[slot] = append(ret[slot], c)
//
//		slot++
//		if slot >= params.ChainConfig.SlotsInEpoch { // reset to first slot
//			slot = 0
//		}
//	}
//
//	return ret
//}

// Block producer is chosen randomly by shuffling a seed + category (block proposer)
//// The previous epoch's seed is used to choose the block producer as the current one (the block's epoch)
//func BlockProposer(state *core.State, slot uint64) (uint64, error) {
//	epoch := ComputeEpochAtSlot(slot)
//	slotInEpoch := slot - epoch * params.ChainConfig.SlotsInEpoch
//
//	// TODO - what seed should we take? last epoch? last finalized epoch?
//	var seed []byte
//	var err error
//	if epoch == 0 {
//		seed, err = GetEpochSeed(state, epoch) // will return genesis seed
//	} else {
//		seed, err = GetEpochSeed(state, epoch - 1) // we always use the seed from previous epoch
//	}
//	if err != nil {
//		return 0, err
//	}
//
//	lst, err := shuffleActiveBPs(
//		GetActiveBlockProducers(state, epoch),
//		SliceToByte32(seed),
//		[]byte("block proposer"),
//	)
//	if err != nil {
//		return 0, err
//	}
//	return lst[slotInEpoch], nil
//}
//
//// Shuffle takes in a list of block producers Ids, a seed and a nonce to create a unique shuffle for that
//// combination by hashing seed + nonce.
//// Changing the nonce for different purposes can be used as "categories" from the same seed
//// TODO - find out if secure
//func shuffleActiveBPs(bps []uint64, seed [32]byte, nonce []byte) ([]uint64, error) {
//	// nonce is used to randomly select multiple types of committees from the same seed
//	seedToUse := seed
//	if nonce != nil {
//		h := sha256.New() // TODO - secure enough?
//		_, err := h.Write(append(seed[:], nonce...))
//		if err != nil {
//			return []uint64{}, err
//		}
//		seedToUse = SliceToByte32(h.Sum(nil))
//	}
//
//	// shuffleActiveBPs
//	return ShuffleList(bps, seedToUse, 60)
//}
