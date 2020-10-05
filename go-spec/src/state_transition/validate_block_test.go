package state_transition

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlockPreValidation(t *testing.T) {
	tests := []struct{
		name string
		body *core.BlockBody
		proposerId uint64
		expectedError error
		useCorretBodyRoot bool
	}{
		{
			name: "valid sig",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2713,
			expectedError: nil,
			useCorretBodyRoot: true,
		},
		{
			name: "invalid sig",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2712,
			expectedError: fmt.Errorf("signature did not verify"),
			useCorretBodyRoot: true,
		},
		{
			name: "wrong proposer",
			body: &core.BlockBody{
				Proposer:             2,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2,
			expectedError: fmt.Errorf("block expectedProposer is worng, expected 2713 but received 2"),
			useCorretBodyRoot: true,
		},
		{
			name: "invalid proposer",
			body: &core.BlockBody{
				Proposer:             4550000000,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2,
			expectedError: fmt.Errorf("block expectedProposer is worng, expected 2713 but received 4550000000"),
			useCorretBodyRoot: true,
		},
		{
			name: "invalid body root",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("signed block root does not match body root"),
			useCorretBodyRoot: false,
		},
		{
			name: "RANDAO too small",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("RANDAO should be 32 byte"),
			useCorretBodyRoot: true,
		},
		{
			name: "RANDAO too big",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6ddd"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("RANDAO should be 32 byte"),
			useCorretBodyRoot: true,
		},
		{
			name: "invalid parent block root",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e42"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("parent block root not found"),
			useCorretBodyRoot: true,
		},
		{
			name: "block from the past",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 33,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2713,
			expectedError: fmt.Errorf("new block's parent block root can't be of a future epoch"),
			useCorretBodyRoot: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := generateTestState(t)

			// block root
			root,err := ssz.HashTreeRoot(test.body)
			require.NoError(t, err)

			// sign
			sk := &bls.SecretKey{}
			sk.SetHexString(hex.EncodeToString([]byte(fmt.Sprintf("%d", test.proposerId))))
			sig := sk.SignByte(root[:])

			if !test.useCorretBodyRoot {
				root[0] = 0
			}

			// header
			head := &core.BlockHeader{
				BlockRoot:            root[:],
				StateRoot:            []byte("stateroot"),
				Signature:            sig.Serialize(),
			}

			st := NewStateTransition()
			if test.expectedError != nil {
				require.EqualError(t, st.PreApplyValidateBlock(state, head, test.body), test.expectedError.Error())
			} else {
				require.NoError(t, st.PreApplyValidateBlock(state, head, test.body))
			}
		})
	}
}

func TestBlockPostValidation(t *testing.T) {
	tests := []struct{
		name string
		body *core.BlockBody
		proposerId uint64
		postStateRoot []byte
		expectedError error
	}{
		{
			name: "valid post state root",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2713,
			postStateRoot: toByte("83912d4e1063873194fbc625a38c16e80c6da5abe4035a9d9af6cd342e017170"),
			expectedError: nil,
		},
		{
			name: "invalid post state root",
			body: &core.BlockBody{
				Proposer:             2713,
				Slot:                 35,
				ParentBlockRoot:      toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
				Randao:               toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
			},
			proposerId: 2713,
			postStateRoot: toByte("4ac7911683b0d4643c289cdd3c45bebaa30e912f28d34a2e7cc0009d65273bd8"),
			expectedError: fmt.Errorf("new block state root is wrong"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := generateTestState(t)

			// block root
			root,err := ssz.HashTreeRoot(test.body)
			require.NoError(t, err)

			// sign
			sk := &bls.SecretKey{}
			sk.SetHexString(hex.EncodeToString([]byte(fmt.Sprintf("%d", test.proposerId))))
			sig := sk.SignByte(root[:])

			// header
			head := &core.BlockHeader{
				BlockRoot:            root[:],
				StateRoot:            test.postStateRoot,
				Signature:            sig.Serialize(),
			}

			st := NewStateTransition()
			newState, err := st.ApplyBlock(state, test.body)
			fmt.Printf("%s\n", hex.EncodeToString(shared.GetStateRoot(newState, test.body.Slot)))
			require.NoError(t, err)

			if test.expectedError != nil {
				require.EqualError(t, st.PostApplyValidateBlock(newState, head, test.body), test.expectedError.Error())
			} else {
				require.NoError(t, st.PostApplyValidateBlock(newState, head, test.body))
			}
		})
	}
}