// Code generated by fastssz. DO NOT EDIT.
// Hash: 02cbfc3e581dd4391beb495f34149023b8d4e13bf82aba3bd385564bfa1ec2f2
// Version: 0.1.3
package history

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the HeaderRecord object
func (h *HeaderRecord) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(h)
}

// MarshalSSZTo ssz marshals the HeaderRecord object to a target array
func (h *HeaderRecord) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'BlockHash'
	if size := len(h.BlockHash); size != 32 {
		err = ssz.ErrBytesLengthFn("HeaderRecord.BlockHash", size, 32)
		return
	}
	dst = append(dst, h.BlockHash...)

	// Field (1) 'TotalDifficulty'
	if size := len(h.TotalDifficulty); size != 32 {
		err = ssz.ErrBytesLengthFn("HeaderRecord.TotalDifficulty", size, 32)
		return
	}
	dst = append(dst, h.TotalDifficulty...)

	return
}

// UnmarshalSSZ ssz unmarshals the HeaderRecord object
func (h *HeaderRecord) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 64 {
		return ssz.ErrSize
	}

	// Field (0) 'BlockHash'
	if cap(h.BlockHash) == 0 {
		h.BlockHash = make([]byte, 0, len(buf[0:32]))
	}
	h.BlockHash = append(h.BlockHash, buf[0:32]...)

	// Field (1) 'TotalDifficulty'
	if cap(h.TotalDifficulty) == 0 {
		h.TotalDifficulty = make([]byte, 0, len(buf[32:64]))
	}
	h.TotalDifficulty = append(h.TotalDifficulty, buf[32:64]...)

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the HeaderRecord object
func (h *HeaderRecord) SizeSSZ() (size int) {
	size = 64
	return
}

// HashTreeRoot ssz hashes the HeaderRecord object
func (h *HeaderRecord) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(h)
}

// HashTreeRootWith ssz hashes the HeaderRecord object with a hasher
func (h *HeaderRecord) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'BlockHash'
	if size := len(h.BlockHash); size != 32 {
		err = ssz.ErrBytesLengthFn("HeaderRecord.BlockHash", size, 32)
		return
	}
	hh.PutBytes(h.BlockHash)

	// Field (1) 'TotalDifficulty'
	if size := len(h.TotalDifficulty); size != 32 {
		err = ssz.ErrBytesLengthFn("HeaderRecord.TotalDifficulty", size, 32)
		return
	}
	hh.PutBytes(h.TotalDifficulty)

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the HeaderRecord object
func (h *HeaderRecord) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(h)
}

// MarshalSSZ ssz marshals the EpochAccumulator object
func (e *EpochAccumulator) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(e)
}

// MarshalSSZTo ssz marshals the EpochAccumulator object to a target array
func (e *EpochAccumulator) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'HeaderRecords'
	if size := len(e.HeaderRecords); size != 8192 {
		err = ssz.ErrVectorLengthFn("EpochAccumulator.HeaderRecords", size, 8192)
		return
	}
	for ii := 0; ii < 8192; ii++ {
		if size := len(e.HeaderRecords[ii]); size != 64 {
			err = ssz.ErrBytesLengthFn("EpochAccumulator.HeaderRecords[ii]", size, 64)
			return
		}
		dst = append(dst, e.HeaderRecords[ii]...)
	}

	return
}

// UnmarshalSSZ ssz unmarshals the EpochAccumulator object
func (e *EpochAccumulator) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 524288 {
		return ssz.ErrSize
	}

	// Field (0) 'HeaderRecords'
	e.HeaderRecords = make([][]byte, 8192)
	for ii := 0; ii < 8192; ii++ {
		if cap(e.HeaderRecords[ii]) == 0 {
			e.HeaderRecords[ii] = make([]byte, 0, len(buf[0:524288][ii*64:(ii+1)*64]))
		}
		e.HeaderRecords[ii] = append(e.HeaderRecords[ii], buf[0:524288][ii*64:(ii+1)*64]...)
	}

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the EpochAccumulator object
func (e *EpochAccumulator) SizeSSZ() (size int) {
	size = 524288
	return
}

