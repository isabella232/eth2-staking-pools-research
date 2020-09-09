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
	require.EqualValues(t, toByte("f9bbaed2f11ce8e9e67d0f519420eadc8fbd919d5ade1759e27fa6d912eb66ef"), newsSeed)
}


func TestCreatedNewPoolReq(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateCreatePoolHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	require.Equal(t, 6, len(newState.Pools))
}