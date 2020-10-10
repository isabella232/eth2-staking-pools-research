package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEpoch2Justification(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t, 95)
	err := populateJustificationAndFinalization(state, 2, 95, 1, &core.Checkpoint{
		Epoch:                2,
		Root:                 toByte("74a631caa345567967d5998fd21c1d17513976c8d53e286525968e52e3e54499"),
	})

	require.NoError(t, err)
	require.NoError(t, processJustificationAndFinalization(state))
	require.EqualValues(t, toByte("74a631caa345567967d5998fd21c1d17513976c8d53e286525968e52e3e54499"), state.CurrentJustifiedCheckpoint.Root)
}
