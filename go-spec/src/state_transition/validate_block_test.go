package state_transition

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidSig(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateValidHeadAndBody(state)
	st := NewStateTransition()

	require.NoError(t, st.PreApplyValidateBlock(state, head, body))
}

func TestInValidSig(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateInvalidSigHeadAndBody(state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "signature did not verify")
}

func TestWrongProposer(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateWrongProposerHeadAndBody(state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "block expectedProposer is worng, expected 456 but received 455")
}

func TestInvalidProposer(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateInvalidProposerHeadAndBody(state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "block expectedProposer is worng, expected 456 but received 4550000000")
}

func TestWrongRoot(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateWrongRootHeadAndBody(state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "signed block root does not match body root")
}

func TestTooSmallRadao(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateTooSmallRandaoHeadAndBody(state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "RANDAO should be 32 byte")
}

func TestTooBigRadao(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateTooBigRandaoHeadAndBody(state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "RANDAO should be 32 byte")
}

func TestValidPostStateRoot(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateValidHeadAndBody(state)

	st := NewStateTransition()
	require.NoError(t, st.PreApplyValidateBlock(state, head, body))
	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)
	require.NoError(t, st.PostApplyValidateBlock(newState, head, body))
}

func TestInvalidPostStateRoot(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateValidHeadAndBody(state)

	head.StateRoot = toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41") // wrong

	st := NewStateTransition()
	require.NoError(t, st.PreApplyValidateBlock(state, head, body))
	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)
	require.EqualError(t, st.PostApplyValidateBlock(newState, head, body), "new block state root is wrong")
}