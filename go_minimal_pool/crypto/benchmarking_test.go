package crypto

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBenchmarkingSigning(t *testing.T) {
	InitBLS()

	amount_of_sigs := 1000

	sks := make([]*bls.Fr, amount_of_sigs)
	sigs := make([]*bls.G2, amount_of_sigs)
	pks := make([]*bls.G1, amount_of_sigs)
	for i := 0 ; i < amount_of_sigs ; i++ {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()
		sks[i] = bls.CastFromSecretKey(sk)

		sigs[i] = bls.CastFromSign(sk.Sign("test"))
		pks[i] = bls.CastFromPublicKey(sk.GetPublicKey())
	}
}

func TestBenchmarkingPkReconstruction(t *testing.T) {
	InitBLS()

	size := 1028

	sk := bls.Fr{}
	sk.SetByCSPRNG()
	degree := uint8(size)
	p, err := NewPolynomial(sk, degree)
	require.NoError(t,err)
	err = p.GenerateRandom()
	require.NoError(t,err)

	// get points
	g1s := make([][]interface{}, size)
	for i := 0 ; i < size ; i++ {
		p,err := p.Evaluate(frPointerFromInt(int64(i + 1))) // evaluate from x=1 forward
 		require.NoError(t, err)

		g1s[i] = []interface{}{
			frFromInt(int64(i + 1)),g1FromFr(*p),
		}
	}

	// Interpolate back
	l := NewG1LagrangeInterpolation(g1s)
	res,err := l.interpolate()
	require.NoError(t, err)

	fmt.Printf(res.GetString(10))
}