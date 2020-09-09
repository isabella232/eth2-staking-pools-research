package core

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/ulule/deepcopier"
	"sort"
)

type IStateTransition interface {
	ValidateBlock(state *State, header *BlockHeader, body *BlockBody) error
	ApplyBlockBody(oldState *State, newBlockHeader *BlockHeader, newBlockBody *BlockBody) (newState *State, err error)

	ProcessExecutionSummaries(state *State, summaries []*ExecutionSummary) error
	ProcessNewPoolRequests(state *State, summaries []*CreateNewPoolRequest) error
}

type StateTransition struct {

}

func (st *StateTransition) ApplyBlockBody(oldState *State, newBlockHeader *BlockHeader, newBlockBody *BlockBody) (newState *State, err error) {
	newState = &State{}
	deepcopier.Copy(oldState).To(newState)

	// validate
	if err := st.ValidateBlock(newState, newBlockHeader, newBlockBody); err != nil {
		return nil,err
	}

	// process
	if err := st.ProcessExecutionSummaries(newState, newBlockBody.ExecutionSummaries); err != nil {
		return nil,err
	}
	if err := st.ProcessNewPoolRequests(newState, newBlockBody.NewPoolReq); err != nil {
		return nil,err
	}

	// bump epoch
	newState.CurrentEpoch += 1
	// apply seed
	newSeed, err := shared.MixSeed(
		shared.SliceToByte32(newState.Seeds[oldState.GetCurrentEpoch()]), // previous seed
		shared.SliceToByte32(newBlockHeader.Signature[:32])) // TODO - use something else than the sig
	if err != nil {
		return nil, err
	}
	newState.Seeds[newState.CurrentEpoch] = newSeed[:]
	// add block root
	root, err := ssz.HashTreeRoot(newBlockBody)
	if err != nil {
		return nil, err
	}
	newState.BlockRoots[newState.CurrentEpoch] = root[:]
	// state root
	root, err = ssz.HashTreeRoot(newState)
	if err != nil {
		return nil, err
	}
	newState.StateRoots[newState.CurrentEpoch] = root[:]

	return newState, nil
}

func (st *StateTransition) ValidateBlock(state *State, header *BlockHeader, body *BlockBody) error {
	// verify proposer is expected proposer
	expectedProposer, err := GetBlockProposer(state, state.CurrentEpoch)
	if err != nil {
		return err
	}
	proposerId := body.GetProposer()
	if expectedProposer != proposerId {
		return fmt.Errorf("block expectedProposer is worng, expected %d but received %d", expectedProposer, proposerId)
	}

	// verify header block root matches
	exectedRoot,err := ssz.HashTreeRoot(body)
	if err != nil {
		return err
	}
	if bytes.Compare(exectedRoot[:], header.GetBlockRoot()) != 0 {
		return fmt.Errorf("signed block root does not match body root")
	}

	// validate signature
	proposer := GetBlockProducer(state, proposerId)
	if proposer == nil {
		return fmt.Errorf("proposer not found")
	}
	sig := &bls.Sign{}
	err = sig.Deserialize(header.Signature)
	if err != nil {
		return err
	}
	pk := &bls.PublicKey{}
	err = pk.Deserialize(proposer.GetPubKey())
	if err != nil {
		return err
	}
	if res := sig.VerifyHash(pk, header.BlockRoot); !res {
		return fmt.Errorf("signature did not verify")
	}

	// TODO - validate block?

	return nil
}

func (st *StateTransition) ProcessExecutionSummaries(state *State, summaries []*ExecutionSummary) error {
	for _, summary := range summaries {
		pool := GetPool(state, summary.GetPoolId())
		if pool != nil {
			return fmt.Errorf("could not find pool %d", summary.GetPoolId())
		}
		if !pool.Active {
			return fmt.Errorf("pool %d is not active", summary.GetPoolId())
		}

		executors := pool.GetSortedCommittee()

		for _, duty := range summary.GetDuties() {
			switch duty.GetType() {
			case 0: // attestation
				for i:=0 ; i < int(TestConfig().PoolExecutorsNumber) ; i++ {
					bp := GetBlockProducer(state, executors[i])
					if bp == nil {
						return fmt.Errorf("BP %d not found", executors[i])
					}

					if !duty.Finalized {
						err := DecreaseBPBalance(bp, 2*TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						participation := duty.GetParticipation()
						if shared.IsBitSet(participation, uint64(i)) {
							err := IncreaseBPBalance(bp, TestConfig().BaseEth2DutyReward)
							if err != nil {
								return err
							}
						} else {
							err := DecreaseBPBalance(bp, TestConfig().BaseEth2DutyReward)
							if err != nil {
								return err
							}
						}
					}
				}
			case 1: // proposal
				for i:=0 ; i < int(TestConfig().PoolExecutorsNumber) ; i++ {
					bp := GetBlockProducer(state, executors[i])
					if bp == nil {
						return fmt.Errorf("BP %d not found", executors[i])
					}

					if !duty.Finalized {
						err := DecreaseBPBalance(bp, 4*TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						participation := duty.GetParticipation()
						if shared.IsBitSet(participation[:], uint64(i)) {
							err := IncreaseBPBalance(bp, 2*TestConfig().BaseEth2DutyReward)
							if err != nil {
								return err
							}
						} else {
							err := DecreaseBPBalance(bp, 2*TestConfig().BaseEth2DutyReward)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (st *StateTransition) ProcessNewPoolRequests(state *State, summaries []*CreateNewPoolRequest) error {
	for _, req := range summaries {
		leader := GetBlockProducer(state, req.StartEpoch)
		if leader == nil {
			return fmt.Errorf("could not find new pool req leader")
		}

		// verify leader is correct
		if req.LeaderBlockProducer != leader.Id {
			return fmt.Errorf("new pool req leader incorrect")
		}
		// TODO - req Id is primary (non duplicate and incremental)
		// TODO - check that network has enough capitalization
		// TODO - check leader is not part of DKG Committee

		// get DKG participants
		participants,err :=  DKGCommittee(state, req.GetId(), req.GetStartEpoch())
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

			// get committee
			committee, err := DKGCommittee(state, req.Id, req.StartEpoch)
			sort.Slice(committee, func(i int, j int) bool {
				return committee[i] < committee[j]
			})

			state.Pools = append(state.Pools, &Pool{
				Id:              uint64(len(state.Pools) + 1),
				PubKey:          pk.Serialize(),
				SortedCommittee: committee,
			})
			if err != nil {
				return err
			}

			// reward/ penalty
			for i := 0 ; i < len(participants) ; i ++ {
				bp := GetBlockProducer(state, participants[i])
				if bp == nil {
					return fmt.Errorf("could not find BP %d", participants[i])
				}
				partic := req.GetParticipation()
				if shared.IsBitSet(partic[:], uint64(i)) {
					err := IncreaseBPBalance(bp, TestConfig().DKGReward)
					if err != nil {
						return err
					}
				} else {
					err := DecreaseBPBalance(bp, TestConfig().DKGReward)
					if err != nil {
						return err
					}
				}
			}

			// special reward for leader
			err = IncreaseBPBalance(leader, 3* TestConfig().DKGReward)
			if err != nil {
				return err
			}
		case 2: // un-successful
			for i := 0 ; i < len(participants) ; i ++ {
				bp := GetBlockProducer(state, participants[i])
				if bp == nil {
					return fmt.Errorf("could not find BP %d", participants[i])
				}
				err := DecreaseBPBalance(bp, TestConfig().DKGReward)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}