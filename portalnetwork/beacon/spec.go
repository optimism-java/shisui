package beacon

import (
	"encoding/json"

	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	"gopkg.in/yaml.v3"
)

type SpecWithForkDigest struct {
	Spec       *common.Spec
	ForkDigest common.ForkDigest
}

type SpecObj interface {
	Deserialize(spec *SpecWithForkDigest, dr *codec.DecodingReader) error
	Serialize(spec *SpecWithForkDigest, w *codec.EncodingWriter) error
	ByteLength(spec *SpecWithForkDigest) uint64
	HashTreeRoot(spec *SpecWithForkDigest, h tree.HashFn) common.Root
	FixedLength(spec *SpecWithForkDigest) uint64
}

type SSZObj interface {
	codec.Serializable
	codec.Deserializable
	codec.FixedLength
	tree.HTR
}

type WrappedSpecObj interface {
	SSZObj
	Unwrap() (*SpecWithForkDigest, SpecObj)
}

type specObj struct {
	spec *SpecWithForkDigest
	des  SpecObj
}

func (s *specObj) Deserialize(dr *codec.DecodingReader) error {
	return s.des.Deserialize(s.spec, dr)
}

func (s *specObj) Serialize(w *codec.EncodingWriter) error {
	return s.des.Serialize(s.spec, w)
}

func (s *specObj) ByteLength() uint64 {
	return s.des.ByteLength(s.spec)
}

func (s *specObj) HashTreeRoot(h tree.HashFn) common.Root {
	return s.des.HashTreeRoot(s.spec, h)
}

func (s *specObj) FixedLength() uint64 {
	return s.des.FixedLength(s.spec)
}

func (s *specObj) Unwrap() (*SpecWithForkDigest, SpecObj) {
	return s.spec, s.des
}

func (s *specObj) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, s.des)
}

func (s *specObj) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.des)
}

func (s *specObj) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(s.des)
}

func (s *specObj) MarshalYAML() (interface{}, error) {
	return s.des, nil
}

// Wraps the object to parametrize with given spec. JSON and YAML functionality is proxied to the inner value.
func (spec *SpecWithForkDigest) Wrap(des SpecObj) SSZObj {
	return &specObj{spec, des}
}

func (spec *SpecWithForkDigest) ForkVersion(slot common.Slot) common.Version {
	epoch := spec.Spec.SlotToEpoch(slot)
	if epoch < spec.Spec.ALTAIR_FORK_EPOCH {
		return spec.Spec.GENESIS_FORK_VERSION
	} else if epoch < spec.Spec.BELLATRIX_FORK_EPOCH {
		return spec.Spec.ALTAIR_FORK_VERSION
	} else if epoch < spec.Spec.CAPELLA_FORK_EPOCH {
		return spec.Spec.BELLATRIX_FORK_VERSION
	} else {
		return spec.Spec.CAPELLA_FORK_VERSION
	}
}
