package beacon

import (
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/capella"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
)

var (
	Bellatrix common.ForkDigest = [4]byte{0x0, 0x0, 0x0, 0x0}
	Capella   common.ForkDigest = [4]byte{0xbb, 0xa4, 0xda, 0x96}
)

// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/light-client/sync-protocol.md
const ExecutionBranchLength = 4

var ExecutionBranchType = view.VectorType(common.Bytes32Type, ExecutionBranchLength)

type ExecutionBranch [ExecutionBranchLength]common.Bytes32

func (eb *ExecutionBranch) Deserialize(dr *codec.DecodingReader) error {
	roots := eb[:]
	return tree.ReadRoots(dr, &roots, 4)
}

func (eb *ExecutionBranch) FixedLength() uint64 {
	return ExecutionBranchType.TypeByteLength()
}

func (eb *ExecutionBranch) Serialize(w *codec.EncodingWriter) error {
	return tree.WriteRoots(w, eb[:])
}

func (eb *ExecutionBranch) ByteLength() (out uint64) {
	return ExecutionBranchType.TypeByteLength()
}

func (eb *ExecutionBranch) HashTreeRoot(hFn tree.HashFn) common.Root {
	return hFn.ComplexVectorHTR(func(i uint64) tree.HTR {
		if i < ExecutionBranchLength {
			return &eb[i]
		}
		return nil
	}, ExecutionBranchLength)
}

type LightClientHeader struct {
	Beacon          common.BeaconBlockHeader
	Execution       capella.ExecutionPayloadHeader
	ExecutionBranch ExecutionBranch
}

var LightClientHeaderType = view.ContainerType("LightClientHeader", []view.FieldDef{
	{Name: "beacon", Type: common.BeaconBlockHeaderType},
	{Name: "execution", Type: capella.ExecutionPayloadHeaderType},
	{Name: "execution_branch", Type: ExecutionBranchType},
})

func (l *LightClientHeader) Deserialize(spec *SpecWithForkDigest, dr *codec.DecodingReader) error {
	if spec.ForkDigest == Capella {
		return dr.Container(&l.Beacon, &l.Execution, &l.ExecutionBranch)
	}
	return dr.FixedLenContainer(&l.Beacon)
}

func (l *LightClientHeader) FixedLength(spec *SpecWithForkDigest) uint64 {
	if spec.ForkDigest == Capella {
		return LightClientHeaderType.TypeByteLength()
	}
	return codec.ContainerLength(&l.Beacon)
}

func (l *LightClientHeader) Serialize(spec *SpecWithForkDigest, w *codec.EncodingWriter) error {
	if spec.ForkDigest == Capella {
		return w.Container(&l.Beacon, &l.Execution, &l.ExecutionBranch)
	}
	return w.FixedLenContainer(&l.Beacon)
}

func (l *LightClientHeader) ByteLength(spec *SpecWithForkDigest) (out uint64) {
	if spec.ForkDigest == Capella {
		return LightClientHeaderType.TypeByteLength()
	}
	return codec.ContainerLength(&l.Beacon)
}

func (l *LightClientHeader) HashTreeRoot(spec *SpecWithForkDigest, h tree.HashFn) common.Root {
	if spec.ForkDigest == Capella {
		return h.HashTreeRoot(&l.Beacon, &l.Execution, &l.ExecutionBranch)
	}
	return h.HashTreeRoot(&l.Beacon)
}

type LightClientBootstrap struct {
	Header                     LightClientHeader
	CurrentSyncCommittee       common.SyncCommittee
	CurrentSyncCommitteeBranch altair.SyncCommitteeProofBranch
}

func NewLightClientBootstrapType(spec *common.Spec) *view.ContainerTypeDef {
	return view.ContainerType("LightClientHeader", []view.FieldDef{
		{Name: "header", Type: LightClientHeaderType},
		{Name: "next_sync_committee", Type: common.SyncCommitteeType(spec)},
		{Name: "next_sync_committee_branch", Type: altair.SyncCommitteeProofBranchType},
	})
}

func (lcb *LightClientBootstrap) FixedLength(spec *SpecWithForkDigest) uint64 {
	return codec.ContainerLength(spec.Wrap(&lcb.Header), spec.Spec.Wrap(&lcb.CurrentSyncCommittee), &lcb.CurrentSyncCommitteeBranch)
}

func (lcb *LightClientBootstrap) Deserialize(spec *SpecWithForkDigest, dr *codec.DecodingReader) error {
	return dr.Container(spec.Wrap(&lcb.Header), spec.Spec.Wrap(&lcb.CurrentSyncCommittee), &lcb.CurrentSyncCommitteeBranch)
}

func (lcb *LightClientBootstrap) Serialize(spec *SpecWithForkDigest, w *codec.EncodingWriter) error {
	return w.Container(spec.Wrap(&lcb.Header), spec.Spec.Wrap(&lcb.CurrentSyncCommittee), &lcb.CurrentSyncCommitteeBranch)
}

func (lcb *LightClientBootstrap) ByteLength(spec *SpecWithForkDigest) uint64 {
	return codec.ContainerLength(spec.Wrap(&lcb.Header), spec.Spec.Wrap(&lcb.CurrentSyncCommittee), &lcb.CurrentSyncCommitteeBranch)
}
