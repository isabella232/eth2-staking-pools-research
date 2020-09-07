package src

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateSummary() *PoolExecutionSummary {
	return &PoolExecutionSummary{
		PoolId:        []byte("0"),
		StartingEpoch: 0,
		EndEpoch:      1,
		Performance:   map[*BeaconDuty][16]byte{
			&BeaconDuty{
				Type:     0,
				Slot:     0,
				Included: false,
			}: [16]byte{},
		},
	}
}

func TestSuccessful(t *testing.T) {
	state := GenerateRandomState()
	summary := GenerateSummary()

	require.NoError(t, summary.ApplyOnState(state))
}