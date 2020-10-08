package state_transition

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/stretchr/testify/require"
	"github.com/ulule/deepcopier"
	"testing"
)

var SK = "59aaaa8f68aad68552512feb1e27438ddbe2730ea416bb3337b579317610d703"
var PK = "846b207d6eb0377ac74db3f7bc295a02340d784431a7cf14dddcd5610c2925facef5763fcaf4358434f92bc9a2906744"

func toByte(str string) []byte {
	ret, _ := hex.DecodeString(str)
	return ret
}

func generateAttestations(
	state *core.State,
	howManyBpSig uint64,
	slot uint64,
	sourceCheckpoint *core.Checkpoint,
	targetCheckpoint *core.Checkpoint,
	committeeIdx uint32,
	finalized bool,
	dutyType int32, // 0 - attestation, 1 - proposal, 2 - aggregation
	) []*core.Attestation {

	data := &core.AttestationData{
		Slot:                 slot,
		CommitteeIndex:       committeeIdx,
		BeaconBlockRoot:      []byte("block root"),
		Source:               sourceCheckpoint,
		Target:               targetCheckpoint,
		ExecutionSummaries:   []*core.ExecutionSummary{
			&core.ExecutionSummary{
				PoolId: 3,
				Epoch:  params.SlotToEpoch(slot),
				Duties:               []*core.BeaconDuty{
					&core.BeaconDuty{
						Type:                 dutyType, // attestation
						Committee:            12,
						Slot:                 342,
						Finalized:            finalized,
						Participation:        []byte{1,3,88},
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

	expectedCommittee, err := shared.SlotCommitteeByIndex(state, data.Slot, uint64(data.CommitteeIndex))
	if err != nil {
		return nil
	}

	var aggregatedSig *bls.Sign
	aggBits := make(bitfield.Bitlist, len(expectedCommittee)) // for bytes
	signed := uint64(0)
	for i, bpId := range expectedCommittee {
		bp := shared.GetBlockProducer(state, bpId)
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

		if signed >= howManyBpSig {
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

func generateTestState(t *testing.T) *core.State {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	pools := make([]*core.Pool, 128)

	// block producers
	bps := make([]*core.BlockProducer, len(pools) * int(params.ChainConfig.VaultSize))
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

	// vaults (pool)
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, params.ChainConfig.VaultSize)
		for j := 0 ; j < int(params.ChainConfig.VaultSize) ; j++ {
			executors[j] = bps[i*int(params.ChainConfig.VaultSize) + j].GetId()
		} // no need to sort as they are already

		sk := &bls.SecretKey{}
		sk.SetHexString(hex.EncodeToString([]byte(fmt.Sprintf("%d", uint64(i)))))

		pools[i] = &core.Pool{
			Id:              uint64(i),
			SortedCommittee: executors,
			PubKey:          sk.GetPublicKey().Serialize(),
			Active:true,
		}
	}

	ret := &core.State {
		CurrentSlot: 0,
		Pools: pools,
		BlockProducers: bps,
		Seeds:          []*core.SlotAndBytes{},
		BlockRoots: 	[]*core.SlotAndBytes{},
		StateRoots: 	[]*core.SlotAndBytes{},
		PreviousEpochAttestations: []*core.PendingAttestation{},
		CurrentEpochAttestations:[]*core.PendingAttestation{},
		JustificationBits: []byte{0},
		CurrentJustifiedCheckpoint: &core.Checkpoint{
			Epoch:                0,
			Root:                 params.ChainConfig.ZeroHash,
		},
	}

	ret, _ = generateAndApplyBlocks(ret, 5)

	return ret
}

// will generate and save blocks from slot 0 until maxBlocks
func generateAndApplyBlocks(state *core.State, maxBlocks int) (*core.State, error) {
	var previousBlockHeader *core.PoolBlockHeader
	for i := 0 ; i < maxBlocks ; i++ {
		// get proposer
		pID, err := shared.BlockProposer(state, uint64(i))
		if err != nil {
			return nil, err
		}

		// state root
		stateRoot,err := ssz.HashTreeRoot(state)
		if err != nil {
			return nil, err
		}

		// parent
		if previousBlockHeader != nil {
			previousBlockHeader.StateRoot =  stateRoot[:]
		}
		parentRoot,err := ssz.HashTreeRoot(previousBlockHeader)
		if err != nil {
			return nil, err
		}

		block := &core.PoolBlock{
			Slot:                 uint64(i),
			Proposer:             pID,
			ParentRoot:           parentRoot[:],
			StateRoot:            params.ChainConfig.ZeroHash, // temp
			Body:                 &core.PoolBlockBody{
				RandaoReveal:         nil,
				Attestations:         []*core.Attestation{},
				NewPoolReq:           []*core.CreateNewPoolRequest{},
			},
		}

		// process
		st := NewStateTransition()

		// compute state root
		root, err := st.ComputeStateRoot(state, &core.SignedPoolBlock{
			Block:                block,
			Signature:            []byte{},
		})
		if err != nil {
			return nil, err
		}
		block.StateRoot = root[:]

		// sign
		sig, err := shared.SignBlock(block, []byte(fmt.Sprintf("%d", pID)), []byte("domain")) // TODO - dynamic domain
		if err != nil {
			return nil, err
		}

		// execute
		state, err = st.ExecuteStateTransition(state, &core.SignedPoolBlock{
			Block:                block,
			Signature:            sig.Serialize(),
		})
		if err != nil {
			return nil, err
		}

		// copy to previousBlockRoot
		previousBlockHeader = &core.PoolBlockHeader{}
		deepcopier.Copy(state.LatestBlockHeader).To(previousBlockHeader)
	}
	return state, nil
}