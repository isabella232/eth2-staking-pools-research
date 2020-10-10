package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEpochJustification(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t, 95)
	err := populateJustificationAndFinalization(state, 2, 95, 1, &core.Checkpoint{
		Epoch:                2,
		Root:                 toByte("155bb576077c9a88f5f71f6ae1c235d368b39e1da2b3d99efd60239ad622e58e"),
	})

	require.NoError(t, err)
	require.NoError(t, processJustificationAndFinalization(state))
	require.EqualValues(t, toByte("155bb576077c9a88f5f71f6ae1c235d368b39e1da2b3d99efd60239ad622e58e"), state.CurrentJustifiedCheckpoint.Root)
}
