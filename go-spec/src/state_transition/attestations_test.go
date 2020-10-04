package state_transition

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessBlockAttestations(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	_, body := GenerateValidHeadAndBody(state)

	st := NewStateTransition()
	require.NoError(t, st.ProcessBlockAttestations(state, body))
}
