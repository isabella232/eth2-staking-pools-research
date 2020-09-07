package src

import (
	"github.com/herumi/bls-eth-go-binary/bls"
)

/**
	Helper functions which are out of spec
 */
type NonSpecFunctions interface {
	GetBlockBody(root []byte) (*BlockBody, error)
	FetchExecutedDuties (pubKey *bls.PublicKey, epoch uint64) ([]*BeaconDuty, error)
	WasDutyIncluded (pubKey *bls.PublicKey, epoch uint64, duty *BeaconDuty) (bool, error)
	PoolExecutionStats (poolId uint64, epoch uint64, duty *BeaconDuty) ([16]byte, error)
}