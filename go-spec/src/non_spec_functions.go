package src

import (
	"encoding/hex"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/state"
	"github.com/herumi/bls-eth-go-binary/bls"
)

/**
	Helper functions which are out of spec
 */
type NonSpecFunctions interface {
	GetBlockBody(root []byte) *BlockBody
	SaveBlockBody(body *BlockBody) error
	FetchExecutedDuties (pubKey *bls.PublicKey, epoch uint64) ([]*BeaconDuty, error)
	WasDutyIncluded (pubKey *bls.PublicKey, epoch uint64, duty *BeaconDuty) (bool, error)
	PoolExecutionStats (poolId uint64, epoch uint64, duty *BeaconDuty) ([16]byte, error)
	SaveState(state *state.State, epoch uint64) error
	GetState(epoch uint64) *state.State
	SeedForEpoch(epoch uint64) [32]byte
}

type SimpleFunctions struct {
	blockBodies map[string]*BlockBody
	states map[uint64]*state.State
}

func NewSimpleFunctions() *SimpleFunctions {
	return &SimpleFunctions{
		blockBodies: make(map[string]*BlockBody),
		states: make(map[uint64]*state.State),
	}
}

func (s *SimpleFunctions) GetBlockBody(root []byte) *BlockBody {
	return s.blockBodies[hex.EncodeToString(root)]
}

func (s *SimpleFunctions) SaveBlockBody(body *BlockBody) error {
	root, err := body.Root()
	if err != nil {
		return err
	}
	s.blockBodies[hex.EncodeToString(root)] = body
	return nil
}

func (s *SimpleFunctions) FetchExecutedDuties (pubKey *bls.PublicKey, epoch uint64) ([]*BeaconDuty, error) {
	return nil, nil
}

func (s *SimpleFunctions) WasDutyIncluded (pubKey *bls.PublicKey, epoch uint64, duty *BeaconDuty) (bool, error) {
	return false, nil
}

func (s *SimpleFunctions) PoolExecutionStats (poolId uint64, epoch uint64, duty *BeaconDuty) ([16]byte, error) {
	return [16]byte{}, nil
}

func (s *SimpleFunctions) SaveState(state *state.State, epoch uint64) error {
	s.states[epoch] = state
	return nil
}

func (s *SimpleFunctions) GetState(epoch uint64) *state.State {
	return s.states[epoch]
}

func (s *SimpleFunctions) SeedForEpoch(epoch uint64) [32]byte {
	state := s.GetState(epoch)
	return state.seed
}