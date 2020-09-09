package state

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidBlockAndHeader(t *testing.T) {
	state := GenerateState(t)

	// mock header and body
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBody := mocks.NewMockIBlockBody(ctrl)
	mockedBody.EXPECT().GetProposer().Return(uint64(456))
	mockedBody.EXPECT().Root().Return([]byte{1,2,3,4,5}, nil)
	mockedBody.EXPECT().Validate().Return(nil)

	mockedheader := mocks.NewMockIBlockHeader(ctrl)
	mockedheader.EXPECT().GetBlockRoot().Return([]byte{1,2,3,4,5})
	mockedheader.EXPECT().Validate(gomock.Any()).Return(nil)

	require.NoError(t, state.ValidateBlock(mockedheader, mockedBody))
}

func TestInvalidProposer(t *testing.T) {
	state := GenerateState(t)

	// mock header and body
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBody := mocks.NewMockIBlockBody(ctrl)
	mockedBody.EXPECT().GetProposer().Return(uint64(0))

	mockedheader := mocks.NewMockIBlockHeader(ctrl)

	require.EqualError(t, state.ValidateBlock(mockedheader, mockedBody), "block expectedProposer is worng, expected 456 but received 0")
}

func TestMismatchBodyRootSig(t *testing.T) {
	state := GenerateState(t)

	// mock header and body
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBody := mocks.NewMockIBlockBody(ctrl)
	mockedBody.EXPECT().GetProposer().Return(uint64(456))
	mockedBody.EXPECT().Root().Return([]byte{1,2,3,4,5}, nil)

	mockedheader := mocks.NewMockIBlockHeader(ctrl)
	mockedheader.EXPECT().GetBlockRoot().Return([]byte{1,2,3,4,6})

	require.EqualError(t, state.ValidateBlock(mockedheader, mockedBody), "signed block root does not match body root")
}
