package src

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateRandomState(t *testing.T) *State {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	pools := make([]*Pool, 5)

	//
	bps := make([]*BlockProducer, len(pools) * int(TestConfig().PoolExecutorsNumber))
	for i := 0 ; i < len(bps) ; i++ {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		bps[i] = &BlockProducer{
			Id:      uint64(i),
			Balance: 1000,
			Stake:   0,
			Slashed: false,
			PubKey: sk.GetPublicKey(),
		}
	}

	//
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, TestConfig().PoolExecutorsNumber)
		for j := 0 ; j < int(TestConfig().PoolExecutorsNumber) ; j++ {
			executors[j] = bps[i*int(TestConfig().PoolExecutorsNumber) + j].Id
		} // no need to sort as they are already

		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		pools[i] = &Pool{
			Id:              uint64(i),
			SortedExecutors: executors,
			PubKey: sk.GetPublicKey(),
		}
	}

	return &State{
		Pools:           pools,
		BlockRoots:      nil,
		HeadBlockHeader: &BlockHeader{
			BlockRoot: nil,
			Signature: nil,
		},
		BlockProducers:  bps,
		Seed:            SliceToByte32([]byte("seedseedseedseedseedseedseedseed")),
	}
}

func TestRandaoSeedMix(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	state := GenerateRandomState(t)

	body := &BlockBody{
		Proposer:              0,
		Number:                1,
		PoolsExecutionSummary: make([]*PoolExecutionSummary, 0),
		StateRoot:             []byte("state"),
		ParentBlockRoot:       []byte("parent"),
	}
	header,err := NewBlockHeader(sk, body)
	require.NoError(t, err)

	helperFunc = NewSimpleFunctions()
	newState, err := state.ProcessNewBlock(header, body)
	require.NoError(t, err)

	expectedSeed,err := MixSeed(state.Seed, SliceToByte32(header.Signature))
	require.NoError(t, err)
	require.EqualValues(t, expectedSeed, newState.Seed)
}


func TestBlockValidation(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))
	require.NoError(t, bls.SetETHmode(bls.EthModeDraft07))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	// generate header and body
	body := &BlockBody{
		Proposer:              0,
		PoolsExecutionSummary: []*PoolExecutionSummary{GenerateAttestationSuccessfulSummary()},
		StateRoot:             []byte("root"),
		ParentBlockRoot:       []byte("parent"),
	}
	root,err := ssz.HashTreeRoot(body)
	require.NoError(t, err)

	// sign
	sig := sk.SignByte(root[:])

	head := &BlockHeader{
		BlockRoot: root[:],
		Signature: sig.Serialize(),
	}

	state := GenerateRandomState(t)

	// set BP
	bp, err := state.GetBlockProducer(0)
	require.NoError(t, err)
	bp.PubKey = sk.GetPublicKey()

	require.NoError(t, state.ValidateBlock(head, body))
}