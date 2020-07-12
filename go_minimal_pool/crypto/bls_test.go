package crypto

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAggregate(t *testing.T) {
	InitBLS()
	secret1 := &bls.Fr{}
	secret1.SetByCSPRNG()
	secret2 := &bls.Fr{}
	secret2.SetByCSPRNG()

	sk1 := bls.CastToSecretKey(secret1)
	pk1 := sk1.GetPublicKey()
	sk2 := bls.CastToSecretKey(secret2)
	pk2 := sk2.GetPublicKey()

	sig1 := sk1.Sign("hello")
	sig2 := sk2.Sign("hello")

	agg_sig := bls.CastToSign(agg_g2(bls.CastFromSign(sig1), bls.CastFromSign(sig2)))
	agg_pk := bls.CastToPublicKey(agg_g1(bls.CastFromPublicKey(pk1),bls.CastFromPublicKey(pk2)))

	require.True(t, agg_sig.Verify(agg_pk, "hello"))
}
