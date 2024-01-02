package history

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/signer/storage"
	"github.com/ethereum/go-ethereum/trie"
	"golang.org/x/crypto/sha3"
)

type ContentType byte

const (
	BlockHeaderType      ContentType = 0x00
	BlockBodyType        ContentType = 0x01
	ReceiptsType         ContentType = 0x02
	EpochAccumulatorType ContentType = 0x03
)

var (
	ErrWithdrawalHashIsNotEqual = errors.New("withdrawals hash is not equal")
	ErrTxHashIsNotEqual         = errors.New("tx hash is not equal")
	ErrUnclesHashIsNotEqual     = errors.New("uncles hash is not equal")
	ErrReceiptsHashIsNotEqual   = errors.New("receipts hash is not equal")
	ErrContentOutOfRange        = errors.New("content out of range")
	ErrHeaderWithProofIsInvalid = errors.New("header proof is invalid")
	ErrInvalidBlockHash = errors.New("invalid block hash")
)

type ContentKey struct {
	selector ContentType
	data     []byte
}

func newContentKey(selector ContentType, hash []byte) *ContentKey {
	return &ContentKey{
		selector: selector,
		data:     hash,
	}
}

func (c *ContentKey) encode() []byte {
	res := make([]byte, 0, len(c.data)+1)
	res = append(res, byte(c.selector))
	res = append(res, c.data...)
	return res
}

type HistoryNetwork struct {
	portalProtocol    *discover.PortalProtocol
	masterAccumulator *MasterAccumulator
}

// FromBlockBodyLegacy convert BlockBodyLegacy to types.Body
func FromBlockBodyLegacy(b *BlockBodyLegacy) (*types.Body, error) {
	transactions := make([]*types.Transaction, 0, len(b.Transactions))
	for _, t := range b.Transactions {
		tran := new(types.Transaction)
		err := rlp.DecodeBytes(t, tran)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tran)
	}
	uncles := make([]*types.Header, 0, len(b.Uncles))
	err := rlp.DecodeBytes(b.Uncles, uncles)
	return &types.Body{
		Uncles:       uncles,
		Transactions: transactions,
	}, err
}

// FromPortalBlockBodyShanghai convert PortalBlockBodyShanghai to types.Body
func FromPortalBlockBodyShanghai(b *PortalBlockBodyShanghai) (*types.Body, error) {
	transactions := make([]*types.Transaction, 0, len(b.Transactions))
	for _, t := range b.Transactions {
		tran := new(types.Transaction)
		err := rlp.DecodeBytes(t, tran)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tran)
	}
	uncles := make([]*types.Header, 0, len(b.Uncles))
	err := rlp.DecodeBytes(b.Uncles, uncles)
	withdrawals := make([]*types.Withdrawal, 0, len(b.Withdrawals))
	for _, w := range b.Withdrawals {
		withdrawal := new(types.Withdrawal)
		err := rlp.DecodeBytes(w, withdrawal)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	return &types.Body{
		Uncles:       uncles,
		Transactions: transactions,
		Withdrawals:  withdrawals,
	}, err
}

// FromPortalReceipts convert PortalReceipts to types.Receipt
func FromPortalReceipts(r *PortalReceipts) ([]*types.Receipt, error) {
	res := make([]*types.Receipt, 0, len(r.Receipts))
	for _, r := range r.Receipts {
		receipt := new(types.Receipt)
		err := rlp.DecodeBytes(r, receipt)
		if err != nil {
			return nil, err
		}
		res = append(res, receipt)
	}
	return res, nil
}

// toBlockBodyLegacy convert types.Body to BlockBodyLegacy
func toBlockBodyLegacy(b *types.Body) (*BlockBodyLegacy, error) {
	txs := make([][]byte, 0, len(b.Transactions))

	for _, tx := range b.Transactions {
		txBytes, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, txBytes)
	}

	uncleBytes, err := rlp.EncodeToBytes(b.Uncles)
	if err != nil {
		return nil, err
	}
	return &BlockBodyLegacy{Uncles: uncleBytes, Transactions: txs}, err
}

