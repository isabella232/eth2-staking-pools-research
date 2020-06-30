package crypto

import (
	"github.com/herumi/bls-eth-go-binary/bls"
)

func initBLS() error {
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
