package src

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateAttestationSuccessfulSummary() *PoolExecutionSummary {
	return &PoolExecutionSummary{
		PoolId:        0,
		Epoch:         1,
		Duties:        []*BeaconDuty{
			&BeaconDuty{
				Type:     0,
				Slot:     0,
				Included: true,
				Executors: [16]byte{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, // the first executor index is set to 1
			},
		},
	}
}

func TestAttestationSuccessful(t *testing.T) {
	state := GenerateRandomState()
	summary := GenerateAttestationSuccessfulSummary()

	require.NoError(t, summary.ApplyOnState(state))

	for _, duty := range summary.Duties {
		pool, err := state.GetPool(summary.PoolId)
		require.NoError(t, err)

		for i:=0 ; i < int(TestConfig().PoolExecutorsNumber) ; i++ {
			bp,err := state.GetBlockProducer(pool.SortedExecutors[i])
			require.NoError(t, err)

			if IsBitSet(duty.Executors[:], uint64(i)) {
				require.EqualValues(t, 1100, bp.Balance)
				require.EqualValues(t, 0, i)
			} else {
				require.EqualValues(t, 900, bp.Balance)
			}
		}
	}
}