package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
)

func (st *StateTransition) validateExecutionSummaries(state *core.State, summaries []*core.ExecutionSummary) error {
	// TODO - validate summaries epoch, should be in some range?
	return nil
}

func (st *StateTransition) processExecutionSummaries(state *core.State, summaries []*core.ExecutionSummary) error {
	if err := st.validateExecutionSummaries(state, summaries); err != nil {
		return err
	}

	// TODO - what if a BP doesn't have enough CDT for penalties?
	for _, summary := range summaries {
		pool := core.GetPool(state, summary.GetPoolId())
		if pool == nil {
			return fmt.Errorf("could not find pool %d", summary.GetPoolId())
		}
		if !pool.Active {
			return fmt.Errorf("pool %d is not active", summary.GetPoolId())
		}

		executors := pool.GetSortedCommittee()

		for _, duty := range summary.GetDuties() {
			switch duty.GetType() {
			case 0: // attestation
				for i:=0 ; i < int(core.TestConfig().VaultSize) ; i++ {
					bp := core.GetBlockProducer(state, executors[i])
					if bp == nil {
						return fmt.Errorf("BP %d not found", executors[i])
					}

					if !duty.Finalized {
						err := core.DecreaseBPBalance(bp, 2*core.TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						participation := duty.GetParticipation()
						if shared.IsBitSet(participation, uint64(i)) {
							core.IncreaseBPBalance(bp, core.TestConfig().BaseEth2DutyReward)
						} else {
							err := core.DecreaseBPBalance(bp, core.TestConfig().BaseEth2DutyReward)
							if err != nil {
								return err
							}
						}
					}
				}
			case 1: // proposal
				for i:=0 ; i < int(core.TestConfig().VaultSize) ; i++ {
					bp := core.GetBlockProducer(state, executors[i])
					if bp == nil {
						return fmt.Errorf("BP %d not found", executors[i])
					}

					if !duty.Finalized {
						err := core.DecreaseBPBalance(bp, 4*core.TestConfig().BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						participation := duty.GetParticipation()
						if shared.IsBitSet(participation[:], uint64(i)) {
							core.IncreaseBPBalance(bp, 2*core.TestConfig().BaseEth2DutyReward)
						} else {
							err := core.DecreaseBPBalance(bp, 2*core.TestConfig().BaseEth2DutyReward)
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
