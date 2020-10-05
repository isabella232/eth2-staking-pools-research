package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAttestationProcessing(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	tests := []struct{
		name string
		blockBody *core.BlockBody
		expectedError error
	}{
		{
			name: "valid block attestation",
			blockBody: &core.BlockBody{
				Slot:                 32,
				Attestations:         generateAttestations(
					state,
					86,
					32,
					&core.Checkpoint{Epoch: 0, Root: params.ChainConfig.ZeroHash},
					&core.Checkpoint{Epoch: 1, Root: []byte{}},
					0,
					true,
					0, /* attestation */
					),
			},
			expectedError: nil,
		},
		{
			name: "threshold sig not achieved",
			blockBody: &core.BlockBody{
				Slot:                 32,
				Attestations:         generateAttestations(
					state,
					85,
					32,
					&core.Checkpoint{Epoch: 0, Root: params.ChainConfig.ZeroHash},
					&core.Checkpoint{Epoch: 1, Root: []byte{}},
					0,
					true,
					0, /* attestation */
					),
			},
			expectedError: fmt.Errorf("attestation did not pass threshold"),
		},
		{
			name: "target epoch invalid",
			blockBody: &core.BlockBody{
				Slot:                 32,
				Attestations:         generateAttestations(
					state,
					128,
					32,
					&core.Checkpoint{Epoch: 0, Root: params.ChainConfig.ZeroHash},
					&core.Checkpoint{Epoch: 5, Root: []byte{}},
					0,
					true,
					0, /* attestation */
				),
			},
			expectedError: fmt.Errorf("taregt not in current/ previous epoch"),
		},
		{
			name: "target slot not in the correct epoch",
			blockBody: &core.BlockBody{
				Slot:                 32,
				Attestations:         generateAttestations(
					state,
					128,
					100,
					&core.Checkpoint{Epoch: 0, Root: params.ChainConfig.ZeroHash},
					&core.Checkpoint{Epoch: 1, Root: []byte{}},
					0,
					true,
					0, /* attestation */
				),
			},
			expectedError: fmt.Errorf("target slot not in the correct epoch"),
		},
		{
			name: "min att. inclusion delay did not pass",
			blockBody: &core.BlockBody{
				Slot:                 32,
				Attestations:         generateAttestations(
					state,
					128,
					33,
					&core.Checkpoint{Epoch: 0, Root: params.ChainConfig.ZeroHash},
					&core.Checkpoint{Epoch: 1, Root: []byte{}},
					0,
					true,
					0, /* attestation */
				),
			},
			expectedError: fmt.Errorf("min att. inclusion delay did not pass"),
		},
		{
			name: "slot to submit att. has passed",
			blockBody: &core.BlockBody{
				Slot:                 32,
				Attestations:         generateAttestations(
					state,
					128,
					0,
					&core.Checkpoint{Epoch: 0, Root: params.ChainConfig.ZeroHash},
					&core.Checkpoint{Epoch: 0, Root: []byte{}},
					0,
					true,
					0, /* attestation */
				),
			},
			expectedError: fmt.Errorf("slot to submit att. has passed"),
		},
		//{ // TODO - complete
		//	name: "committee index out of range",
		//	blockBody: &core.BlockBody{
		//		Slot:                 32,
		//		Attestations:         generateAttestations(
		//			state,
		//			86,
		//			32,
		//			0,
		//			1,
		//			1000000,
		//			true,
		//			0, /* attestation */
		//		),
		//	},
		//	expectedError: fmt.Errorf("slot to submit att. has passed"),
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func (t *testing.T) {
			state := generateTestState(t)
			st := NewStateTransition()

			if test.expectedError != nil {
				require.EqualError(t, st.ProcessBlockAttestations(state, test.blockBody), test.expectedError.Error())
			} else {
				require.NoError(t, st.ProcessBlockAttestations(state, test.blockBody))
			}
		})
	}
}