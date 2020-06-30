package crypto

import "github.com/herumi/bls-eth-go-binary/bls"

/**
	This builds a polynomial for a particular secret and generates shares for distribution
 */
type DKG struct {
	polynomial *Polynomial
}

func NewDKG(secret bls.Fr, degree uint8) (*DKG,error) {
	p, err := NewPolynomial(secret, degree)
	if err != nil {
		return nil, err
	}

	return &DKG{polynomial:p}, nil
}

func (dkg *DKG) shares(indexes []*bls.Fr) (map[string]*bls.Fr, error) {
	ret := make(map[string]*bls.Fr)
	for _, idx := range indexes {
		p,err := dkg.polynomial.Evaluate(idx)
		if err != nil {
			return nil, err
		}

		ret[idx.GetString(10)] = p
	}

	return ret, nil
}
