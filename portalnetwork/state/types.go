package state

import (
	"errors"
	"fmt"

	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
)

var _ common.SSZObj = (*Nibbles)(nil)

type Nibbles struct {
	Nibbles []byte
}

func (n *Nibbles) Serialize(w *codec.EncodingWriter) error {
	if len(n.Nibbles)%2 == 0 {
		err := w.WriteByte(0)
		if err != nil {
			return err
		}

		for i := 0; i < len(n.Nibbles); i += 2 {
			err = w.WriteByte(n.Nibbles[i]<<4 | n.Nibbles[i+1])
			if err != nil {
				return err
			}
		}
	} else {
		err := w.WriteByte(0x10 | n.Nibbles[0])
		if err != nil {
			return err
		}

		for i := 1; i < len(n.Nibbles); i += 2 {
			err = w.WriteByte(n.Nibbles[i]<<4 | n.Nibbles[i+1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Nibbles) ByteLength() uint64 {
	return uint64(len(n.Nibbles)/2 + 1)
}

func (n *Nibbles) FixedLength() uint64 {
	return 0
}

func (n *Nibbles) Deserialize(dr *codec.DecodingReader) error {
	firstByte, err := dr.ReadByte()
	if err != nil {
		return err
	}

	packedNibbles := make([]byte, dr.Scope())
	_, err = dr.Read(packedNibbles)
	if err != nil {
		return err
	}
	flag, first := unpackNibblePair(firstByte)
	nibbles := make([]byte, 0, 1+2*len(packedNibbles))

	if flag == 0 {
		if first != 0 {
			return fmt.Errorf("nibbles: The lowest 4 bits of the first byte must be 0, but was: %x", first)
		}
	} else if flag == 1 {
		nibbles = append(nibbles, first)
	} else {
		return fmt.Errorf("nibbles: The highest 4 bits of the first byte must be 0 or 1, but was: %x", flag)
	}

	for _, b := range packedNibbles {
		left, right := unpackNibblePair(b)
		nibbles = append(nibbles, left, right)
	}

	unpackedNibbles, err := FromUnpackedNibbles(nibbles)
	if err != nil {
		return err
	}

	*n = *unpackedNibbles
	return nil
}

func (n *Nibbles) HashTreeRoot(h tree.HashFn) tree.Root {
	//TODO implement me
	panic("implement me")
}

func FromUnpackedNibbles(nibbles []byte) (*Nibbles, error) {
	if len(nibbles) > 64 {
		return nil, errors.New("too many nibbles")
	}

	for _, nibble := range nibbles {
		if nibble > 0xf {
			return nil, errors.New("nibble out of range")
		}
	}

	return &Nibbles{Nibbles: nibbles}, nil
}

func unpackNibblePair(pair byte) (byte, byte) {
	return pair >> 4, pair & 0xf
}

// test data from
// https://github.com/ethereum/portal-network-specs/blob/master/state/state-network-test-vectors.md
type AccountTrieNodeKey struct {
	Path     Nibbles
	NodeHash common.Bytes32
}

func (a *AccountTrieNodeKey) Deserialize(dr *codec.DecodingReader) error {
	return dr.Container(
		&a.Path,
		&a.NodeHash,
	)
}

func (a *AccountTrieNodeKey) Serialize(w *codec.EncodingWriter) error {
	return w.Container(
		&a.Path,
		&a.NodeHash,
	)
}

func (a *AccountTrieNodeKey) ByteLength(spec *common.Spec) uint64 {
	return codec.ContainerLength(
		&a.Path,
		&a.NodeHash,
	)
}

func (a *AccountTrieNodeKey) FixedLength(spec *common.Spec) uint64 {
	return 0
}

func (a *AccountTrieNodeKey) HashTreeRoot(spec *common.Spec, hFn tree.HashFn) common.Root {
	return hFn.HashTreeRoot(
		&a.Path,
		&a.NodeHash,
	)
}

type ContractStorageTrieNodeKey struct {
	Address  common.Eth1Address
	Path     Nibbles
	NodeHash common.Bytes32
}

func (c *ContractStorageTrieNodeKey) Deserialize(dr *codec.DecodingReader) error {
	return dr.Container(
		&c.Address,
		&c.Path,
		&c.NodeHash,
	)
}

func (c *ContractStorageTrieNodeKey) Serialize(w *codec.EncodingWriter) error {
	return w.Container(
		&c.Address,
		&c.Path,
		&c.NodeHash,
	)
}

func (c *ContractStorageTrieNodeKey) ByteLength(spec *common.Spec) uint64 {
	return codec.ContainerLength(
		&c.Address,
		&c.Path,
		&c.NodeHash,
	)
}

func (c *ContractStorageTrieNodeKey) FixedLength(spec *common.Spec) uint64 {
	return 0
}

func (c *ContractStorageTrieNodeKey) HashTreeRoot(spec *common.Spec, hFn tree.HashFn) common.Root {
	return hFn.HashTreeRoot(
		&c.Address,
		&c.Path,
		&c.NodeHash,
	)
}

type ContractBytecodeKey struct {
	Address  common.Eth1Address
	NodeHash common.Bytes32
}

func (c *ContractBytecodeKey) Deserialize(dr *codec.DecodingReader) error {
	return dr.FixedLenContainer(
		&c.Address,
		&c.NodeHash,
	)
}

func (c *ContractBytecodeKey) Serialize(w *codec.EncodingWriter) error {
	return w.FixedLenContainer(
		&c.Address,
		&c.NodeHash,
	)
}

func (c *ContractBytecodeKey) ByteLength(spec *common.Spec) uint64 {
	return codec.ContainerLength(
		&c.Address,
		&c.NodeHash,
	)
}

func (c *ContractBytecodeKey) FixedLength(spec *common.Spec) uint64 {
	return 0
}

func (c *ContractBytecodeKey) HashTreeRoot(spec *common.Spec, hFn tree.HashFn) common.Root {
	return hFn.HashTreeRoot(
		&c.Address,
		&c.NodeHash,
	)
}

const MaxTrieNodeLength = 1026
// A content value type, used when retrieving a trie node.
type TrieNode struct {
	Node []byte
} 

func (t *TrieNode) Deserialize(dr *codec.DecodingReader) error {
	return dr.ByteList((*[]byte)(&t.Node), uint64(MaxTrieNodeLength))
}

func (t TrieNode) Serialize(w *codec.EncodingWriter) error {
	return w.Write(t.Node)
}

func (t TrieNode) ByteLength() (out uint64) {
	return uint64(len(t.Node))
}

func (t *TrieNode) FixedLength() uint64 {
	return 0
}

func (t TrieNode) HashTreeRoot(hFn tree.HashFn) tree.Root {
	return hFn.ByteListHTR(t.Node, MaxTrieNodeLength)
}

const MaxTrieProofLength = 65

type TrieProof []TrieNode

func (r *TrieProof) Deserialize(dr *codec.DecodingReader) error {
	return dr.List(func() codec.Deserializable {
		i := len(*r)
		*r = append(*r, TrieNode{})
		return &((*r)[i])
	}, 0, MaxTrieProofLength)
}

func (r TrieProof) Serialize(w *codec.EncodingWriter) error {
	return w.List(func(i uint64) codec.Serializable {
		return &r[i]
	}, 0, uint64(len(r)))
}

func (r TrieProof) ByteLength() (out uint64) {
	for _, v := range r {
		out += v.ByteLength() + codec.OFFSET_SIZE
	}
	return
}

func (r *TrieProof) FixedLength(_ *common.Spec) uint64 {
	return 0
}

func (r TrieProof) HashTreeRoot(hFn tree.HashFn) common.Root {
	length := uint64(len(r))
	return hFn.ComplexListHTR(func(i uint64) tree.HTR {
		if i < length {
			return &r[i]
		}
		return nil
	}, length, MaxTrieProofLength)
}

const MaxContractBytecodeLength = 32768
// A content value type, used when retrieving contract's bytecode.
type ContractBytecode struct {
	Code []byte
} 

func (t *ContractBytecode) Deserialize(dr *codec.DecodingReader) error {
	return dr.ByteList((*[]byte)(&t.Code), uint64(MaxContractBytecodeLength))
}

func (t ContractBytecode) Serialize(w *codec.EncodingWriter) error {
	return w.Write(t.Code)
}

func (t ContractBytecode) ByteLength() (out uint64) {
	return uint64(len(t.Code))
}

func (t *ContractBytecode) FixedLength() uint64 {
	return 0
}

func (t ContractBytecode) HashTreeRoot(hFn tree.HashFn) tree.Root {
	return hFn.ByteListHTR(t.Code, MaxContractBytecodeLength)
}


// A content value type, used when offering a trie node from the account trie.
type AccountTrieNodeWithProof struct{
    /// An proof for the account trie node.
    Proof TrieProof
    /// A block at which the proof is anchored.
    BlockHash common.Bytes32
}

// A content value type, used when offering a trie node from the contract storage trie.
type ContractStorageTrieNodeWithProof struct {
	// A proof for the contract storage trie node.
	StoregeProof TrieProof
	// A proof for the account state.
	AccountProof TrieProof
	// A block at which the proof is anchored.
	BlockHash common.Bytes32
}

// A content value type, used when offering contract's bytecode.
type ContractBytecodeWithProof struct {
	// A contract's bytecode.
	Code ContractBytecode
	// A proof for the account state of the corresponding contract.
	AccountProof TrieProof
	// A block at which the proof is anchored.
	BlockHash common.Bytes32
}