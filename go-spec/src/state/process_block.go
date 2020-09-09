package state

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func (state *State) ValidateBlock(header core.IBlockHeader, body core.IBlockBody) error {
	currentEpoch := state.GetCurrentEpoch()

	// verify proposer is expected proposer
	expectedProposer, err := state.GetBlockProposer(currentEpoch)
	if err != nil {
		return err
	}
	proposer := body.GetProposer()
	if expectedProposer != proposer {
		return fmt.Errorf("block expectedProposer is worng, expected %d but received %d", expectedProposer, proposer)
	}

	// signing committee
	//committeeIds, err := state.BlockVotingCommittee(currentEpoch)
	//if err != nil {
	//	return err
	//}
	//committee := make([]core.IBlockProducer, len(committeeIds))
	//for i, id := range committeeIds {
	//	committee[i] = state.GetBlockProducer(id)
	//}

	// verify header block root matches
	exectedRoot,err := body.Root()
	if err != nil {
		return err
	}
	if bytes.Compare(exectedRoot, header.GetBlockRoot()) != 0 {
		return fmt.Errorf("signed block root does not match body root")
	}

	//
	err = header.Validate(state.GetBlockProducer(proposer))
	if err != nil {
		return err
	}

	err = body.Validate()
	if err != nil {
		return err
	}

	return nil
}

// Applies every pool performance to its relevant executors, decreasing and increasing balances.
func (state *State) ProcessPoolExecutions(summaries []core.IExecutionSummary) error {
	for _, summary := range summaries {
		pool := state.GetPool(summary.GetPoolId())
		if pool == nil {
			return fmt.Errorf("could not find pool %d", summary.GetPoolId())
		}

		if !pool.IsActive() {
			return fmt.Errorf("pool %d not Active", summary.GetPoolId())
		}

		if err := summary.ApplyOnState(state); err != nil {
			return err
		}
	}
	return nil
}

func (state *State) ProcessNewPoolRequests(requests []core.ICreatePoolRequest) error {
	currentBP := state.GetBlockProducer(state.GetCurrentEpoch())
	if currentBP == nil {
		return fmt.Errorf("could not find current proposer")
	}

	for _, req := range requests {
		if err := req.Validate(state, currentBP); err != nil {
			return err
		}

		// get DKG participants
		participants,err :=  state.DKGCommittee(req.GetId(), req.GetStartEpoch())
		if err != nil {
			return err
		}

		switch req.GetStatus() {
		case 0:
			// TODO if i'm the DKDG leader act uppon it
		case 1: // successful
			pk := &bls.PublicKey{}
			err := pk.Deserialize(req.GetCreatePubKey())
			if err != nil {
				return err
			}

			err = state.AddNewPool(&Pool{
				Id:              uint64(len(state.Pools) + 1),
				PubKey:          pk.Serialize(),
				SortedExecutors: []uint64{}, // TODO - POPULAT
			})
			if err != nil {
				return err
			}

			// reward/ penalty
			for i := 0 ; i < len(participants) ; i ++ {
				bp := state.GetBlockProducer(participants[i])
				if bp == nil {
					return fmt.Errorf("could not find BP %d", participants[i])
				}
				partic := req.GetParticipation()
				if shared.IsBitSet(partic[:], uint64(i)) {
					_, err := bp.IncreaseBalance(core.TestConfig().DKGReward)
					if err != nil {
						return err
					}
				} else {
					_, err := bp.DecreaseBalance(core.TestConfig().DKGReward)
					if err != nil {
						return err
					}
				}
			}

			// special reward for leader
			_, err = currentBP.IncreaseBalance(3* core.TestConfig().DKGReward)
			if err != nil {
				return err
			}
		case 2: // un-successful
			for i := 0 ; i < len(participants) ; i ++ {
				bp := state.GetBlockProducer(participants[i])
				if bp == nil {
					return fmt.Errorf("could not find BP %d", participants[i])
				}
				_, err := bp.DecreaseBalance(core.TestConfig().DKGReward)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// called when a new block was proposed
func (oldState *State) ProcessNewBlock(newBlockHeader core.IBlockHeader, newBlockBody core.IBlockBody) (newState core.IState, error error) {
	// copy the oldState to apply oldState transition on
	newState, err := oldState.Copy()
	if err != nil {
		return nil, err
	}

	// bump epoch
	newState.SetCurrentEpoch(newState.GetCurrentEpoch() + 1)

	// validate block header and body
	err = oldState.ValidateBlock(newBlockHeader, newBlockBody)
	if err != nil {
		return nil, err
	}

	//
	err = newState.ProcessPoolExecutions(newBlockBody.GetExecutionSummaries())
	if err != nil {
		return nil, err
	}

	//
	err = newState.ProcessNewPoolRequests(newBlockBody.GetNewPoolRequests())
	if err != nil {
		return nil, err
	}

	// update internal oldState vars
	newSeed, err := shared.MixSeed(newState.GetSeed(oldState.GetCurrentEpoch()), shared.SliceToByte32(newBlockHeader.GetSignature()[:32])) // TODO - use something else than the sig
	if err != nil {
		return nil, err
	}
	newState.SetSeed(newSeed, newState.GetCurrentEpoch())
	newState.SetCurrentEpoch(newState.GetCurrentEpoch())
	newState.SetBlockRoot(shared.SliceToByte32(newBlockHeader.GetBlockRoot()), newState.GetCurrentEpoch())
	//newState.SetHeadBlockHeader(newBlockHeader)

	// the CurrentEpoch's oldState root is not included inside the oldState root as it creates
	// a recursive dependency.
	newStateRoot, err := newState.Root()
	if err != nil {
		return nil, err
	}
	newState.SetStateRoot(newStateRoot, newState.GetCurrentEpoch())

	return newState, nil
}
