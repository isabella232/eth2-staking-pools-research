package state_transition

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/stretchr/testify/require"
	"testing"
)

var SK = "59aaaa8f68aad68552512feb1e27438ddbe2730ea416bb3337b579317610d703"
var PK = "846b207d6eb0377ac74db3f7bc295a02340d784431a7cf14dddcd5610c2925facef5763fcaf4358434f92bc9a2906744"

func toByte(str string) []byte {
	ret, _ := hex.DecodeString(str)
	return ret
}

func GenerateValidHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		35,
		17,
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		generateAttestations(state,128, 35,0,true, 0 /* attestation */),
	)
}

func GenerateWrongProposerHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		455, // wrong
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateInvalidProposerHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		4550000000, // invalid
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateWrongRootHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"73aa0c267311b8c49f0b9812f7f2f845c55b0d4921c1b40a38f0d82d471d9bcf", // wrong
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}


func GenerateInvalidSigHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		"59aaaa8f68aad68552512feb1e27438ddbe2730ea416bb3337b579317610d702", // wrong
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateCreatePoolHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		1,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateNotCreatePoolHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		2,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateFinalizedAttestationPoolHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		2,
		129,
		false,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		generateAttestations(state,128, 0,0,true, 0 /* attestation */),
	)
}

func GenerateNotFinalizedAttestationPoolHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		2,
		129,
		false,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		generateAttestations(state,128, 0,0,false, 0 /* attestation */),
	)
}

func GenerateFinalizedProposalPoolHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		2,
		129,
		false,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		generateAttestations(state,128, 0,0,true, 1 /* proposal */),
	)
}

func GenerateNotFinalizedProposalPoolHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		2,
		129,
		false,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		generateAttestations(state,128, 0,0,false, 1 /* proposal */),
	)
}

func GenerateCreatePoolWithExistingIdHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		1,
		4,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateTooSmallRandaoHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateTooBigRandaoHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6ddd"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateInvalidParentBlockRootHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e30"),
		nil,
	)
}

func GenerateInvalidBlockEpochRootHeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		12,
		17,
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateInvalidETH1HeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func GenerateInvalidETH2HeadAndBody(state *core.State)(*core.BlockHeader, *core.BlockBody) {
	return generateHeaderAndBody(
		state,
		33,
		17,
		SK,
		"",
		0,
		129,
		true,
		toByte("97c4116516e77c522344aa3c3c223db0c14bad05aa005be63aadd19341e0cc6d"),
		toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
		nil,
	)
}

func generateAttestations(
	state *core.State,
	howManyBpSig uint64,
	slot uint64,
	committeeIdx uint32,
	finalized bool,
	dutyType int32, // 0 - attestation, 1 - proposal, 2 - aggregation
	) []*core.Attestation {

	data := &core.AttestationData{
		Slot:                 slot,
		CommitteeIndex:       committeeIdx,
		BeaconBlockRoot:      []byte("block root"),
		Source:               &core.Checkpoint{
			Epoch:                0,
			Root:                 []byte{},
		},
		Target:               &core.Checkpoint{
			Epoch:                1,
			Root:                 []byte{},
		},
		ExecutionSummaries:   []*core.ExecutionSummary{
			&core.ExecutionSummary{
				PoolId:               3,
				Epoch:                5,
				Duties:               []*core.BeaconDuty{
					&core.BeaconDuty{
						Type:                 dutyType, // attestation
						Committee:            12,
						Slot:                 342,
						Finalized:            finalized,
						Participation:        []byte{1,3,88,12,43,12,89,35,1,0,99,16,63,13,33,0},
					},
				},
			},
		},
	}

	// sign
	root, err := ssz.HashTreeRoot(data)
	if err != nil {
		return nil
	}

	expectedCommittee, err := shared.SlotCommittee(state, data.Slot, uint64(data.CommitteeIndex))
	if err != nil {
		return nil
	}

	var aggregatedSig *bls.Sign
	aggBits := make(bitfield.Bitlist, core.TestConfig().MaxAttestationCommitteeSize) // for bytes
	signed := uint64(0)
	for i, bpId := range expectedCommittee {
		bp := core.GetBlockProducer(state, bpId)
		sk := &bls.SecretKey{}
		sk.SetHexString(hex.EncodeToString([]byte(fmt.Sprintf("%d", bp.Id))))

		// sign
		if aggregatedSig == nil {
			aggregatedSig = sk.SignByte(root[:])
		} else {
			aggregatedSig.Add(sk.SignByte(root[:]))
		}
		aggBits.SetBitAt(uint64(i), true)
		signed ++

		if signed > howManyBpSig {
			break
		}
	}

	return []*core.Attestation{
		{
			Data:            data,
			Signature:       aggregatedSig.Serialize(),
			AggregationBits: aggBits,
		},
	}
}