// HashTreeRoot ssz hashes the EpochAccumulator object
func (e *EpochAccumulator) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(e)
}

// HashTreeRootWith ssz hashes the EpochAccumulator object with a hasher
func (e *EpochAccumulator) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'HeaderRecords'
	{
		if size := len(e.HeaderRecords); size != 8192 {
			err = ssz.ErrVectorLengthFn("EpochAccumulator.HeaderRecords", size, 8192)
			return
		}
		subIndx := hh.Index()
		for _, i := range e.HeaderRecords {
			if len(i) != 64 {
				err = ssz.ErrBytesLength
				return
			}
			hh.PutBytes(i)
		}
		hh.Merkleize(subIndx)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the EpochAccumulator object
func (e *EpochAccumulator) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(e)
}

// MarshalSSZ ssz marshals the BlockBodyLegacy object
func (b *BlockBodyLegacy) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(b)
}

// MarshalSSZTo ssz marshals the BlockBodyLegacy object to a target array
func (b *BlockBodyLegacy) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(8)

	// Offset (0) 'Transactions'
	dst = ssz.WriteOffset(dst, offset)
	for ii := 0; ii < len(b.Transactions); ii++ {
		offset += 4
		offset += len(b.Transactions[ii])
	}

	// Offset (1) 'Uncles'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(b.Uncles)

	// Field (0) 'Transactions'
	if size := len(b.Transactions); size > 16384 {
		err = ssz.ErrListTooBigFn("BlockBodyLegacy.Transactions", size, 16384)
		return
	}
	{
		offset = 4 * len(b.Transactions)
		for ii := 0; ii < len(b.Transactions); ii++ {
			dst = ssz.WriteOffset(dst, offset)
			offset += len(b.Transactions[ii])
		}
	}
	for ii := 0; ii < len(b.Transactions); ii++ {
		if size := len(b.Transactions[ii]); size > 16777216 {
			err = ssz.ErrBytesLengthFn("BlockBodyLegacy.Transactions[ii]", size, 16777216)
			return
		}
		dst = append(dst, b.Transactions[ii]...)
	}

	// Field (1) 'Uncles'
	if size := len(b.Uncles); size > 131072 {
		err = ssz.ErrBytesLengthFn("BlockBodyLegacy.Uncles", size, 131072)
		return
	}
	dst = append(dst, b.Uncles...)

	return
}

