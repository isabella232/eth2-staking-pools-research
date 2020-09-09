package state_transition

import (
	"encoding/hex"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/stretchr/testify/require"
	"testing"
)

var SK = "59aaaa8f68aad68552512feb1e27438ddbe2730ea416bb3337b579317610d703"

func toByte(str string) []byte {
	ret, _ := hex.DecodeString(str)
	return ret
}

func GenerateValidHeadAndBody(t *testing.T)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
			t,
		456,
		SK,
		"",
		)
}

func GenerateWrongProposerHeadAndBody(t *testing.T)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		455, // wrong
		SK,
		"",
	)
}

func GenerateInvalidProposerHeadAndBody(t *testing.T)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		4550000000, // invalid
		SK,
		"",
	)
}

func GenerateWrongRootHeadAndBody(t *testing.T)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		456, // invalid
		SK,
		"73aa0c267311b8c49f0b9812f7f2f845c55b0d4921c1b40a38f0d82d471d9bcf", // wrong
	)
}


func GenerateInvalidSigHeadAndBody(t *testing.T)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		456,
		"59aaaa8f68aad68552512feb1e27438ddbe2730ea416bb3337b579317610d702", // wrong
		"",
	)
}

func generateHeaderAndBody(t *testing.T, proposer uint64, skStr string, headerBodyRoot string) (*core.BlockHeader, *core.BlockBody) {
	body := &core.BlockBody{
		Proposer:           proposer,
		Epoch:              5,
		ExecutionSummaries: []*core.ExecutionSummary{
			&core.ExecutionSummary{
				PoolId:        12,
				Epoch:         5,
				Duties:        []*core.BeaconDuty {
					&core.BeaconDuty{
						Type:          0, // attestation
						Committee:     12,
						Slot:         342,
						Finalized:     true,
						Participation: []byte{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0},
					},
					&core.BeaconDuty{
						Type:          1, // proposal
						Committee:     0,
						Slot:         343,
						Finalized:     true,
						Participation: []byte{},
					},
				},
			},
		},
		NewPoolReq:         []*core.CreateNewPoolRequest{
			&core.CreateNewPoolRequest{
				Id:                  3,
				Status:              0, // started
				StartEpoch:          5,
				EndEpoch:            6,
				LeaderBlockProducer: 15,
				CreatePubKey:        toByte("public key"),
				Participation:       []byte{43,12,89,35,99,16,63,13,33,0,1,3,88,12,43,1},
			},
		},
		StateRoot:          toByte("state root state root state root state root state root"),
		ParentBlockRoot:    toByte("parent block root parent block root parent block root parent block root"),
	}

	sk := &bls.SecretKey{}
	sk.SetHexString(skStr)

	root, _ := ssz.HashTreeRoot(body)
	if len(headerBodyRoot) > 0 {
		root = shared.SliceToByte32( toByte(headerBodyRoot))
	}
	sig := sk.SignByte(root[:])

	return &core.BlockHeader{
		BlockRoot:            root[:],
		Signature:            sig.Serialize(),
	}, body
}

func generateTestState(t *testing.T) *core.State {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	pools := make([]*core.Pool, 5)

	//
	bps := make([]*core.BlockProducer, len(pools) * int(core.TestConfig().PoolExecutorsNumber))
	for i := 0 ; i < len(bps) ; i++ {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		if i == 456 { // is the block producer for this state
			sk.SetHexString(SK)
		}

		bps[i] = &core.BlockProducer{
			Id:      uint64(i),
			Balance: 1000,
			Stake:   0,
			Slashed: false,
			Active:  true,
			PubKey:  sk.GetPublicKey().Serialize(),
		}
	}

	//
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, core.TestConfig().PoolExecutorsNumber)
		for j := 0 ; j < int(core.TestConfig().PoolExecutorsNumber) ; j++ {
			executors[j] = bps[i*int(core.TestConfig().PoolExecutorsNumber) + j].GetId()
		} // no need to sort as they are already

		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		pools[i] = &core.Pool{
			Id:              uint64(i),
			SortedCommittee: executors,
			PubKey:          sk.GetPublicKey().Serialize(),
		}
	}

	return &core.State{
		Pools: pools,
		BlockProducers: bps,
		Seeds:          [][]byte{[]byte("seedseedseedseedseedseedseedseed")},
	}
}
