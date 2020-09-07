package src

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateRandomState() *State {
	pools := make([]*Pool, 5)

	//
	bps := make([]*BlockProducer, len(pools) * int(TestConfig().PoolExecutorsNumber))
	for i := 0 ; i < len(bps) ; i++ {
		bps[i] = &BlockProducer{
			Id:      uint64(i),
			Balance: 1000,
			Stake:   0,
			Slashed: false,
		}
	}

	//
	for i := 0 ; i < len(pools) ; i++ {
		executors := make([]uint64, TestConfig().PoolExecutorsNumber)
		for j := 0 ; j < int(TestConfig().PoolExecutorsNumber) ; j++ {
			executors[j] = bps[i*int(TestConfig().PoolExecutorsNumber) + j].Id
		} // no need to sort as they are already

		pools[i] = &Pool{
			Id:              uint64(i),
			SortedExecutors: executors,
		}
	}

	return &State{
		Pools:           pools,
		BlockRoots:      nil,
		HeadBlockHeader: nil,
		BlockProducers:  bps,
		Seed:            []byte("seed"),
	}
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

	state := GenerateRandomState()

	// set BP
	bp, err := state.GetBlockProducer(0)
	require.NoError(t, err)
	bp.PubKey = sk.GetPublicKey()

	require.NoError(t, state.ValidateBlock(head, body))
}