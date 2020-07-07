package crypto

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func frPointerRandom() *bls.Fr {
	p := &bls.Fr{}
	p.SetByCSPRNG()
	return p
}

func TestRedistribuition(t *testing.T) {
	InitBLS()

	tests := []struct{
		testName string
		sks []*bls.Fr
		skStr string
	}{
		{
			testName: "y=5x^2+3",
			sks: []*bls.Fr{
				frPointerFromInt(8),
				frPointerFromInt(23),
				frPointerFromInt(48),
			},
			skStr: "3",
		},
		{
			testName: "y=2x^2+3",
			sks: []*bls.Fr{
				frPointerFromInt(5),
				frPointerFromInt(11),
				frPointerFromInt(21),
			},
			skStr: "3",
		},
		{
			testName: "y=x^2+30",
			sks: []*bls.Fr{
				frPointerFromInt(31),
				frPointerFromInt(34),
				frPointerFromInt(39),
			},
			skStr: "30",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			sk1 := test.sks[0]
			sk2 := test.sks[1]
			sk3 := test.sks[2]

			distro1,err := NewRedistribuition(3, sk1)
			require.NoError(t,err)
			distro2,err := NewRedistribuition(3, sk2)
			require.NoError(t,err)
			distro3,err := NewRedistribuition(3, sk3)
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
			sk1Interpolated,err := p1.interpolate()
			require.NoError(t,err)
			p2 := NewLagrangeInterpolation(shares2)
			sk2Interpolated,err := p2.interpolate()
			require.NoError(t,err)
			p3 := NewLagrangeInterpolation(shares3)
			sk3Interpolated,err := p3.interpolate()
			require.NoError(t,err)

			// reconstruct group secret from the 3 re-distributed shares
			group := [][]bls.Fr {
				{frFromInt(1), *sk1Interpolated},
				{frFromInt(2), *sk2Interpolated},
				{frFromInt(3), *sk3Interpolated},
			}
			p := NewLagrangeInterpolation(group)
			groupSk, err := p.interpolate()
			require.NoError(t,err)

			require.Equal(t, test.skStr, groupSk.GetString(10))
		})
	}
}
