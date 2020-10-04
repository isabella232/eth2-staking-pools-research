package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessBlockAttestations(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	body := &core.BlockBody{
		Slot:                 35,
		Attestations:         generateAttestations(state,86, 35,0,true, 0 /* attestation */),
	}

	st := NewStateTransition()
	require.NoError(t, st.ProcessBlockAttestations(state, body))
}

func TestProcessBlockAttestationsWithoutThreshold(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	body := &core.BlockBody{
		Slot:                 35,
		Attestations:         generateAttestations(state,85, 35,0,true, 0 /* attestation */),
	}

	st := NewStateTransition()
	require.Error(t, st.ProcessBlockAttestations(state, body), "attestation did not pass threshold")
}
