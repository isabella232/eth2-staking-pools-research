package crypto

import "github.com/herumi/bls-eth-go-binary/bls"

// responsible for generating shares for redistribution
// https://github.com/bloxapp/eth2-staking-pools-research/blob/master/pool_rotation.md
type Redistribuition struct {
	degree uint8
	originalSk *bls.Fr
	polynomial *Polynomial
}

func NewRedistribuition(degree uint8, originalSk *bls.Fr) (*Redistribuition,error) {
	p, err := NewPolynomial(*originalSk, degree)
	if err != nil {
		return nil, err
	}

	return &Redistribuition{
		degree:     degree,
		originalSk: originalSk,
		polynomial:p,
	}, nil
}

func (distro *Redistribuition)GenerateShares(indexes []uint32) (map[uint32]*bls.Fr, error) {
	ret := make(map[uint32]*bls.Fr)
	for _, share_idx := range indexes {
		share_idx_fr := &bls.Fr{}
		share_idx_fr.SetInt64(int64(share_idx))
		p,err := distro.polynomial.Evaluate(share_idx_fr)
		if err != nil {
			return nil, err
		}

		ret[share_idx] = p
	}

	return ret, nil
}