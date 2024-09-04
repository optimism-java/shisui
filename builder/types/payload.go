package types

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

var SigningDomainBuilderV1 = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

type PayloadRequestV1 struct {
	Slot       uint64      `json:"slot"`
	ParentHash common.Hash `json:"parentHash"`
}

// MarshalJSON marshals as JSON.
func (p *PayloadRequestV1) MarshalJSON() ([]byte, error) {
	type PayloadRequestV1 struct {
		Slot       hexutil.Uint64 `json:"slot"`
		ParentHash common.Hash    `json:"parentHash"`
	}

	var enc PayloadRequestV1
	enc.Slot = hexutil.Uint64(p.Slot)
	enc.ParentHash = p.ParentHash
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (p *PayloadRequestV1) UnmarshalJSON(input []byte) error {
	type PayloadRequestV1 struct {
		Slot       *hexutil.Uint64 `json:"slot"`
		ParentHash *common.Hash    `json:"parentHash"`
	}

	var dec PayloadRequestV1
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.Slot == nil {
		return errors.New("missing required field 'slot' for PayloadRequestV1")
	}
	p.Slot = uint64(*dec.Slot)
	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for PayloadRequestV1")
	}
	p.ParentHash = *dec.ParentHash
	return nil
}

type BuilderPayloadRequest struct {
	Message   PayloadRequestV1 `json:"message"`
	Signature []byte           `json:"signature"`
}

// MarshalJSON marshals as JSON.
func (p *BuilderPayloadRequest) MarshalJSON() ([]byte, error) {
	type BuilderPayloadRequest struct {
		Message   PayloadRequestV1 `json:"message"`
		Signature hexutil.Bytes    `json:"signature"`
	}

	var enc BuilderPayloadRequest
	enc.Message = p.Message
	enc.Signature = p.Signature
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (p *BuilderPayloadRequest) UnmarshalJSON(input []byte) error {
	type BuilderPayloadRequest struct {
		Message   *PayloadRequestV1 `json:"message"`
		Signature *hexutil.Bytes    `json:"signature"`
	}

	var dec BuilderPayloadRequest
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.Message == nil {
		return errors.New("missing required field 'message' for PayloadRequestV1")
	}
	p.Message = *dec.Message
	if dec.Signature == nil {
		return errors.New("missing required field 'signature' for PayloadRequestV1")
	}
	p.Signature = *dec.Signature
	return nil
}

// BidTrace represents a bid trace.
type BidTrace struct {
	Slot                 uint64         `json:"slot"`
	ParentHash           common.Hash    `json:"parentHash"`
	BlockHash            common.Hash    `json:"blockHash"`
	BuilderAddress       common.Address `json:"builderAddress"`
	ProposerAddress      common.Address `json:"proposerAddress"`
	ProposerFeeRecipient common.Address `json:"proposerFeeRecipient"`
	GasLimit             uint64         `json:"gasLimit"`
	GasUsed              uint64         `json:"gasUsed"`
	Value                *big.Int       `json:"value"`
}

// MarshalJSON marshals as JSON.
func (b *BidTrace) MarshalJSON() ([]byte, error) {
	type BidTrace struct {
		Slot                 hexutil.Uint64 `json:"slot"`
		ParentHash           common.Hash    `json:"parentHash"`
		BlockHash            common.Hash    `json:"blockHash"`
		BuilderAddress       common.Address `json:"builderAddress"`
		ProposerAddress      common.Address `json:"proposerAddress"`
		ProposerFeeRecipient common.Address `json:"proposerFeeRecipient"`
		GasLimit             hexutil.Uint64 `json:"gasLimit"`
		GasUsed              hexutil.Uint64 `json:"gasUsed"`
		Value                *hexutil.Big   `json:"value"`
	}

	var enc BidTrace
	enc.Slot = hexutil.Uint64(b.Slot)
	enc.ParentHash = b.ParentHash
	enc.BlockHash = b.BlockHash
	enc.BuilderAddress = b.BuilderAddress
	enc.ProposerAddress = b.ProposerAddress
	enc.ProposerFeeRecipient = b.ProposerFeeRecipient
	enc.GasLimit = hexutil.Uint64(b.GasLimit)
	enc.GasUsed = hexutil.Uint64(b.GasUsed)
	enc.Value = (*hexutil.Big)(b.Value)
	return json.Marshal(enc)
}

// UnmarshalJSON unmarshals from JSON.
func (b *BidTrace) UnmarshalJSON(input []byte) error {
	type BidTrace struct {
		Slot                 *hexutil.Uint64 `json:"slot"`
		ParentHash           *common.Hash    `json:"parentHash"`
		BlockHash            *common.Hash    `json:"blockHash"`
		BuilderAddress       *common.Address `json:"builderAddress"`
		ProposerAddress      *common.Address `json:"proposerAddress"`
		ProposerFeeRecipient *common.Address `json:"proposerFeeRecipient"`
		GasLimit             *hexutil.Uint64 `json:"gasLimit"`
		GasUsed              *hexutil.Uint64 `json:"gasUsed"`
		Value                *hexutil.Big    `json:"value"`
	}

	var dec BidTrace
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.Slot == nil {
		return errors.New("missing required field 'slot' for BidTrace")
	}
	b.Slot = uint64(*dec.Slot)
	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for BidTrace")
	}
	b.ParentHash = *dec.ParentHash
	if dec.BlockHash == nil {
		return errors.New("missing required field 'blockHash' for BidTrace")
	}
	b.BlockHash = *dec.BlockHash
	if dec.BuilderAddress == nil {
		return errors.New("missing required field 'builderAddress' for BidTrace")
	}
	b.BuilderAddress = *dec.BuilderAddress
	if dec.ProposerAddress == nil {
		return errors.New("missing required field 'proposerAddress' for BidTrace")
	}
	b.ProposerAddress = *dec.ProposerAddress
	if dec.ProposerFeeRecipient == nil {
		return errors.New("missing required field 'proposerFeeRecipient' for BidTrace")
	}
	b.ProposerFeeRecipient = *dec.ProposerFeeRecipient
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for BidTrace")
	}
	b.GasLimit = uint64(*dec.GasLimit)
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for BidTrace")
	}
	b.GasUsed = uint64(*dec.GasUsed)
	if dec.Value == nil {
		return errors.New("missing required field 'value' for BidTrace")
	}
	b.Value = (*big.Int)(dec.Value)
	if b.Value == nil {
		return errors.New("missing required field 'value' for BidTrace")
	}
	return nil
}

