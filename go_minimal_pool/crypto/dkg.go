package crypto

import "github.com/herumi/bls-eth-go-binary/bls"

/**
	This builds a polynomial for a particular secret and generates shares for distribution
 */
type DKG struct {
	polynomials map[*bls.Fr]*Polynomial
	degree uint8
}

func NewDKG(degree uint8, indexes []*bls.Fr) (*DKG,error) {
	polynomials := make(map[*bls.Fr]*Polynomial)
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

func (dkg *DKG) GroupSecrets(indexes []*bls.Fr) (map[*bls.Fr]*bls.Fr, error) {
	ret := make(map[*bls.Fr][]*bls.Fr)
	for p_idx := range dkg.polynomials {
		poly := dkg.polynomials[p_idx]
		for _, share_idx := range indexes {
			p,err := poly.Evaluate(share_idx)
			if err != nil {
				return nil, err
			}

			ret[share_idx] = append(ret[share_idx], p)
		}
	}

	return dkg.sumShares(ret), nil
}

func (dkg *DKG) sumShares(shares map[*bls.Fr][]*bls.Fr) map[*bls.Fr]*bls.Fr {
	ret := make(map[*bls.Fr]*bls.Fr)
	for pIdx, shares := range shares {
		sum := &bls.Fr{}
		sum.SetInt64(0)
		for _, s := range shares {
			bls.FrAdd(sum, sum, s)
		}

		ret[pIdx] = sum
	}

	return ret
}