// toPortalBlockBodyShanghai convert types.Body to PortalBlockBodyShanghai
func toPortalBlockBodyShanghai(b *types.Body) (*PortalBlockBodyShanghai, error) {
	legacy, err := toBlockBodyLegacy(b)
	if err != nil {
		return nil, err
	}
	withdrawals := make([][]byte, 0, len(b.Withdrawals))
	for _, w := range b.Withdrawals {
		b, err := rlp.EncodeToBytes(w)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, b)
	}
	return &PortalBlockBodyShanghai{Transactions: legacy.Transactions, Uncles: legacy.Uncles, Withdrawals: withdrawals}, nil
}

// ToPortalReceipts convert types.Receipt to PortalReceipts
func ToPortalReceipts(receipts []*types.Receipt) (*PortalReceipts, error) {
	res := make([][]byte, 0, len(receipts))
	for _, r := range receipts {
		b, err := rlp.EncodeToBytes(r)
		if err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return &PortalReceipts{Receipts: res}, nil
}

// EncodeBlockBody encode types.Body to ssz bytes
func EncodeBlockBody(body *types.Body) ([]byte, error) {
	if body.Withdrawals != nil && len(body.Withdrawals) > 0 {
		blockShanghai, err := toPortalBlockBodyShanghai(body)
		if err != nil {
			return nil, err
		}
		return blockShanghai.MarshalSSZ()
	} else {
		legacyBlock, err := toBlockBodyLegacy(body)
		if err != nil {
			return nil, err
		}
		return legacyBlock.MarshalSSZ()
	}
}

// EncodeReceipts encode []*types.Receipt to ssz bytes
func EncodeReceipts(receipts []*types.Receipt) ([]byte, error) {
	portalReceipts, err := ToPortalReceipts(receipts)
	if err != nil {
		return nil, err
	}
	return portalReceipts.MarshalSSZ()
}

func ValidateBlockHeaderBytes(headerBytes []byte, blockHash []byte) (*types.Header, error) {
	header := new(types.Header)
	err := rlp.DecodeBytes(headerBytes, header)
	if err != nil {
		return nil, err
	}
	if header.ExcessBlobGas != nil {
		return nil, errors.New("EIP-4844 not yet implemented")
	}
	hash := header.Hash()
	if !bytes.Equal(hash[:], blockHash) {
		return nil, ErrInvalidBlockHash
	}
	return header, nil
}

func validateBlockBodyLegacy(body *BlockBodyLegacy, header *types.Header) error {
	keccak := sha3.NewLegacyKeccak256()
	uncleHash := keccak.Sum(body.Uncles)
	if !bytes.Equal(uncleHash, header.UncleHash[:]) {
		return ErrUnclesHashIsNotEqual
	}
	keccak.Reset()
	for _, tx := range body.Transactions {
		_, err := keccak.Write(tx)
		if err != nil {
			return err
		}
	}
	txHash := keccak.Sum(nil)
	if !bytes.Equal(txHash, header.TxHash.Bytes()) {
		return ErrTxHashIsNotEqual
	}
	return nil
}

func validatePortalBlockBodyShanghai(body *PortalBlockBodyShanghai, header *types.Header) error {
	legacy := &BlockBodyLegacy{
		Transactions: body.Transactions,
		Uncles:       body.Uncles,
	}
	if err := validateBlockBodyLegacy(legacy, header); err != nil {
		return err
	}
	keccak := sha3.NewLegacyKeccak256()
	for _, withdrawal := range body.Withdrawals {
		_, err := keccak.Write(withdrawal)
		if err != nil {
			return err
		}
	}
	if !bytes.Equal(keccak.Sum(nil), header.WithdrawalsHash.Bytes()) {
		return ErrWithdrawalHashIsNotEqual
	}
	return nil
}

func validateBlockBody(body *types.Body, header *types.Header) error {
	if hash := types.CalcUncleHash(body.Uncles); !bytes.Equal(hash[:], header.UncleHash.Bytes()) {
		return ErrUnclesHashIsNotEqual
	}

	if hash := types.DeriveSha(types.Transactions(body.Transactions), trie.NewStackTrie(nil)); !bytes.Equal(hash[:], header.TxHash.Bytes()) {
		return ErrTxHashIsNotEqual
	}
	if body.Withdrawals == nil {
		return nil
	}
	if hash := types.DeriveSha(types.Withdrawals(body.Withdrawals), trie.NewStackTrie(nil)); !bytes.Equal(hash[:], header.WithdrawalsHash.Bytes()) {
		return ErrWithdrawalHashIsNotEqual
	}
	return nil
}

func DecodePortalBlockBodyBytes(bodyBytes []byte) (*types.Body, error) {
	blockBodyShanghai := new(PortalBlockBodyShanghai)
	err := blockBodyShanghai.UnmarshalSSZ(bodyBytes)
	if err == nil {
		return FromPortalBlockBodyShanghai(blockBodyShanghai)
	}

	blockBodyLegacy := new(BlockBodyLegacy)
	err = blockBodyLegacy.UnmarshalSSZ(bodyBytes)
	if err == nil {
		return FromBlockBodyLegacy(blockBodyLegacy)
	}
	return nil, errors.New("all portal block body decodings failed")
}

func ValidateBlockBodyBytes(bodyBytes []byte, header *types.Header) (*types.Body, error) {
	// TODO check shanghai, pos and legacy block
	body, err := DecodePortalBlockBodyBytes(bodyBytes)
	if err != nil {
		return nil, err
	}
	err = validateBlockBody(body, header)
	return body, err
}

func caclRootHash(input [][]byte) ([]byte, error) {
	keccak := sha3.NewLegacyKeccak256()
	for _, i := range input {
		_, err := keccak.Write(i)
		if err != nil {
			return nil, err
		}
	}
	hash := keccak.Sum(nil)
	return hash, nil
}

func ValidatePortalReceipts(receipt *PortalReceipts, receiptsRoot []byte) error {
	hash, err := caclRootHash(receipt.Receipts)
	if err != nil {
		return err
	}
	if !bytes.Equal(hash, receiptsRoot) {
		return ErrReceiptsHashIsNotEqual
	}
	return nil
}

func ValidatePortalReceiptsBytes(receiptBytes, receiptsRoot []byte) ([]*types.Receipt, error) {
	portalReceipts := new(PortalReceipts)
	err := portalReceipts.UnmarshalSSZ(receiptBytes)
	if err != nil {
		return nil, err
	}
	return FromPortalReceipts(portalReceipts)
}

func NewHistoryNetwork(portalProtocol *discover.PortalProtocol, accu *MasterAccumulator) *HistoryNetwork {
	return &HistoryNetwork{
		portalProtocol:    portalProtocol,
		masterAccumulator: accu,
	}
}

func (h *HistoryNetwork) Start() error {
	err := h.portalProtocol.Start()
	if err != nil {
		return err
	}
	go h.processContentLoop()
	return nil
}

func (h *HistoryNetwork) verifyHeader(header *types.Header, proof BlockHeaderProof) (bool, error) {
	return h.masterAccumulator.VerifyHeader(*header, proof)
}

// Currently doing 4 retries on lookups but only when the validation fails.
const requestRetries = 4

func (h *HistoryNetwork) GetBlockHeader(blockHash []byte) (*types.Header, error) {
	contentKey := newContentKey(BlockHeaderType, blockHash).encode()
	contentId := h.portalProtocol.ToContentId(contentKey)
	if !h.portalProtocol.InRange(contentId) {
		return nil, ErrContentOutOfRange
	}

	res, err := h.portalProtocol.Get(contentId)
	// other error
	if err != nil && err != storage.ErrNotFound {
		return nil, err
	}
	// no error
	if err == nil {
		blockHeaderWithProof, err := DecodeBlockHeaderWithProof(res)
		if err != nil {
			return nil, err
		}
		header := new(types.Header)
		err = rlp.DecodeBytes(blockHeaderWithProof.Header, header)
		return header, err
	}
	// no content in local storage

	for retries := 0; retries < requestRetries; retries++ {
		content, err := h.portalProtocol.ContentLookup(contentKey)
		if err != nil {
			return nil, err
		}

		headerWithProof, err := DecodeBlockHeaderWithProof(res)
		if err != nil {
			continue
		}

		header, err := ValidateBlockHeaderBytes(headerWithProof.Header, blockHash)
		if err != nil {
			continue
		}
		valid, err := h.verifyHeader(header, *headerWithProof.Proof)
		if err != nil || !valid {
			continue
		}
		// TODO handle the error
		_ = h.portalProtocol.Put(contentId, content)
		return header, nil
	}
	return nil, ErrContentOutOfRange
}

func (h *HistoryNetwork) GetBlockBody(blockHash []byte) (*types.Body, error) {
	header, err := h.GetBlockHeader(blockHash)
	if err != nil {
		return nil, err
	}
	contentKey := newContentKey(BlockBodyType, blockHash).encode()
	contentId := h.portalProtocol.ToContentId(contentKey)

	if !h.portalProtocol.InRange(contentId) {
		return nil, ErrContentOutOfRange
	}

	res, err := h.portalProtocol.Get(contentId)
	// other error
	if err != nil && err != storage.ErrNotFound {
		return nil, err
	}
	// no error
	if err == nil {
		body, err := DecodePortalBlockBodyBytes(res)
		return body, err
	}
	// no content in local storage

	for retries := 0; retries < requestRetries; retries++ {
		content, err := h.portalProtocol.ContentLookup(contentKey)
		if err != nil {
			return nil, err
		}
		body, err := DecodePortalBlockBodyBytes(res)
		if err != nil {
			continue
		}

		err = validateBlockBody(body, header)
		if err != nil {
			continue
		}
		// TODO handle the error
		_ = h.portalProtocol.Put(contentId, content)
		return body, nil
	}
	return nil, ErrContentOutOfRange
}

func (h *HistoryNetwork) GetReceipts(blockHash []byte) ([]*types.Receipt, error) {
	header, err := h.GetBlockHeader(blockHash)
	if err != nil {
		return nil, err
	}
	contentKey := newContentKey(ReceiptsType, blockHash).encode()
	contentId := h.portalProtocol.ToContentId(contentKey)

	if !h.portalProtocol.InRange(contentId) {
		return nil, ErrContentOutOfRange
	}

	res, err := h.portalProtocol.Get(contentId)
	// other error
	if err != nil && err != storage.ErrNotFound {
		return nil, err
	}
	// no error
	if err == nil {
		portalReceipte := new(PortalReceipts)
		err := portalReceipte.UnmarshalSSZ(res)
		if err != nil {
			return nil, err
		}
		receipts, err := FromPortalReceipts(portalReceipte)
		return receipts, err
	}
	// no content in local storage

	for retries := 0; retries < requestRetries; retries++ {
		content, err := h.portalProtocol.ContentLookup(contentKey)
		if err != nil {
			return nil, err
		}
		receipts, err := ValidatePortalReceiptsBytes(content, header.ReceiptHash.Bytes())
		if err != nil {
			continue
		}
		// TODO handle the error
		_ = h.portalProtocol.Put(contentId, content)
		return receipts, nil
	}
	return nil, ErrContentOutOfRange
}

func decodeEpochAccumulator(data []byte) (*EpochAccumulator, error) {
	epochAccu := new(EpochAccumulator)
	err := epochAccu.UnmarshalSSZ(data)
	return epochAccu, err
}

func (h *HistoryNetwork) GetEpochAccumulator(epochHash []byte) (*EpochAccumulator, error) {
	contentKey := newContentKey(EpochAccumulatorType, epochHash).encode()
	contentId := h.portalProtocol.ToContentId(contentKey)

	res, err := h.portalProtocol.Get(contentId)
	// other error
	if err != nil && err != storage.ErrNotFound {
		return nil, err
	}
	// no error
	if err == nil {
		epochAccu, err := decodeEpochAccumulator(res)
		return epochAccu, err
	}
	for retries := 0; retries < requestRetries; retries++ {
		content, err := h.portalProtocol.ContentLookup(contentKey)
		if err != nil {
			return nil, err
		}
		epochAccu, err := decodeEpochAccumulator(res)
		if err != nil {
			continue
		}
		hash, err := epochAccu.HashTreeRoot()
		if err != nil {
			continue
		}
		mixHash := MixInLength(hash, epochSize)
		if !bytes.Equal(mixHash, epochHash) {
			continue
		}
		// TODO handle the error
		_ = h.portalProtocol.Put(contentId, content)
		return epochAccu, nil
	}
	return nil, ErrContentOutOfRange
}

func DecodeBlockHeaderWithProof(content []byte) (*BlockHeaderWithProof, error) {
	headerWithProof := new(BlockHeaderWithProof)
	err := headerWithProof.UnmarshalSSZ(content)
	return headerWithProof, err
}

func (h *HistoryNetwork) validateContent(contentKey []byte, content []byte) error {

	switch ContentType(contentKey[0]) {
	case BlockHeaderType:
		headerWithProof, err := DecodeBlockHeaderWithProof(content)
		if err != nil {
			return err
		}
		header, err := ValidateBlockHeaderBytes(headerWithProof.Header, contentKey[1:])
		if err != nil {
			return err
		}
		valid, err := h.verifyHeader(header, *headerWithProof.Proof)
		if err != nil {
			return err
		}
		if !valid {
			return ErrHeaderWithProofIsInvalid
		}
		return err
	case BlockBodyType:
		header, err := h.GetBlockHeader(contentKey[1:])
		if err != nil {
			return err
		}
		_, err = ValidateBlockBodyBytes(content, header)
		return err
	case ReceiptsType:
		header, err := h.GetBlockHeader(contentKey[1:])
		if err != nil {
			return err
		}
		_, err = ValidatePortalReceiptsBytes(content, header.ReceiptHash.Bytes())
		return err
	case EpochAccumulatorType:
		if !h.masterAccumulator.Contains(contentKey[1:]) {
			return errors.New("epoch hash is not existed")
		}

		epochAcc, err := decodeEpochAccumulator(content)
		if err != nil {
			return err
		}
		hash, err := epochAcc.HashTreeRoot()
		if err != nil {
			return err
		}

		epochHash := MixInLength(hash, epochSize)
		if !bytes.Equal(contentKey[1:], epochHash) {
			return errors.New("epoch accumulator has invalid root hash")
		}
		return nil
	}
	return errors.New("unknown content type")
}

func (h *HistoryNetwork) validateContents(contentKeys [][]byte, contents [][]byte) error {
	for i, content := range contents {
		contentKey := contentKeys[i]
		err := h.validateContent(contentKey, content)
		if err != nil {
			return fmt.Errorf("content validate failed with content key %v", contentKey)
		}
		contentId := h.portalProtocol.ToContentId(contentKey)
		_ = h.portalProtocol.Put(contentId, content)
	}
	return nil
}

func (h *HistoryNetwork) processContentLoop() {
	contentChan := h.portalProtocol.GetContent()
	for contentElement := range contentChan {
		err := h.validateContents(contentElement.ContentKeys, contentElement.Contents)
		if err != nil {
			continue
		}
		// TODO gossip the validate content
	}
}
