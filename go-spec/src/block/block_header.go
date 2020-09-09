package block

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type BlockHeader struct {
	BlockRoot 			[]byte
	Signature			[]byte // TODO - checking validity + how many voted?
}

func (header *BlockHeader) Copy() core.IBlockHeader {
	return &BlockHeader{
		BlockRoot: header.BlockRoot,
		Signature: header.Signature,
	}
}

func (header *BlockHeader) Validate(bp core.IBlockProducer) error {
	sig := &bls.Sign{}
	err := sig.Deserialize(header.Signature)
	if err != nil {
		return err
	}

	pk, err := bp.GetPubKey()
	if err != nil {
		return err
	}

	if res := sig.VerifyHash(pk, header.BlockRoot); !res {
		return fmt.Errorf("signature did not verify")
	}
	return nil
}

func (header *BlockHeader) GetBlockRoot() []byte {
	return header.BlockRoot
}

func (header *BlockHeader) GetSignature() []byte {
	return header.Signature
}

func NewBlockHeader(sk *bls.SecretKey, body *BlockBody) (*BlockHeader,error) {
	root,err := body.Root()
	if err != nil {
		return nil, err
	}

	sig := sk.SignByte(root)

	return &BlockHeader{
		BlockRoot: root,
		Signature: sig.Serialize(),
	}, nil
}