package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandaoSeedMix(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateValidHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	newsSeed, err := core.GetSeed(newState, newState.CurrentEpoch)
	require.NoError(t, err)
	require.EqualValues(t, toByte("e232e8c5886ce4f8f89628766e2f5e75b8564be3e897d3ccfc8e57cc9ea9215d"), newsSeed)
}