// VersionedBuilderPayloadResponse contains a versioned signed builder payload.
type VersionedBuilderPayloadResponse struct {
	Version          SpecVersion            `json:"version"`
	Message          *BidTrace              `json:"message"`
	ExecutionPayload *engine.ExecutableData `json:"executionPayload"`
	Signature        []byte                 `json:"signature"`
}

// MarshalJSON marshals as JSON.
func (p *VersionedBuilderPayloadResponse) MarshalJSON() ([]byte, error) {
	type VersionedBuilderPayloadResponse struct {
		Version          *SpecVersion           `json:"version"`
		Message          *BidTrace              `json:"message"`
		ExecutionPayload *engine.ExecutableData `json:"executionPayload"`
		Signature        hexutil.Bytes          `json:"signature"`
	}

	var enc VersionedBuilderPayloadResponse
	enc.Version = &p.Version
	enc.Message = p.Message
	enc.ExecutionPayload = p.ExecutionPayload
	enc.Signature = p.Signature
	return json.Marshal(enc)
}

// UnmarshalJSON unmarshals from JSON.
func (p *VersionedBuilderPayloadResponse) UnmarshalJSON(input []byte) error {
	type VersionedBuilderPayloadResponse struct {
		Version          *SpecVersion           `json:"version"`
		Message          *BidTrace              `json:"message"`
		ExecutionPayload *engine.ExecutableData `json:"executionPayload"`
		Signature        hexutil.Bytes          `json:"signature"`
	}

	var dec VersionedBuilderPayloadResponse
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.Version == nil {
		return errors.New("missing required field 'version' for VersionedBuilderPayloadResponse")
	}
	p.Version = *dec.Version
	if dec.Message == nil {
		return errors.New("missing required field 'message' for VersionedBuilderPayloadResponse")
	}
	p.Message = dec.Message
	if dec.ExecutionPayload == nil {
		return errors.New("missing required field 'executionPayload' for VersionedBuilderPayloadResponse")
	}
	p.ExecutionPayload = dec.ExecutionPayload
	if dec.Signature == nil {
		return errors.New("missing required field 'signature' for VersionedBuilderPayloadResponse")
	}
	p.Signature = dec.Signature
	return nil
}
