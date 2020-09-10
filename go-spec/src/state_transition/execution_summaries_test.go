package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFinalizedAttestation(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateFinalizedAttestationPoolHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	// check rewards
	participation := []byte{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0}
	committee := core.GetPool(newState, 3).SortedCommittee
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber; i++ { // pool id = 3
		bp := core.GetBlockProducer(newState, committee[i])
		if shared.IsBitSet(participation, i) {
			require.EqualValues(t, 1100, bp.Balance)
		} else {
			require.EqualValues(t, 900, bp.Balance)
		}
	}
}

func TestNotFinalizedAttestation(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateNotFinalizedAttestationPoolHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	// check rewards
	committee := core.GetPool(newState, 3).SortedCommittee
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber; i++ { // pool id = 3
		bp := core.GetBlockProducer(newState, committee[i])
		require.EqualValues(t, 800, bp.Balance)
	}
}

func TestFinalizedProposal(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateFinalizedProposalPoolHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	// check rewards
	participation := []byte{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0}
	committee := core.GetPool(newState, 3).SortedCommittee
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber; i++ { // pool id = 3
		bp := core.GetBlockProducer(newState, committee[i])
		if shared.IsBitSet(participation, i) {
			require.EqualValues(t, 1200, bp.Balance)
		} else {
			require.EqualValues(t, 800, bp.Balance)
		}
	}
}

func TestNotFinalizedProposal(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	head, body := GenerateNotFinalizedProposalPoolHeadAndBody(t)

	st := NewStateTransition()

	newState, err := st.ApplyBlockBody(state, head, body)
	require.NoError(t, err)

	// check rewards
	committee := core.GetPool(newState, 3).SortedCommittee
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber; i++ { // pool id = 3
		bp := core.GetBlockProducer(newState, committee[i])
		require.EqualValues(t, 600, bp.Balance)
	}
}