// UnmarshalSSZ ssz unmarshals the BlockBodyLegacy object
func (b *BlockBodyLegacy) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 8 {
		return ssz.ErrSize
	}

	tail := buf
	var o0, o1 uint64

	// Offset (0) 'Transactions'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 8 {
		return ssz.ErrInvalidVariableOffset
	}

	// Offset (1) 'Uncles'
	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
		return ssz.ErrOffset
	}

	// Field (0) 'Transactions'
	{
		buf = tail[o0:o1]
		num, err := ssz.DecodeDynamicLength(buf, 16384)
		if err != nil {
			return err
		}
		b.Transactions = make([][]byte, num)
		err = ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
			if len(buf) > 16777216 {
				return ssz.ErrBytesLength
			}
			if cap(b.Transactions[indx]) == 0 {
				b.Transactions[indx] = make([]byte, 0, len(buf))
			}
			b.Transactions[indx] = append(b.Transactions[indx], buf...)
			return nil
		})
		if err != nil {
			return err
		}
	}

	// Field (1) 'Uncles'
	{
		buf = tail[o1:]
		if len(buf) > 131072 {
			return ssz.ErrBytesLength
		}
		if cap(b.Uncles) == 0 {
			b.Uncles = make([]byte, 0, len(buf))
		}
		b.Uncles = append(b.Uncles, buf...)
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the BlockBodyLegacy object
func (b *BlockBodyLegacy) SizeSSZ() (size int) {
	size = 8

	// Field (0) 'Transactions'
	for ii := 0; ii < len(b.Transactions); ii++ {
		size += 4
		size += len(b.Transactions[ii])
	}

	// Field (1) 'Uncles'
	size += len(b.Uncles)

	return
}

// HashTreeRoot ssz hashes the BlockBodyLegacy object
func (b *BlockBodyLegacy) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

// HashTreeRootWith ssz hashes the BlockBodyLegacy object with a hasher
func (b *BlockBodyLegacy) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Transactions'
	{
		subIndx := hh.Index()
		num := uint64(len(b.Transactions))
		if num > 16384 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range b.Transactions {
			{
				elemIndx := hh.Index()
				byteLen := uint64(len(elem))
				if byteLen > 16777216 {
					err = ssz.ErrIncorrectListSize
					return
				}
				hh.AppendBytes32(elem)
				hh.MerkleizeWithMixin(elemIndx, byteLen, (16777216+31)/32)
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16384)
	}

	// Field (1) 'Uncles'
	{
		elemIndx := hh.Index()
		byteLen := uint64(len(b.Uncles))
		if byteLen > 131072 {
			err = ssz.ErrIncorrectListSize
			return
		}
		hh.Append(b.Uncles)
		hh.MerkleizeWithMixin(elemIndx, byteLen, (131072+31)/32)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the BlockBodyLegacy object
func (b *BlockBodyLegacy) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(b)
}

// MarshalSSZ ssz marshals the PortalBlockBodyShanghai object
func (p *PortalBlockBodyShanghai) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(p)
}

// MarshalSSZTo ssz marshals the PortalBlockBodyShanghai object to a target array
func (p *PortalBlockBodyShanghai) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(12)

	// Offset (0) 'Transactions'
	dst = ssz.WriteOffset(dst, offset)
	for ii := 0; ii < len(p.Transactions); ii++ {
		offset += 4
		offset += len(p.Transactions[ii])
	}

	// Offset (1) 'Uncles'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(p.Uncles)

	// Offset (2) 'Withdrawals'
	dst = ssz.WriteOffset(dst, offset)
	for ii := 0; ii < len(p.Withdrawals); ii++ {
		offset += 4
		offset += len(p.Withdrawals[ii])
	}

	// Field (0) 'Transactions'
	if size := len(p.Transactions); size > 16384 {
		err = ssz.ErrListTooBigFn("PortalBlockBodyShanghai.Transactions", size, 16384)
		return
	}
	{
		offset = 4 * len(p.Transactions)
		for ii := 0; ii < len(p.Transactions); ii++ {
			dst = ssz.WriteOffset(dst, offset)
			offset += len(p.Transactions[ii])
		}
	}
	for ii := 0; ii < len(p.Transactions); ii++ {
		if size := len(p.Transactions[ii]); size > 16777216 {
			err = ssz.ErrBytesLengthFn("PortalBlockBodyShanghai.Transactions[ii]", size, 16777216)
			return
		}
		dst = append(dst, p.Transactions[ii]...)
	}

	// Field (1) 'Uncles'
	if size := len(p.Uncles); size > 131072 {
		err = ssz.ErrBytesLengthFn("PortalBlockBodyShanghai.Uncles", size, 131072)
		return
	}
	dst = append(dst, p.Uncles...)

	// Field (2) 'Withdrawals'
	if size := len(p.Withdrawals); size > 16 {
		err = ssz.ErrListTooBigFn("PortalBlockBodyShanghai.Withdrawals", size, 16)
		return
	}
	{
		offset = 4 * len(p.Withdrawals)
		for ii := 0; ii < len(p.Withdrawals); ii++ {
			dst = ssz.WriteOffset(dst, offset)
			offset += len(p.Withdrawals[ii])
		}
	}
	for ii := 0; ii < len(p.Withdrawals); ii++ {
		if size := len(p.Withdrawals[ii]); size > 192 {
			err = ssz.ErrBytesLengthFn("PortalBlockBodyShanghai.Withdrawals[ii]", size, 192)
			return
		}
		dst = append(dst, p.Withdrawals[ii]...)
	}

	return
}

