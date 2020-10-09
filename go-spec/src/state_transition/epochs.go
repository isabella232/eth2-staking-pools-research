package state_transition

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
)

//def process_epoch(state: BeaconState) -> None:
//	process_justification_and_finalization(state)
//	process_rewards_and_penalties(state)
//	process_registry_updates(state)
//	process_slashings(state)
//	process_final_updates(state)
func processEpoch(state *core.State) error {
	return nil
}

// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#justification-and-finalization
func processJustificationAndFinalization(state *core.State) error {
	if shared.GetCurrentEpoch(state) <= params.ChainConfig.GenesisEpoch + 1 {
		return nil
	}

	currentEpoch := shared.GetCurrentEpoch(state)
	previousEpoch := shared.GetPreviousEpoch(state)

	oldPrevJustificationPoint := state.PreviousJustifiedCheckpoint
	oldCurrentJustificationPoint := state.CurrentJustifiedCheckpoint

	// process justifications
	state.PreviousJustifiedCheckpoint = state.CurrentJustifiedCheckpoint
	newBits := state.JustificationBits
	newBits.Shift(1)
	state.JustificationBits = newBits

	prev, current, err := calculateAttestingBalances(state)
	if err != nil {
		return err
	}
	if prev.AttestingBalance * 3 >= prev.ActiveBalance * 2 {
		state.CurrentJustifiedCheckpoint = &core.Checkpoint{
			Epoch:                previousEpoch,
			Root:                 shared.GetBlockRoot(state, previousEpoch).Bytes,
		}
		newBits.SetBitAt(1, true)
		state.JustificationBits = newBits
	}
	if current.AttestingBalance * 3 >= current.ActiveBalance * 2 {
		state.CurrentJustifiedCheckpoint = &core.Checkpoint{
			Epoch:                currentEpoch,
			Root:                 shared.GetBlockRoot(state, currentEpoch).Bytes,
		}
		newBits.SetBitAt(0, true)
		state.JustificationBits = newBits
	}

	// process finalization
	justification := state.JustificationBits.Bytes()[0]

	// 2nd/3rd/4th (0b1110) most recent epochs are justified, the 2nd using the 4th as source.
	if justification&0x0E == 0x0E && (oldPrevJustificationPoint.Epoch+3) == currentEpoch {
		state.FinalizedCheckpoint = oldPrevJustificationPoint
	}

	// 2nd/3rd (0b0110) most recent epochs are justified, the 2nd using the 3rd as source.
	if justification&0x06 == 0x06 && (oldPrevJustificationPoint.Epoch+2) == currentEpoch {
		state.FinalizedCheckpoint = oldPrevJustificationPoint
	}

	// 1st/2nd/3rd (0b0111) most recent epochs are justified, the 1st using the 3rd as source.
	if justification&0x07 == 0x07 && (oldCurrentJustificationPoint.Epoch+2) == currentEpoch {
		state.FinalizedCheckpoint = oldCurrentJustificationPoint
	}

	// The 1st/2nd (0b0011) most recent epochs are justified, the 1st using the 2nd as source
	if justification&0x03 == 0x03 && (oldCurrentJustificationPoint.Epoch+1) == currentEpoch {
		state.FinalizedCheckpoint = oldCurrentJustificationPoint
	}
	return nil
}

type Balances struct {
	Epoch uint64
	ActiveBalance uint64
	AttestingIndexes []uint64
	AttestingBalance uint64
}

func calculateAttestingBalances(state *core.State) (prev *Balances, current *Balances, err error) {
	// TODO - assert epoch in [currentEpoch, previousEpoch]

	calc := func(attestations []*core.PendingAttestation, epoch uint64) (*Balances,error) {
		// filter matching att. by target root
		matchingAtt := make([]*core.PendingAttestation, 0)
		for _, att := range attestations {
			if blkRoot := shared.GetBlockRoot(state, epoch); blkRoot != nil {
				if bytes.Equal(att.Data.Target.Root, blkRoot.Bytes) {
					matchingAtt = append(matchingAtt, att)
				}
			} else {
				return nil, fmt.Errorf("could not find block root for epoch %d", epoch)
			}
		}

		ret := &Balances{
			Epoch:            epoch,
			ActiveBalance:    0,
			AttestingIndexes: []uint64{},
			AttestingBalance: 0,
		}

		// calculate attesting balance and indices
		for _, att := range matchingAtt {
			committee, err := shared.SlotCommitteeByIndex(state, att.Data.Slot, uint64(att.Data.CommitteeIndex))
			if err != nil {
				return nil, err
			}
			attestingIndices := shared.AttestingIndices(att.AggregationBits, committee)
			for _, idx := range attestingIndices {
				bp := shared.GetBlockProducer(state, idx)
				if bp != nil && !bp.Slashed {
					ret.AttestingIndexes = append(ret.AttestingIndexes, idx)
					ret.AttestingBalance += bp.Stake
				}
			}
		}

		// get active balance
		activeBps := shared.GetActiveBlockProducers(state, shared.GetCurrentEpoch(state))
		for _, idx := range activeBps {
			bp := shared.GetBlockProducer(state, idx)
			if bp != nil {
				ret.ActiveBalance += bp.Stake
			}
		}

		return ret, nil
	}

	currentEpoch := shared.GetCurrentEpoch(state)
	previousEpoch := shared.GetPreviousEpoch(state)
	prev, err = calc(state.PreviousEpochAttestations, previousEpoch)
	if err != nil {
		return nil, nil, err
	}
	current, err = calc(state.CurrentEpochAttestations, currentEpoch)
	if err != nil {
		return nil, nil, err
	}


	return prev, current, nil
}