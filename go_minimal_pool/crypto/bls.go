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