// UnmarshalSSZ ssz unmarshals the PortalBlockBodyShanghai object
func (p *PortalBlockBodyShanghai) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 12 {
		return ssz.ErrSize
	}

	tail := buf
	var o0, o1, o2 uint64

	// Offset (0) 'Transactions'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 12 {
		return ssz.ErrInvalidVariableOffset
	}

	// Offset (1) 'Uncles'
	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
		return ssz.ErrOffset
	}

	// Offset (2) 'Withdrawals'
	if o2 = ssz.ReadOffset(buf[8:12]); o2 > size || o1 > o2 {
		return ssz.ErrOffset
	}

	// Field (0) 'Transactions'
	{
		buf = tail[o0:o1]
		num, err := ssz.DecodeDynamicLength(buf, 16384)
		if err != nil {
			return err
		}
		p.Transactions = make([][]byte, num)
		err = ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
			if len(buf) > 16777216 {
				return ssz.ErrBytesLength
			}
			if cap(p.Transactions[indx]) == 0 {
				p.Transactions[indx] = make([]byte, 0, len(buf))
			}
			p.Transactions[indx] = append(p.Transactions[indx], buf...)
			return nil
		})
		if err != nil {
			return err
		}
	}

	// Field (1) 'Uncles'
	{
		buf = tail[o1:o2]
		if len(buf) > 131072 {
			return ssz.ErrBytesLength
		}
		if cap(p.Uncles) == 0 {
			p.Uncles = make([]byte, 0, len(buf))
		}
		p.Uncles = append(p.Uncles, buf...)
	}

	// Field (2) 'Withdrawals'
	{
		buf = tail[o2:]
		num, err := ssz.DecodeDynamicLength(buf, 16)
		if err != nil {
			return err
		}
		p.Withdrawals = make([][]byte, num)
		err = ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
			if len(buf) > 192 {
				return ssz.ErrBytesLength
			}
			if cap(p.Withdrawals[indx]) == 0 {
				p.Withdrawals[indx] = make([]byte, 0, len(buf))
			}
			p.Withdrawals[indx] = append(p.Withdrawals[indx], buf...)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the PortalBlockBodyShanghai object
func (p *PortalBlockBodyShanghai) SizeSSZ() (size int) {
	size = 12

	// Field (0) 'Transactions'
	for ii := 0; ii < len(p.Transactions); ii++ {
		size += 4
		size += len(p.Transactions[ii])
	}

	// Field (1) 'Uncles'
	size += len(p.Uncles)

	// Field (2) 'Withdrawals'
	for ii := 0; ii < len(p.Withdrawals); ii++ {
		size += 4
		size += len(p.Withdrawals[ii])
	}

	return
}

// HashTreeRoot ssz hashes the PortalBlockBodyShanghai object
func (p *PortalBlockBodyShanghai) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(p)
}

