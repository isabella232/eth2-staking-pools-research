package block

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	state2 "github.com/bloxapp/eth2-staking-pools-research/go-spec/src/state"
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateAttestationSuccessfulSummary() *PoolExecutionSummary {
	return &PoolExecutionSummary{
		PoolId:        0,
		Epoch:         1,
		Duties:        []*BeaconDuty{
			&BeaconDuty{
				dutyType:      0,
				slot:          0,
				finalized:     true,
				participation: [16]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, // the first executor index is set to 1
			},
		},
	}
}

func TestAttestationSuccessful(t *testing.T) {
	state := state2.GenerateRandomState(t)
	summary := GenerateAttestationSuccessfulSummary()

	require.NoError(t, summary.ApplyOnState(state))

	for _, duty := range summary.Duties {
		pool := state.GetPool(summary.PoolId)

		for i:=0 ; i < int(core.TestConfig().PoolExecutorsNumber) ; i++ {
			bp,err := state.GetBlockProducer(pool.SortedExecutors[i])
			require.NoError(t, err)

			if shared.IsBitSet(duty.Executors[:], uint64(i)) {
				require.EqualValues(t, 1100, bp.Balance)
				require.EqualValues(t, 0, i)
			} else {
				require.EqualValues(t, 900, bp.Balance)
			}
		}
	}
}