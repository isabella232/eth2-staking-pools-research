package shared

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared/params"
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
		return fmt.Errorf("block sig not verified")
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

// Spec pseudocode definition:
//  def get_domain(state: BeaconState, domain_type: DomainType, epoch: Epoch=None) -> Domain:
//    """
//    Return the signature domain (fork version concatenated with domain type) of a message.
//    """
//    epoch = get_current_epoch(state) if epoch is None else epoch
//    fork_version = state.fork.previous_version if epoch < state.fork.epoch else state.fork.current_version
//    return compute_domain(domain_type, fork_version, state.genesis_validators_root)
func Domain(epoch uint64, domainType []byte, genesisRoot []byte) ([]byte, error) {
	// TODO - add fork version
	return ComputeDomain(domainType, nil, genesisRoot), nil
}

// def compute_domain(domain_type: DomainType, fork_version: Version=None, genesis_validators_root: Root=None) -> Domain:
//    """
//    Return the domain for the ``domain_type`` and ``fork_version``.
//    """
//    if fork_version is None:
//        fork_version = GENESIS_FORK_VERSION
//    if genesis_validators_root is None:
//        genesis_validators_root = Root()  # all bytes zero by default
//    fork_data_root = compute_fork_data_root(fork_version, genesis_validators_root)
//    return Domain(domain_type + fork_data_root[:28])
func ComputeDomain(domainType []byte, forkVersion []byte, genesisValidatorRoot []byte) []byte {
	domainBytes := [4]byte{}
	copy(domainBytes[:], domainType[0:4])

	if forkVersion == nil {
		forkVersion = params.ChainConfig.GenesisForkVersion
	}
	if genesisValidatorRoot == nil {
		genesisValidatorRoot = params.ChainConfig.ZeroHash
	}
	forkBytes := make([]byte, 4)
	copy(forkBytes[:], forkVersion)
	forkDataRoot := make([]byte, 32) // TODO - fork data root

	var b []byte
	b = append(b, domainType[:4]...)
	b = append(b, forkDataRoot[:28]...)
	return b
}