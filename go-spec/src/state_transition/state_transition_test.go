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
	_, body := GenerateValidHeadAndBody(t, state)

	st := NewStateTransition()

	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	newsSeed, err := core.GetSeed(newState, newState.CurrentEpoch)
	require.NoError(t, err)
	require.EqualValues(t, toByte("e4a17401658219365021cf584f4758d4b22ec861d9653e8249c8b4f73285a909"), newsSeed)
}
