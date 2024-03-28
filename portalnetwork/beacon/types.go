package beacon

import (
	"github.com/ethereum/go-ethereum/portalnetwork/portal_ssz"
	"github.com/ledgerwatch/erigon-lib/types/clonable"
	"github.com/ledgerwatch/erigon/cl/clparams"
	"github.com/ledgerwatch/erigon/cl/cltypes"
	"github.com/ledgerwatch/erigon/cl/cltypes/solid"
	"github.com/ledgerwatch/erigon/cl/merkle_tree"
)

//go:generate sszgen --path types.go --exclude-objs ForkedLightClientUpdate,ForkedLightClientUpdateRange

type LightClientUpdateKey struct {
	StartPeriod uint64
	Count       uint64
}

type LightClientBootstrapKey struct {
	BlockHash []byte `ssz-size:"32"`
}

type LightClientFinalityUpdateKey struct {
	FinalizedSlot uint64
}

type LightClientOptimisticUpdateKey struct {
	OptimisticSlot uint64
}

type ForkedLightClientUpdate struct {
	LightClientUpdate *cltypes.LightClientUpdate
	Version           clparams.StateVersion
}

func NewForkedLightClientUpdate(version clparams.StateVersion) *ForkedLightClientUpdate {
	return &ForkedLightClientUpdate{
		LightClientUpdate: cltypes.NewLightClientUpdate(version),
		Version:           version,
	}
}

func (flgu *ForkedLightClientUpdate) EncodeSSZ(buf []byte) ([]byte, error) {
	dst := buf
	dst = append(dst, portal_ssz.ForkDigest(flgu.Version)...)
	dst, err := flgu.LightClientUpdate.EncodeSSZ(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (flgu *ForkedLightClientUpdate) EncodingSizeSSZ() int {
	return len(portal_ssz.ForkDigest(flgu.Version)) + flgu.LightClientUpdate.EncodingSizeSSZ()
}

func (flgu *ForkedLightClientUpdate) DecodeSSZ(buf []byte, version int) error {
	flgu.Version = clparams.StateVersion(version)
	err := flgu.LightClientUpdate.DecodeSSZ(buf[4:], version)
	if err != nil {
		return err
	}
	return nil
}

func (flgu *ForkedLightClientUpdate) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(portal_ssz.ForkDigest(flgu.Version), flgu.LightClientUpdate)
}

func (flgu *ForkedLightClientUpdate) Clone() clonable.Clonable {
	return NewForkedLightClientUpdate(flgu.Version)
}

type ForkedLightClientUpdateRange *solid.ListSSZ[*ForkedLightClientUpdate]
