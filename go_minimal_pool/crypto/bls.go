package crypto

import (
	"github.com/herumi/bls-eth-go-binary/bls"
)

func InitBLS() error {
	err := bls.Init(bls.BLS12_381)
	if err != nil {
		return err
	}
	err = bls.SetETHmode(bls.EthModeDraft07)
	if err != nil {
		return err
	}
	return nil
}

func Sign(secret *bls.Fr, msg []byte) *bls.G2 {
	sk := bls.CastToSecretKey(secret)
	sig := sk.SignByte(msg)
	return bls.CastFromSign(sig)
}


func agg_g1(a *bls.G1, b *bls.G1) *bls.G1 {
	out := &bls.G1{}
	bls.G1Add(out, a, b)

	return out
}

func agg_g2(a *bls.G2, b *bls.G2) *bls.G2 {
	out := &bls.G2{}
	bls.G2Add(out, a, b)

	return out
}