// HashTreeRootWith ssz hashes the PortalBlockBodyShanghai object with a hasher
func (p *PortalBlockBodyShanghai) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Transactions'
	{
		subIndx := hh.Index()
		num := uint64(len(p.Transactions))
		if num > 16384 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range p.Transactions {
			{
				elemIndx := hh.Index()
				byteLen := uint64(len(elem))
				if byteLen > 16777216 {
					err = ssz.ErrIncorrectListSize
					return
				}
				hh.AppendBytes32(elem)
				hh.MerkleizeWithMixin(elemIndx, byteLen, (16777216+31)/32)
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16384)
	}

	// Field (1) 'Uncles'
	{
		elemIndx := hh.Index()
		byteLen := uint64(len(p.Uncles))
		if byteLen > 131072 {
			err = ssz.ErrIncorrectListSize
			return
		}
		hh.Append(p.Uncles)
		hh.MerkleizeWithMixin(elemIndx, byteLen, (131072+31)/32)
	}

	// Field (2) 'Withdrawals'
	{
		subIndx := hh.Index()
		num := uint64(len(p.Withdrawals))
		if num > 16 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range p.Withdrawals {
			{
				elemIndx := hh.Index()
				byteLen := uint64(len(elem))
				if byteLen > 192 {
					err = ssz.ErrIncorrectListSize
					return
				}
				hh.AppendBytes32(elem)
				hh.MerkleizeWithMixin(elemIndx, byteLen, (192+31)/32)
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the PortalBlockBodyShanghai object
func (p *PortalBlockBodyShanghai) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(p)
}

// MarshalSSZ ssz marshals the BlockHeaderWithProof object
func (b *BlockHeaderWithProof) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(b)
}

// MarshalSSZTo ssz marshals the BlockHeaderWithProof object to a target array
func (b *BlockHeaderWithProof) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(8)

	// Offset (0) 'Header'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(b.Header)

	// Offset (1) 'Proof'
	dst = ssz.WriteOffset(dst, offset)
	if b.Proof == nil {
		b.Proof = new(BlockHeaderProof)
	}
	offset += b.Proof.SizeSSZ()

	// Field (0) 'Header'
	if size := len(b.Header); size > 8192 {
		err = ssz.ErrBytesLengthFn("BlockHeaderWithProof.Header", size, 8192)
		return
	}
	dst = append(dst, b.Header...)

	// Field (1) 'Proof'
	if dst, err = b.Proof.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the BlockHeaderWithProof object
func (b *BlockHeaderWithProof) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 8 {
		return ssz.ErrSize
	}

	tail := buf
	var o0, o1 uint64

	// Offset (0) 'Header'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 8 {
		return ssz.ErrInvalidVariableOffset
	}

	// Offset (1) 'Proof'
	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
		return ssz.ErrOffset
	}

	// Field (0) 'Header'
	{
		buf = tail[o0:o1]
		if len(buf) > 8192 {
			return ssz.ErrBytesLength
		}
		if cap(b.Header) == 0 {
			b.Header = make([]byte, 0, len(buf))
		}
		b.Header = append(b.Header, buf...)
	}

	// Field (1) 'Proof'
	{
		buf = tail[o1:]
		if b.Proof == nil {
			b.Proof = new(BlockHeaderProof)
		}
		if err = b.Proof.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the BlockHeaderWithProof object
func (b *BlockHeaderWithProof) SizeSSZ() (size int) {
	size = 8

	// Field (0) 'Header'
	size += len(b.Header)

	// Field (1) 'Proof'
	if b.Proof == nil {
		b.Proof = new(BlockHeaderProof)
	}
	size += b.Proof.SizeSSZ()

	return
}

// HashTreeRoot ssz hashes the BlockHeaderWithProof object
func (b *BlockHeaderWithProof) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

// HashTreeRootWith ssz hashes the BlockHeaderWithProof object with a hasher
func (b *BlockHeaderWithProof) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Header'
	{
		elemIndx := hh.Index()
		byteLen := uint64(len(b.Header))
		if byteLen > 8192 {
			err = ssz.ErrIncorrectListSize
			return
		}
		hh.Append(b.Header)
		hh.MerkleizeWithMixin(elemIndx, byteLen, (8192+31)/32)
	}

	// Field (1) 'Proof'
	if err = b.Proof.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the BlockHeaderWithProof object
func (b *BlockHeaderWithProof) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(b)
}

// MarshalSSZ ssz marshals the SSZProof object
func (s *SSZProof) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(s)
}

// MarshalSSZTo ssz marshals the SSZProof object to a target array
func (s *SSZProof) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(36)

	// Field (0) 'Leaf'
	if size := len(s.Leaf); size != 32 {
		err = ssz.ErrBytesLengthFn("SSZProof.Leaf", size, 32)
		return
	}
	dst = append(dst, s.Leaf...)

	// Offset (1) 'Witnesses'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(s.Witnesses) * 32

	// Field (1) 'Witnesses'
	if size := len(s.Witnesses); size > 65536 {
		err = ssz.ErrListTooBigFn("SSZProof.Witnesses", size, 65536)
		return
	}
	for ii := 0; ii < len(s.Witnesses); ii++ {
		if size := len(s.Witnesses[ii]); size != 32 {
			err = ssz.ErrBytesLengthFn("SSZProof.Witnesses[ii]", size, 32)
			return
		}
		dst = append(dst, s.Witnesses[ii]...)
	}

	return
}

