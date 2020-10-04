package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFinalizedAttestation(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	_, body := GenerateFinalizedAttestationPoolHeadAndBody(state)

	st := NewStateTransition()

	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	// check rewards
	participation := bitfield.Bitlist{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0}
	committee := core.GetPool(newState, 3).SortedCommittee
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber; i++ { // pool id = 3
		bp := core.GetBlockProducer(newState, committee[i])
		if participation.BitAt(i) {
			require.EqualValues(t, 1100, bp.CDTBalance)
		} else {
			require.EqualValues(t, 900, bp.CDTBalance)
		}
	}
}

func TestNotFinalizedAttestation(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	_, body := GenerateNotFinalizedAttestationPoolHeadAndBody(state)

	st := NewStateTransition()

	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	// check rewards
	committee := core.GetPool(newState, 3).SortedCommittee
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber; i++ { // pool id = 3
		bp := core.GetBlockProducer(newState, committee[i])
		require.EqualValues(t, 800, bp.CDTBalance)
	}
}

func TestFinalizedProposal(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	_, body := GenerateFinalizedProposalPoolHeadAndBody(state)

	st := NewStateTransition()

	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	// check rewards
	participation := bitfield.Bitlist{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0}
	committee := core.GetPool(newState, 3).SortedCommittee
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber; i++ { // pool id = 3
		bp := core.GetBlockProducer(newState, committee[i])
		if participation.BitAt(i) {
			require.EqualValues(t, 1200, bp.CDTBalance)
		} else {
			require.EqualValues(t, 800, bp.CDTBalance)
		}
	}
}

func TestNotFinalizedProposal(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	_, body := GenerateNotFinalizedProposalPoolHeadAndBody(state)

	st := NewStateTransition()

	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	// check rewards
	committee := core.GetPool(newState, 3).SortedCommittee
	require.NoError(t, err)

	// test penalties/ rewards
	for i := uint64(0) ; i < core.TestConfig().DKGParticipantsNumber; i++ { // pool id = 3
		bp := core.GetBlockProducer(newState, committee[i])
		require.EqualValues(t, 600, bp.CDTBalance)
	}
}