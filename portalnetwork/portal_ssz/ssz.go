package portal_ssz

import (
	"github.com/ledgerwatch/erigon-lib/types/ssz"
	"github.com/ledgerwatch/erigon/cl/clparams"
)

var (
	Bellatrix = []byte{0x0, 0x0, 0x0, 0x0}
	Capella   = []byte{0xbb, 0xa4, 0xda, 0x96}
)

func ForkVersion(digest []byte) clparams.StateVersion {
	if digest[0] == 0xbb && digest[1] == 0xa4 && digest[2] == 0xda && digest[3] == 0x96 {
		return clparams.CapellaVersion
	}
	return clparams.BellatrixVersion
}

func ForkDigest(version clparams.StateVersion) []byte {
	if version < clparams.CapellaVersion {
		return Bellatrix
	}
	return Capella
}

func DecodeDynamicListForkedObject[T ssz.Unmarshaler](bytes []byte, start, end uint32, max uint64) ([]T, error) {
	if start > end || len(bytes) < int(end) {
		return nil, ssz.ErrBadOffset
	}
	buf := bytes[start:end]
	var elementsNum, currentOffset uint32
	if len(buf) > 4 {
		currentOffset = ssz.DecodeOffset(buf)
		elementsNum = currentOffset / 4
	}
	inPos := 4
	if uint64(elementsNum) > max {
		return nil, ssz.ErrTooBigList
	}
	objs := make([]T, elementsNum)
	for i := range objs {
		endOffset := uint32(len(buf))
		if i != len(objs)-1 {
			if len(buf[inPos:]) < 4 {
				return nil, ssz.ErrLowBufferSize
			}
			endOffset = ssz.DecodeOffset(buf[inPos:])
		}
		inPos += 4
		if endOffset < currentOffset || len(buf) < int(endOffset) {
			return nil, ssz.ErrBadOffset
		}
		objs[i] = objs[i].Clone().(T)
		if err := objs[i].DecodeSSZ(buf[currentOffset:endOffset], int(ForkVersion(buf[currentOffset:currentOffset+4]))); err != nil {
			return nil, err
		}
		currentOffset = endOffset
	}
	return objs, nil
}
