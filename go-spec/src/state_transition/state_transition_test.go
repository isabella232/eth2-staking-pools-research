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
	require.EqualValues(t, toByte("f3bd5f1c8baac2307ee3388e2e7dd7bcee2ab9ee141a6ffbc897f3a2b5e42170"), newsSeed)
}
