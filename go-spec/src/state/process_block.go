package state

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func (state *State) ValidateBlock(header core.IBlockHeader, body core.IBlockBody) error {
	bp := state.GetBlockProducer(body.GetProposer())
	if bp == nil {
		return fmt.Errorf("could not find BP %d", body.GetProposer())
	}

	err := header.Validate(bp)
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
func (state *State) ProcessNewBlock(newBlockHeader core.IBlockHeader, newBlockBody core.IBlockBody) (newState core.IState, error error) {
	previousEpoch := state.GetCurrentEpoch()
	currentEpoch := previousEpoch + 1

	expectedProposer, err := state.GetBlockProposer(previousEpoch) // TODO - should it be the previous?
	if err != nil {
		return nil, err
	}
	proposer := newBlockBody.GetProposer()
	if expectedProposer != proposer {
		return nil, fmt.Errorf("block expectedProposer is worng, expected %d but received %d", expectedProposer, newBlockBody.GetProposer())
	}

	bp := state.GetBlockProducer(proposer)
	if bp == nil {
		return nil, fmt.Errorf("could not find BP %d", newBlockBody.GetProposer())
	}

	// copy the state to apply state transition on
	stateCopy, err := state.Copy()
	if err != nil {
		return nil, err
	}

	err = stateCopy.ProcessPoolExecutions(newBlockBody.GetExecutionSummaries())
	if err != nil {
		return nil, err
	}

	err = stateCopy.ProcessNewPoolRequests(newBlockBody.GetNewPoolRequests())
	if err != nil {
		return nil, err
	}

	// update internal state vars
	newSeed, err := shared.MixSeed(stateCopy.GetSeed(previousEpoch), shared.SliceToByte32(newBlockHeader.GetSignature()[:32])) // TODO - use something else than the sig
	if err != nil {
		return nil, err
	}
	stateCopy.SetSeed(newSeed, currentEpoch)
	stateCopy.SetCurrentEpoch(currentEpoch)
	stateCopy.SetBlockRoot(shared.SliceToByte32(newBlockHeader.GetBlockRoot()), currentEpoch)
	//stateCopy.SetHeadBlockHeader(newBlockHeader)

	// the CurrentEpoch's state root is not included inside the state root as it creates
	// a recursive dependency.
	newStateRoot, err := stateCopy.Root()
	if err != nil {
		return nil, err
	}
	stateCopy.SetStateRoot(newStateRoot, currentEpoch)

	return stateCopy, nil
}
