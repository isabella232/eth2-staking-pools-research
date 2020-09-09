package state_transition

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidSig(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateValidHeadAndBody(t)
	st := NewStateTransition()

	require.NoError(t, st.ValidateBlock(state, head, body))
}

func TestInValidSig(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateInvalidSigHeadAndBody(t)
	st := NewStateTransition()

	require.EqualError(t, st.ValidateBlock(state, head, body), "signature did not verify")
}

func TestWrongProposer(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateWrongProposerHeadAndBody(t)
	st := NewStateTransition()

	require.EqualError(t, st.ValidateBlock(state, head, body), "block expectedProposer is worng, expected 456 but received 455")
}

func TestInvalidProposer(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateInvalidProposerHeadAndBody(t)
	st := NewStateTransition()

	require.EqualError(t, st.ValidateBlock(state, head, body), "block expectedProposer is worng, expected 456 but received 4550000000")
}

func TestWrongRoot(t *testing.T) {
	state := generateTestState(t)
	head, body := GenerateWrongRootHeadAndBody(t)
	st := NewStateTransition()

	require.EqualError(t, st.ValidateBlock(state, head, body), "signed block root does not match body root")
}