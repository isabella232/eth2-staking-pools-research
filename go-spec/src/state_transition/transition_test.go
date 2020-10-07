package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandaoSeedMix(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	state := generateTestState(t)
	block := &core.PoolBlock{
		Proposer:        2713,
		Slot:            35,
		Body: &core.PoolBlockBody{
			RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		},
		ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
	}

	require.NoError(t, processRANDAO(state, block))

	newsSeed, err := shared.GetSeed(state, 35)
	require.NoError(t, err)
	require.EqualValues(t, toByte("e4a17401658219365021cf584f4758d4b22ec861d9653e8249c8b4f73285a909"), newsSeed)
}


func TestBlockPostValidation(t *testing.T) {
	tests := []struct{
		name          string
		block         *core.PoolBlock
		proposerId    uint64
		postStateRoot []byte
		expectedError error
	}{
		{
			name: "valid post state root",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId:    2713,
			postStateRoot: toByte("1fbf8ba42600abfd3185f95ed278b4169e141bdb96c0e8e0c6bdf257a8eb1701"),
			expectedError: nil,
		},
		{
			name: "invalid post state root",
			block: &core.PoolBlock{
				Proposer:        2713,
				Slot:            35,
				Body: &core.PoolBlockBody{
					RandaoReveal:         toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
				},
				ParentRoot: toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
			proposerId:    2713,
			postStateRoot: toByte("4ac7911683b0d4643c289cdd3c45bebaa30e912f28d34a2e7cc0009d65273bd8"),
			expectedError: fmt.Errorf("new block state root is wrong"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := generateTestState(t)

			// sign
			sk := []byte(fmt.Sprintf("%d", test.proposerId))
			sig, err := shared.SignBlock(test.block, sk, []byte("domain")) // TODO - dynamic domain

			// header
			signed := &core.SignedPoolBlock{
				Block:                test.block,
				Signature:            sig.Serialize(),
			}

			st := NewStateTransition()
			_, err = st.ExecuteStateTransition(state, signed)
			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