func generateHeaderAndBody(
	state *core.State,
	slot uint64,
	proposer uint64,
	skStr string,
	headerBodyRoot string,
	createPoolStatus uint32,
	createPoolReqId uint64,
	includeCreatePool bool,
	randao []byte,
	parentBlockRoot []byte,
	attestations []*core.Attestation,
	) (*core.BlockHeader, *core.BlockBody) {
	body := &core.BlockBody{
		Proposer:           proposer,
		Slot:               slot,
		ParentBlockRoot:    parentBlockRoot,
		Randao: randao,
	}

	if includeCreatePool {
		body.NewPoolReq = append(body.NewPoolReq, &core.CreateNewPoolRequest{
			Id:                  createPoolReqId,
			Status:              createPoolStatus, // started
			StartEpoch:          1,
			EndEpoch:            2,
			LeaderBlockProducer: 1,
			CreatePubKey:        toByte("a3b9110ec26cbb02e6182fab4dcb578d17411f26e41f16aad99cfce51e9bc76ce5e7de00a831bbcadd1d7bc0235c945d"), // priv: 3ef5411174c7d9672652bf4ffc342af3720cc23e52c377b95927871645435f41
			Participation:       []byte{43,12,89,35,99,16,63,13,33,0,1,3,88,12,43,1},
		})
	}

	if attestations != nil {
		body.Attestations = attestations
	}

	// calculate and set state root after applying block
	stateRoot, err := CalculateAndInsertStateRootToBlock(state ,body)
	if err != nil {
		return nil, nil
	}

	sk := &bls.SecretKey{}
	sk.SetHexString(skStr)

	root, _ := ssz.HashTreeRoot(body)
	if len(headerBodyRoot) > 0 {
		root = shared.SliceToByte32(toByte(headerBodyRoot))
	}
	sig := sk.SignByte(root[:])

	return &core.BlockHeader{
		BlockRoot:            root[:],
		StateRoot: 			  stateRoot,
		Signature:            sig.Serialize(),
	}, body
}

func generateTestState(t *testing.T) *core.State {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	pools := make([]*core.Pool, 128)

	//
	bps := make([]*core.BlockProducer, len(pools) * int(core.TestConfig().VaultSize))
	for i := 0 ; i < len(bps) ; i++ {
		sk := &bls.SecretKey{}
		sk.SetHexString(hex.EncodeToString([]byte(fmt.Sprintf("%d", uint64(i)))))

		if i == 17 { // is the block producer for this state
			sk.SetHexString(SK)
		}

		bps[i] = &core.BlockProducer{
			Id:      uint64(i),
			CDTBalance: 1000,
			Stake:   0,
			Slashed: false,
			Active:  true,
			PubKey:  sk.GetPublicKey().Serialize(),
		}
	}

	//
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, core.TestConfig().VaultSize)
		for j := 0 ; j < int(core.TestConfig().VaultSize) ; j++ {
			executors[j] = bps[i*int(core.TestConfig().VaultSize) + j].GetId()
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

	ret := &core.State {
		CurrentEpoch:1,
		Pools: pools,
		BlockProducers: bps,
		Seeds:          []*core.EpochAndBytes{
			&core.EpochAndBytes{
				Epoch:                0,
				Bytes:                []byte("seedseedseedseedseedseedseedseed"),
			},
		},
		BlockRoots: 	[]*core.EpochAndBytes{
			&core.EpochAndBytes{
				Epoch:                0,
				Bytes:                toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e40"),
			},
		},
		StateRoots: 	[]*core.EpochAndBytes{
			&core.EpochAndBytes{
				Epoch:                0,
				Bytes:                toByte("75141b2e032f1b045ab9c7998dfd7238044e40eed0b2c526c33340643e871e41"),
			},
		},
	}

	root, err := ssz.HashTreeRoot(ret)
	if err != nil {
		return nil
	}

	ret.StateRoots = append(ret.StateRoots, &core.EpochAndBytes{
		Epoch:                0,
		Bytes:                root[:],
	})

	return ret
}