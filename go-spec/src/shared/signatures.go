package shared

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/prysmaticlabs/go-ssz"
)

func SignBlock(block *core.PoolBlock, sk []byte, domain []byte) (*bls.Sign, error) {
	root, err := BlockSigningRoot(block, domain)
	if err != nil {
		return nil, err
	}

	privKey := bls.SecretKey{}
	err = privKey.Deserialize(sk)
	if err != nil {
		return nil, err
	}
	sig := privKey.SignByte(root[:])
	return sig, nil
}

func VerifyBlockSigningRoot(block *core.PoolBlock, pubKey []byte, sigByts []byte, domain []byte) error {
	root, err := BlockSigningRoot(block, domain)

	// pk
	pk := &bls.PublicKey{}
	err = pk.Deserialize(pubKey)
	if err != nil {
		return err
	}

	// sig
	sig := &bls.Sign{}
	err = sig.Deserialize(sigByts)
	if err != nil {
		return err
	}

	if !sig.VerifyByte(pk, root[:]) {
		return fmt.Errorf("sig not verified")
	}

	return nil
}

func BlockSigningRoot(block *core.PoolBlock, domain []byte) ([32]byte, error) {
	root, err := ssz.HashTreeRoot(block)
	if err != nil {
		return [32]byte{}, err
	}
	container := struct {
		ObjectRoot []byte
		Domain []byte
	}{
		root[:],
		domain,
	}
	return ssz.HashTreeRoot(container)
}