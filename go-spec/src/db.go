package src

import (
	"fmt"
	"reflect"
)

func GetBlockBody(root []byte) (*BlockBody, error) {
	return nil, nil
}

func GetPool(state *State, id []byte) (*Pool, error) {
	for _, p := range state.Pools {
		if reflect.DeepEqual(p.Id, id) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("pool not found")
}

func GetBlockProducer(state *State, id uint64) (*BlockProducer, error) {
	for _, bp := range state.BlockProducers {
		if bp.Id == id {
			return bp, nil
		}
	}
	return nil, fmt.Errorf("block producer not found")
}