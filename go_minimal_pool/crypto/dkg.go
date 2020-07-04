package crypto

import "github.com/herumi/bls-eth-go-binary/bls"

/**
	This builds a polynomial for a particular secret and generates shares for distribution
 */
type DKG struct {
	polynomials map[uint32]*Polynomial
	degree uint8
}

func NewDKG(degree uint8, indexes []uint32) (*DKG,error) {
	polynomials := make(map[uint32]*Polynomial)
	for _, idx := range indexes {
		secret := &bls.Fr{}
		secret.SetByCSPRNG()
		p, err := NewPolynomial(*secret, degree)
		if err != nil {
			return nil, err
		}

		polynomials[idx] = p
	}

	return &DKG{polynomials:polynomials, degree:degree}, nil
}

func (dkg *DKG) GroupSecrets(indexes []uint32) (map[uint32]*bls.Fr, error) {
	ret := make(map[uint32][]*bls.Fr)
	for p_idx := range dkg.polynomials {
		poly := dkg.polynomials[p_idx]
		for _, share_idx := range indexes {
			share_idx_fr := &bls.Fr{}
			share_idx_fr.SetInt64(int64(share_idx))
			p,err := poly.Evaluate(share_idx_fr)
			if err != nil {
				return nil, err
			}

			ret[share_idx] = append(ret[share_idx], p)
		}
	}

	return dkg.sumShares(ret), nil
}

func (dkg *DKG) sumShares(shares map[uint32][]*bls.Fr) map[uint32]*bls.Fr {
	ret := make(map[uint32]*bls.Fr)
	for pIdx, shares := range shares {
		pIdx_fr := &bls.Fr{}
		pIdx_fr.SetInt64(int64(pIdx))

		sum := &bls.Fr{}
		sum.SetInt64(0)
		for _, s := range shares {
			bls.FrAdd(sum, sum, s)
		}

		ret[pIdx] = sum
	}

	return ret
}
