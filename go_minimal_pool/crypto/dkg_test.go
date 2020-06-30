package crypto

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDKG(t *testing.T) {
	initBLS()

	// generate a polynomial
	poly := &Polynomial{
		Degree: 2,
		Coefficients: []bls.Fr{
			frFromInt(6), // free coefficient
			frFromInt(0), // x^1
			frFromInt(1), // x^2
		},
	}

	dkg := DKG{polynomial:poly}

	// gete shares
	shares,err := dkg.shares([]*bls.Fr{
		frPointerFromInt(1),
		frPointerFromInt(2),
		frPointerFromInt(3),
		frPointerFromInt(4),
	})
	require.NoError(t, err)

	// verify
	require.Len(t, shares, 4)
	require.Equal(t, "7", shares["1"].GetString(10))
	require.Equal(t, "10", shares["2"].GetString(10))
	require.Equal(t, "15", shares["3"].GetString(10))
	require.Equal(t, "22", shares["4"].GetString(10))
}
