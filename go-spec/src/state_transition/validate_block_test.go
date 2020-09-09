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