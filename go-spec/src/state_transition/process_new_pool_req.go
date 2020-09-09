package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"sort"
)

func (st *StateTransition) ProcessNewPoolRequests(state *core.State, summaries []*core.CreateNewPoolRequest) error {
	for _, req := range summaries {
		leader := core.GetBlockProducer(state, req.StartEpoch)
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
		participants,err :=  core.DKGCommittee(state, req.GetId(), req.GetStartEpoch())
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
			committee, err := core.DKGCommittee(state, req.Id, req.StartEpoch)
			sort.Slice(committee, func(i int, j int) bool {
				return committee[i] < committee[j]
			})

			state.Pools = append(state.Pools, &core.Pool{
				Id:              uint64(len(state.Pools) + 1),
				PubKey:          pk.Serialize(),
				SortedCommittee: committee,
			})
			if err != nil {
				return err
			}

			// reward/ penalty
			for i := 0 ; i < len(participants) ; i ++ {
				bp := core.GetBlockProducer(state, participants[i])
				if bp == nil {
					return fmt.Errorf("could not find BP %d", participants[i])
				}
				partic := req.GetParticipation()
				if shared.IsBitSet(partic[:], uint64(i)) {
					err := core.IncreaseBPBalance(bp, core.TestConfig().DKGReward)
					if err != nil {
						return err
					}
				} else {
					err := core.DecreaseBPBalance(bp, core.TestConfig().DKGReward)
					if err != nil {
						return err
					}
				}
			}

			// special reward for leader
			err = core.IncreaseBPBalance(leader, 3* core.TestConfig().DKGReward)
			if err != nil {
				return err
			}
		case 2: // un-successful
			for i := 0 ; i < len(participants) ; i ++ {
				bp := core.GetBlockProducer(state, participants[i])
				if bp == nil {
					return fmt.Errorf("could not find BP %d", participants[i])
				}
				err := core.DecreaseBPBalance(bp, core.TestConfig().DKGReward)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