// UnmarshalSSZ ssz unmarshals the SSZProof object
func (s *SSZProof) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 36 {
		return ssz.ErrSize
	}

	tail := buf
	var o1 uint64

	// Field (0) 'Leaf'
	if cap(s.Leaf) == 0 {
		s.Leaf = make([]byte, 0, len(buf[0:32]))
	}
	s.Leaf = append(s.Leaf, buf[0:32]...)

	// Offset (1) 'Witnesses'
	if o1 = ssz.ReadOffset(buf[32:36]); o1 > size {
		return ssz.ErrOffset
	}

	if o1 < 36 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (1) 'Witnesses'
	{
		buf = tail[o1:]
		num, err := ssz.DivideInt2(len(buf), 32, 65536)
		if err != nil {
			return err
		}
		s.Witnesses = make([][]byte, num)
		for ii := 0; ii < num; ii++ {
			if cap(s.Witnesses[ii]) == 0 {
				s.Witnesses[ii] = make([]byte, 0, len(buf[ii*32:(ii+1)*32]))
			}
			s.Witnesses[ii] = append(s.Witnesses[ii], buf[ii*32:(ii+1)*32]...)
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the SSZProof object
func (s *SSZProof) SizeSSZ() (size int) {
	size = 36

	// Field (1) 'Witnesses'
	size += len(s.Witnesses) * 32

	return
}

// HashTreeRoot ssz hashes the SSZProof object
func (s *SSZProof) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(s)
}

// HashTreeRootWith ssz hashes the SSZProof object with a hasher
func (s *SSZProof) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Leaf'
	if size := len(s.Leaf); size != 32 {
		err = ssz.ErrBytesLengthFn("SSZProof.Leaf", size, 32)
		return
	}
	hh.PutBytes(s.Leaf)

	// Field (1) 'Witnesses'
	{
		if size := len(s.Witnesses); size > 65536 {
			err = ssz.ErrListTooBigFn("SSZProof.Witnesses", size, 65536)
			return
		}
		subIndx := hh.Index()
		for _, i := range s.Witnesses {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				return
			}
			hh.Append(i)
		}
		numItems := uint64(len(s.Witnesses))
		hh.MerkleizeWithMixin(subIndx, numItems, 65536)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the SSZProof object
func (s *SSZProof) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(s)
}

// MarshalSSZ ssz marshals the MasterAccumulator object
func (m *MasterAccumulator) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(m)
}

// MarshalSSZTo ssz marshals the MasterAccumulator object to a target array
func (m *MasterAccumulator) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(4)

	// Offset (0) 'HistoricalEpochs'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(m.HistoricalEpochs) * 32

	// Field (0) 'HistoricalEpochs'
	if size := len(m.HistoricalEpochs); size > 1897 {
		err = ssz.ErrListTooBigFn("MasterAccumulator.HistoricalEpochs", size, 1897)
		return
	}
	for ii := 0; ii < len(m.HistoricalEpochs); ii++ {
		if size := len(m.HistoricalEpochs[ii]); size != 32 {
			err = ssz.ErrBytesLengthFn("MasterAccumulator.HistoricalEpochs[ii]", size, 32)
			return
		}
		dst = append(dst, m.HistoricalEpochs[ii]...)
	}

	return
}

