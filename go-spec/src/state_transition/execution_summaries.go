package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
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
		pool := shared.GetPool(state, summary.GetPoolId())
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
				for i:=0 ; i < int(params.ChainConfig.VaultSize) ; i++ {
					bp := shared.GetBlockProducer(state, executors[i])
					if bp == nil {
						return fmt.Errorf("BP %d not found", executors[i])
					}

					if !duty.Finalized {
						err := shared.DecreaseBPBalance(bp, 2*params.ChainConfig.BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						participation := duty.GetParticipation()
						if participation.BitAt(uint64(i)) {
							shared.IncreaseBPBalance(bp, params.ChainConfig.BaseEth2DutyReward)
						} else {
							err := shared.DecreaseBPBalance(bp, params.ChainConfig.BaseEth2DutyReward)
							if err != nil {
								return err
							}
						}
					}
				}
			case 1: // proposal
				for i:=0 ; i < int(params.ChainConfig.VaultSize) ; i++ {
					bp := shared.GetBlockProducer(state, executors[i])
					if bp == nil {
						return fmt.Errorf("BP %d not found", executors[i])
					}

					if !duty.Finalized {
						err := shared.DecreaseBPBalance(bp, 4*params.ChainConfig.BaseEth2DutyReward)
						if err != nil {
							return err
						}
					} else {
						participation := duty.GetParticipation()
						if participation[:].BitAt(uint64(i)) {
							shared.IncreaseBPBalance(bp, 2*params.ChainConfig.BaseEth2DutyReward)
						} else {
							err := shared.DecreaseBPBalance(bp, 2*params.ChainConfig.BaseEth2DutyReward)
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
