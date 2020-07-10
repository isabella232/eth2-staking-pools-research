package crypto

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDKG(t *testing.T) {
	InitBLS()

	// generate a polynomial
	poly1 := &Polynomial{
		Degree: 2,
		Coefficients: []bls.Fr{
			frFromInt(6), // free coefficient
			frFromInt(0), // x^1
			frFromInt(1), // x^2
		},
	}
	poly2 := &Polynomial{
		Degree: 2,
		Coefficients: []bls.Fr{
			frFromInt(3), // free coefficient
			frFromInt(0), // x^1
			frFromInt(1), // x^2
		},
	}
	poly3 := &Polynomial{
		Degree: 2,
		Coefficients: []bls.Fr{
			frFromInt(9), // free coefficient
			frFromInt(0), // x^1
			frFromInt(1), // x^2
		},
	}

	dkg := DKG{polynomials:map[uint32]*Polynomial{
		1: poly1,
		2: poly2,
		3: poly3,
	}, degree:2}

	// gete sks
	sks,err := dkg.GroupSecrets([]uint32{
		1,
		2,
		3,
	})
	require.NoError(t, err)
	// verify
	require.Len(t, sks, 3)
	require.Equal(t, "21", sks[1].GetString(10))
	require.Equal(t, "30", sks[2].GetString(10))
	require.Equal(t, "45", sks[3].GetString(10))


	// group pk
	pk,err := dkg.GroupPK(sks)
	require.NoError(t, err)

	// verify
	expectedGroupSecret := &bls.Fr{}
	expectedGroupSecret.SetInt64(18)
	expectedGroupSk := bls.CastToSecretKey(expectedGroupSecret)

	require.Equal(t, expectedGroupSk.GetPublicKey().GetHexString(), pk.GetHexString())
}