// UnmarshalSSZ ssz unmarshals the MasterAccumulator object
func (m *MasterAccumulator) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 4 {
		return ssz.ErrSize
	}

	tail := buf
	var o0 uint64

	// Offset (0) 'HistoricalEpochs'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 4 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (0) 'HistoricalEpochs'
	{
		buf = tail[o0:]
		num, err := ssz.DivideInt2(len(buf), 32, 1897)
		if err != nil {
			return err
		}
		m.HistoricalEpochs = make([][]byte, num)
		for ii := 0; ii < num; ii++ {
			if cap(m.HistoricalEpochs[ii]) == 0 {
				m.HistoricalEpochs[ii] = make([]byte, 0, len(buf[ii*32:(ii+1)*32]))
			}
			m.HistoricalEpochs[ii] = append(m.HistoricalEpochs[ii], buf[ii*32:(ii+1)*32]...)
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the MasterAccumulator object
func (m *MasterAccumulator) SizeSSZ() (size int) {
	size = 4

	// Field (0) 'HistoricalEpochs'
	size += len(m.HistoricalEpochs) * 32

	return
}

// HashTreeRoot ssz hashes the MasterAccumulator object
func (m *MasterAccumulator) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(m)
}

// HashTreeRootWith ssz hashes the MasterAccumulator object with a hasher
func (m *MasterAccumulator) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'HistoricalEpochs'
	{
		if size := len(m.HistoricalEpochs); size > 1897 {
			err = ssz.ErrListTooBigFn("MasterAccumulator.HistoricalEpochs", size, 1897)
			return
		}
		subIndx := hh.Index()
		for _, i := range m.HistoricalEpochs {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				return
			}
			hh.Append(i)
		}
		numItems := uint64(len(m.HistoricalEpochs))
		hh.MerkleizeWithMixin(subIndx, numItems, 1897)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the MasterAccumulator object
func (m *MasterAccumulator) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(m)
}

// MarshalSSZ ssz marshals the PortalReceipts object
func (p *PortalReceipts) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(p)
}

// MarshalSSZTo ssz marshals the PortalReceipts object to a target array
func (p *PortalReceipts) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(4)

	// Offset (0) 'Receipts'
	dst = ssz.WriteOffset(dst, offset)
	for ii := 0; ii < len(p.Receipts); ii++ {
		offset += 4
		offset += len(p.Receipts[ii])
	}

	// Field (0) 'Receipts'
	if size := len(p.Receipts); size > 16384 {
		err = ssz.ErrListTooBigFn("PortalReceipts.Receipts", size, 16384)
		return
	}
	{
		offset = 4 * len(p.Receipts)
		for ii := 0; ii < len(p.Receipts); ii++ {
			dst = ssz.WriteOffset(dst, offset)
			offset += len(p.Receipts[ii])
		}
	}
	for ii := 0; ii < len(p.Receipts); ii++ {
		if size := len(p.Receipts[ii]); size > 134217728 {
			err = ssz.ErrBytesLengthFn("PortalReceipts.Receipts[ii]", size, 134217728)
			return
		}
		dst = append(dst, p.Receipts[ii]...)
	}

	return
}

// UnmarshalSSZ ssz unmarshals the PortalReceipts object
func (p *PortalReceipts) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 4 {
		return ssz.ErrSize
	}

	tail := buf
	var o0 uint64

	// Offset (0) 'Receipts'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 4 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (0) 'Receipts'
	{
		buf = tail[o0:]
		num, err := ssz.DecodeDynamicLength(buf, 16384)
		if err != nil {
			return err
		}
		p.Receipts = make([][]byte, num)
		err = ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
			if len(buf) > 134217728 {
				return ssz.ErrBytesLength
			}
			if cap(p.Receipts[indx]) == 0 {
				p.Receipts[indx] = make([]byte, 0, len(buf))
			}
			p.Receipts[indx] = append(p.Receipts[indx], buf...)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the PortalReceipts object
func (p *PortalReceipts) SizeSSZ() (size int) {
	size = 4

	// Field (0) 'Receipts'
	for ii := 0; ii < len(p.Receipts); ii++ {
		size += 4
		size += len(p.Receipts[ii])
	}

	return
}

// HashTreeRoot ssz hashes the PortalReceipts object
func (p *PortalReceipts) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(p)
}

// HashTreeRootWith ssz hashes the PortalReceipts object with a hasher
func (p *PortalReceipts) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Receipts'
	{
		subIndx := hh.Index()
		num := uint64(len(p.Receipts))
		if num > 16384 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range p.Receipts {
			{
				elemIndx := hh.Index()
				byteLen := uint64(len(elem))
				if byteLen > 134217728 {
					err = ssz.ErrIncorrectListSize
					return
				}
				hh.AppendBytes32(elem)
				hh.MerkleizeWithMixin(elemIndx, byteLen, (134217728+31)/32)
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16384)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the PortalReceipts object
func (p *PortalReceipts) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(p)
}
