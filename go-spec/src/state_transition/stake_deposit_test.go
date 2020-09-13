package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidStakeDeposit(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	_, body := GenerateFinalizedAttestationPoolHeadAndBody(state)

	st := NewStateTransition()

	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	// verify stake added
	require.EqualValues(t, 100, core.GetBlockProducer(newState, 17).Stake)
}
