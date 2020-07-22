package crypto

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

// will generate a shared secret via dkg, then re-distributes the shares and checks
// that shares from the original DKG can't be combined with the shares from the re-distribution to
// re-construct the secret.
func TestShuffleIntegrity(t *testing.T) {
	InitBLS()

	///
	///	Step 1 - first DKG
	///
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
	// shared secret is 18

	dkg := DKG{polynomials:map[uint32]*Polynomial{
		1: poly1,
		2: poly2,
		3: poly3,
	}, degree:2}

	// dkg shares
	sks,err := dkg.GroupSecrets([]uint32{
		1,
		2,
		3,
	})
	require.NoError(t, err)

	///
	/// Step 2 - redistribute shares
	///
	distro1,err := NewRedistribuition(2, sks[1])
	require.NoError(t,err)
	distro2,err := NewRedistribuition(2, sks[2])
	require.NoError(t,err)
	distro3,err := NewRedistribuition(2, sks[3])
	require.NoError(t,err)

	// generate shares to everyone
	sharesFrom1,err := distro1.GenerateShares([]uint32{1,2,3})
	require.NoError(t,err)
	sharesFrom2,err := distro2.GenerateShares([]uint32{1,2,3})
	require.NoError(t,err)
	sharesFrom3,err := distro3.GenerateShares([]uint32{1,2,3})
	require.NoError(t,err)

	shares1 := [][]bls.Fr {
		{frFromInt(1), *sharesFrom1[1]},
		{frFromInt(2), *sharesFrom2[1]},
		{frFromInt(3), *sharesFrom3[1]},
	}
	shares2 := [][]bls.Fr {
		{frFromInt(1), *sharesFrom1[2]},
		{frFromInt(2), *sharesFrom2[2]},
		{frFromInt(3), *sharesFrom3[2]},
	}
	shares3 := [][]bls.Fr {
		{frFromInt(1), *sharesFrom1[3]},
		{frFromInt(2), *sharesFrom2[3]},
		{frFromInt(3), *sharesFrom3[3]},
	}

	// reconstruct individual group sk for 1,2,3
	p1 := NewLagrangeInterpolation(shares1)
	sk1Interpolated,err := p1.Interpolate()
	require.NoError(t,err)
	p2 := NewLagrangeInterpolation(shares2)
	sk2Interpolated,err := p2.Interpolate()
	require.NoError(t,err)
	p3 := NewLagrangeInterpolation(shares3)
	sk3Interpolated,err := p3.Interpolate()
	require.NoError(t,err)


	///
	/// Step 3 - verify redistribute and original shares reconstruct to the same secreet
	///
	group := [][]bls.Fr {
		{frFromInt(1), *sk1Interpolated},
		{frFromInt(2), *sk2Interpolated},
		{frFromInt(3), *sk3Interpolated},
	}
	p := NewLagrangeInterpolation(group)
	redistribuitedGroupSk, err := p.Interpolate()
	require.NoError(t,err)

	require.Equal(t, "18", redistribuitedGroupSk.GetString(10))


	///
	/// Step 4 - try to re-construct the secreet by combining shares from the 2 groups
	///
	t.Run("shuffled 1", func(t *testing.T) {
		group = [][]bls.Fr {
			{frFromInt(1), *sk1Interpolated},
			{frFromInt(2), *sks[2]},
			{frFromInt(3), *sk3Interpolated},
		}
		p = NewLagrangeInterpolation(group)
		shuffledGroupSk, err := p.Interpolate()
		require.NoError(t,err)
		require.NotEqual(t, "18", shuffledGroupSk.GetString(10))
	})
	t.Run("shuffled 2", func(t *testing.T) {
		group = [][]bls.Fr {
			{frFromInt(1), *sk1Interpolated},
			{frFromInt(2), *sks[2]},
			{frFromInt(3), *sks[3]},
		}
		p = NewLagrangeInterpolation(group)
		shuffledGroupSk, err := p.Interpolate()
		require.NoError(t,err)
		require.NotEqual(t, "18", shuffledGroupSk.GetString(10))
	})
	t.Run("shuffled 2", func(t *testing.T) {
		group = [][]bls.Fr {
			{frFromInt(1), *sk1Interpolated},
			{frFromInt(2), *sk2Interpolated},
			{frFromInt(3), *sks[3]},
		}
		p = NewLagrangeInterpolation(group)
		shuffledGroupSk, err := p.Interpolate()
		require.NoError(t,err)
		require.NotEqual(t, "18", shuffledGroupSk.GetString(10))
	})
}
