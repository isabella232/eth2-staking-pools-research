package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStateCopying(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateValidHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	// test copying
	bp := core.GetBlockProducer(state, 0)
	bp.Balance = 100000
	require.NotEqualValues(t, core.GetBlockProducer(newState, 0).Balance, 100000)

	core.GetPool(state, 0).Active = false
	require.NotEqualValues(t, core.GetPool(newState, 0).Active, core.GetPool(state, 0).Active)
}
