package block

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/mocks"
	"github.com/golang/mock/gomock"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidSig(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()

	root := []byte("root root root root root root root root root root")

	// mock BP
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBP := mocks.NewMockIBlockProducer(ctrl)
	mockedBP.EXPECT().GetPubKey().Return(sk1.GetPublicKey(), nil)

	header := &BlockHeader{
		BlockRoot: root,
		Signature: sk1.SignByte(root).Serialize(),
	}

	require.NoError(t, header.Validate(mockedBP))
}

func TestBPReturnErrorOnGetPubKey(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()

	wrongSk := &bls.SecretKey{}
	wrongSk.SetByCSPRNG()

	root := []byte("root root root root root root root root root root")

	// mock BP
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBP := mocks.NewMockIBlockProducer(ctrl)
	mockedBP.EXPECT().GetPubKey().Return(nil, fmt.Errorf("no PK"))

	header := &BlockHeader{
		BlockRoot: root,
		Signature: wrongSk.SignByte(root).Serialize(),
	}

	require.EqualError(t, header.Validate(mockedBP), "no PK")
}

func TestBPNoPubKey(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()

	wrongSk := &bls.SecretKey{}
	wrongSk.SetByCSPRNG()

	root := []byte("root root root root root root root root root root")

	// mock BP
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBP := mocks.NewMockIBlockProducer(ctrl)
	mockedBP.EXPECT().GetPubKey().Return(nil, nil)

	header := &BlockHeader{
		BlockRoot: root,
		Signature: wrongSk.SignByte(root).Serialize(),
	}

	require.EqualError(t, header.Validate(mockedBP), "signature did not verify")
}

func TestInvalidSig(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()

	wrongSk := &bls.SecretKey{}
	wrongSk.SetByCSPRNG()

	root := []byte("root root root root root root root root root root")

	// mock BP
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBP := mocks.NewMockIBlockProducer(ctrl)
	mockedBP.EXPECT().GetPubKey().Return(sk1.GetPublicKey(), nil)

	header := &BlockHeader{
		BlockRoot: root,
		Signature: wrongSk.SignByte(root).Serialize(),
	}

	require.EqualError(t, header.Validate(mockedBP), "signature did not verify")
}
