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

func GenerateValidHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		0,
		6,
		true,
		true,
		true,
		true,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateWrongProposerHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		455, // wrong
		SK,
		"",
		0,
		6,
		true,
		true,
		true,
		true,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateInvalidProposerHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		4550000000, // invalid
		SK,
		"",
		0,
		6,
		true,
		true,
		true,
		true,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateWrongRootHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456, // invalid
		SK,
		"73aa0c267311b8c49f0b9812f7f2f845c55b0d4921c1b40a38f0d82d471d9bcf", // wrong
		0,
		6,
		true,
		true,
		true,
		true,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}


func GenerateInvalidSigHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		"59aaaa8f68aad68552512feb1e27438ddbe2730ea416bb3337b579317610d702", // wrong
		"",
		0,
		6,
		true,
		true,
		true,
		true,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateCreatePoolHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		1,
		6,
		true,
		true,
		false,
		false,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateNotCreatePoolHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		2,
		6,
		true,
		true,
		false,
		false,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateFinalizedAttestationPoolHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		2,
		6,
		true,
		true,
		true,
		false,
		false,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateNotFinalizedAttestationPoolHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		2,
		6,
		false,
		true,
		true,
		false,
		false,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateFinalizedProposalPoolHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		2,
		6,
		true,
		true,
		false,
		true,
		false,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateNotFinalizedProposalPoolHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		2,
		6,
		false,
		false,
		false,
		true,
		false,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateCreatePoolWithExistingIdHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		1,
		4,
		true,
		true,
		false,
		false,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
	)
}

func GenerateTooSmallRandaoHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		0,
		6,
		true,
		true,
		true,
		true,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6"),
	)
}

func GenerateTooBigRandaoHeadAndBody(t *testing.T, state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		t,
		state,
		456,
		SK,
		"",
		0,
		6,
		true,
		true,
		true,
		true,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6ddd"),
	)
}

func generateHeaderAndBody(
	t *testing.T,
	state *core.State,
	proposer uint64,
	skStr string,
	headerBodyRoot string,
	createPoolStatus int32,
	createPoolReqId uint64,
	attestationDutyFinalized bool,
	proposalDutyFinalized bool,
	includeBeaconAttestationDuty bool,
	includeBeaconProposalDuty bool,
	includeCreatePool bool,
	randao []byte,
	) (*core.BlockHeader, *core.BlockBody) {
	body := &core.BlockBody{
		Proposer:           proposer,
		Epoch:              5,
		ExecutionSummaries: []*core.ExecutionSummary{
			&core.ExecutionSummary{
				PoolId:        3,
				Epoch:         5,
				Duties:        []*core.BeaconDuty {},
			},
		},
		NewPoolReq:         []*core.CreateNewPoolRequest{},
		StateRoot:          nil,// will be calculated later toByte("eb4a4ad8124e7c6511326a6951d3c9095077f9d7917750dc779b4beaee24e8e2"),
		ParentBlockRoot:    toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		Randao: randao,
	}

	if includeBeaconAttestationDuty {
		body.ExecutionSummaries[0].Duties = append(body.ExecutionSummaries[0].Duties, &core.BeaconDuty{
			Type:          0, // attestation
			Committee:     12,
			Slot:         342,
			Finalized:     attestationDutyFinalized,
			Participation: []byte{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0},
		})
	}
	if includeBeaconProposalDuty {
		body.ExecutionSummaries[0].Duties = append(body.ExecutionSummaries[0].Duties, &core.BeaconDuty{
			Type:          1, // proposal
			Committee:     0,
			Slot:         343,
			Finalized:     proposalDutyFinalized,
			Participation: []byte{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0},
		})
	}
	if includeCreatePool {
		body.NewPoolReq = append(body.NewPoolReq, &core.CreateNewPoolRequest{
			Id:                  createPoolReqId,
			Status:              createPoolStatus, // started
			StartEpoch:          0,
			EndEpoch:            1,
			LeaderBlockProducer: 0,
			CreatePubKey:        toByte("a3b9110ec26cbb02e6182fab4dcb578d17411f26e41f16aad99cfce51e9bc76ce5e7de00a831bbcadd1d7bc0235c945d"), // priv: 3ef5411174c7d9672652bf4ffc342af3720cc23e52c377b95927871645435f41
			Participation:       []byte{43,12,89,35,99,16,63,13,33,0,1,3,88,12,43,1},
		})
	}

	// calculate and set state root after applying block
	require.NoError(t, CalculateAndInsertStateRootToBlock(state ,body))

	sk := &bls.SecretKey{}
	sk.SetHexString(skStr)

	root, _ := ssz.HashTreeRoot(body)
	if len(headerBodyRoot) > 0 {
		root = shared.SliceToByte32(toByte(headerBodyRoot))
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
			Active:true,
		}
	}

	return &core.State{
		Pools: pools,
		BlockProducers: bps,
		Seeds:          []*core.EpochAndBytes{
			&core.EpochAndBytes{
				Epoch:                0,
				Bytes:                []byte("seedseedseedseedseedseedseedseed"),
			},
		},
		BlockRoots: 	[]*core.EpochAndBytes{},
		// TODO - add state roots
	}
}