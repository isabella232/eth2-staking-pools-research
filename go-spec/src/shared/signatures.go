package shared

import (
	"encoding/binary"
	"encoding/hex"
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
	err = privKey.SetHexString(hex.EncodeToString(sk))
	if err != nil {
		return nil, err
	}
	sig := privKey.SignByte(root[:])
	return sig, nil
}

func VerifyBlockSigningRoot(block *core.PoolBlock, pubKey []byte, sigByts []byte, domain []byte) error {
	root, err := BlockSigningRoot(block, domain)
	if err != nil {
		return err
	}

	res, err := VerifySignature(root[:], pubKey, sigByts)
	if err != nil {
		return err
	}
	if !res {
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

func RANDAOSigningData(state *core.State) (data []byte, domain []byte, err error)  {
	epoch := GetCurrentEpoch(state)
	data = make([]byte, 8) // 64 bit
	binary.LittleEndian.PutUint64(data, epoch)

	domain = []byte("randao") // TODO - change to dynamic domain

	return data, domain, nil
}

func RandaoSigningRoot(data []byte, domain []byte) ([32]byte, error) {
	container := struct {
		ObjectRoot []byte
		Domain []byte
	}{
		data,
		domain,
	}
	return ssz.HashTreeRoot(container)
}

func SignRandao(data []byte, domain []byte, sk []byte) (*bls.Sign, error) {
	root, err := RandaoSigningRoot(data, domain)
	if err != nil {
		return nil, err
	}

	privKey := bls.SecretKey{}
	err = privKey.SetHexString(hex.EncodeToString(sk))
	if err != nil {
		return nil, err
	}
	sig := privKey.SignByte(root[:])
	return sig, nil
}

func VerifyRandaoRevealSignature(data []byte, domain []byte, pubKey []byte, sigByts []byte) (bool, error)  {
	root, err := RandaoSigningRoot(data, domain)
	if err != nil {
		return false, err
	}
	return VerifySignature(root[:], pubKey, sigByts)
}

func VerifySignature(root []byte, pubKey []byte, sigByts []byte) (bool, error) {
	// pk
	pk := &bls.PublicKey{}
	err := pk.Deserialize(pubKey)
	if err != nil {
		return false, err
	}

	// sig
	sig := &bls.Sign{}
	err = sig.Deserialize(sigByts)
	if err != nil {
		return false, err
	}

	// verify
	if !sig.VerifyByte(pk, root) {
		return false, nil
	}
	return true, nil
}