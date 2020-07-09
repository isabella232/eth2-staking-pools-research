package crypto

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func g1FromFr(fr bls.Fr) *bls.G1 {
	sk := &bls.SecretKey{}
	sk.SetDecString(fr.GetString(10))

	pk := sk.GetPublicKey()
	return bls.CastFromPublicKey(pk)
}

func g2FromFr(fr bls.Fr) *bls.G2 {
	sk := &bls.SecretKey{}
	sk.SetDecString(fr.GetString(10))

	sig := sk.Sign("test")
	return bls.CastFromSign(sig)
}

func TestG2Interpolation(t *testing.T) {
	InitBLS()

	// generate a polynomial
	poly := Polynomial{
		Degree:              2,
		Coefficients:        []bls.Fr{
			frFromInt(6), // free coefficient
			frFromInt(1), // x^2
			frFromInt(0), // x^1
		},
	}

	// evaluate it and transform to g1 results
	res1,err := poly.Evaluate(frPointerFromInt(int64(1)))
	require.NoError(t, err)
	res2,err := poly.Evaluate(frPointerFromInt(int64(2)))
	require.NoError(t, err)
	res3,err := poly.Evaluate(frPointerFromInt(int64(3)))
	require.NoError(t, err)
	p1 := g2FromFr(*res1)
	p2 := g2FromFr(*res2)
	p3 := g2FromFr(*res3)

	// Interpolate back to a polynomial
	points := [][]interface{} {
		{frFromInt(1),p1},
		{frFromInt(2),p2},
		{frFromInt(3),p3},
	}
	l := NewG2LagrangeInterpolation(points)
	res,err := l.interpolate()
	require.NoError(t, err)

	// compare results
	expected := g2FromFr(frFromInt(6))
	require.Equal(t, expected.GetString(10), res.GetString(10))
}

func TestG1Interpolation(t *testing.T) {
	InitBLS()

	// generate a polynomial
	poly := Polynomial{
		Degree:              2,
		Coefficients:        []bls.Fr{
			frFromInt(6), // free coefficient
			frFromInt(0), // x^1
			frFromInt(1), // x^2
		},
	}

	// evaluate it and transform to g1 results
	res1,err := poly.Evaluate(frPointerFromInt(int64(1)))
	require.NoError(t, err)
	res2,err := poly.Evaluate(frPointerFromInt(int64(2)))
	require.NoError(t, err)
	res3,err := poly.Evaluate(frPointerFromInt(int64(3)))
	require.NoError(t, err)
	p1 := g1FromFr(*res1)
	p2 := g1FromFr(*res2)
	p3 := g1FromFr(*res3)

	// Interpolate back to a polynomial
	points := [][]interface{} {
		{frFromInt(1),p1},
		{frFromInt(2),p2},
		{frFromInt(3),p3},
	}
	l := NewG1LagrangeInterpolation(points)
	res,err := l.interpolate()
	require.NoError(t, err)

	// compare results
	expected := g1FromFr(frFromInt(6))
	require.Equal(t, expected.GetString(10), res.GetString(10))
}