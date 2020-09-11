package state_transition

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidSig(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateValidHeadAndBody(t, state)
	st := NewStateTransition()

	require.NoError(t, st.PreApplyValidateBlock(state, head, body))
}

func TestInValidSig(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateInvalidSigHeadAndBody(t, state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "signature did not verify")
}

func TestWrongProposer(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateWrongProposerHeadAndBody(t, state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "block expectedProposer is worng, expected 456 but received 455")
}

func TestInvalidProposer(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateInvalidProposerHeadAndBody(t, state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "block expectedProposer is worng, expected 456 but received 4550000000")
}

func TestWrongRoot(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateWrongRootHeadAndBody(t, state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "signed block root does not match body root")
}

func TestTooSmallRadao(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateTooSmallRandaoHeadAndBody(t, state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "RANDAO should be 32 byte")
}

func TestTooBigRadao(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateTooBigRandaoHeadAndBody(t, state)
	st := NewStateTransition()

	require.EqualError(t, st.PreApplyValidateBlock(state, head, body), "RANDAO should be 32 byte")
}

func TestValidPostStateRoot(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateValidHeadAndBody(t, state)

	st := NewStateTransition()
	newState, err := st.ApplyBlock(state, body)
	require.NoError(t, err)

	require.NoError(t, st.PreApplyValidateBlock(state, head, body))
	require.NoError(t, st.PostApplyValidateBlock(newState, head, body))
}