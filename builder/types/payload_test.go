package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestBuilderPayloadRequest_MarshalJSON(t *testing.T) {
	signature := [65]byte{0x1, 0x2}
	payload := BuilderPayloadRequest{
		Message: PayloadRequestV1{
			Slot:       1,
			ParentHash: common.Hash{0x2, 03},
		},
		Signature: signature[:],
	}
	json, err := payload.MarshalJSON()
	require.NoError(t, err)
	payloadStr := `{"message":{"slot":"0x1","parentHash":"0x0203000000000000000000000000000000000000000000000000000000000000"},"signature":"0x0102000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"}`
	require.Equal(t, payloadStr, string(json))
}

func TestBuilderPayloadRequest_UnmarshalJSON(t *testing.T) {
	payloadStr := `{"message":{"slot":"0xa","parentHash":"0x0405000000000000000000000000000000000000000000000000000000000000"},"signature":"0x0203000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"}`
	payload := BuilderPayloadRequest{}
	signature := [65]byte{0x2, 0x3}
	err := payload.UnmarshalJSON([]byte(payloadStr))
	require.NoError(t, err)
	require.Equal(t, payload.Message.Slot, uint64(10))
	require.Equal(t, payload.Message.ParentHash, common.Hash{0x4, 05})
	require.Equal(t, payload.Signature, signature[:])
}

func TestVersionedBuilderPayloadResponse_MarshalJSON(t *testing.T) {
	bidTrace := BidTrace{
		Slot:                 uint64(1),
		ParentHash:           common.Hash{0x1},
		BlockHash:            common.Hash{0x2},
		BuilderAddress:       common.Address{0x3},
		ProposerAddress:      common.Address{0x4},
		ProposerFeeRecipient: common.Address{0x5},
		GasLimit:             uint64(100),
		GasUsed:              uint64(90),
		Value:                big.NewInt(30),
	}
	zero := uint64(0)
	executionPayload := &engine.ExecutableData{
		ParentHash:   common.Hash{0x02, 0x03},
		FeeRecipient: common.Address{0x04, 0x05},
		StateRoot:    common.Hash{0x07, 0x16},
		ReceiptsRoot: common.Hash{0x08, 0x20},
		Random:       common.Hash{0x09, 0x30},
		LogsBloom:    types.Bloom{}.Bytes(),
		Number:       uint64(10),
		GasLimit:     uint64(100),
		GasUsed:      uint64(100),
		Timestamp:    uint64(105),
		ExtraData:    hexutil.MustDecode("0x0042fafc"),

		BaseFeePerGas: big.NewInt(16),

		BlockHash:     common.Hash{0x09, 0x30},
		Transactions:  [][]byte{},
		Withdrawals:   types.Withdrawals{},
		BlobGasUsed:   &zero,
		ExcessBlobGas: &zero,
	}
	versionedBuilderPayloadResponse := VersionedBuilderPayloadResponse{
		Version:          SpecVersionEcotone,
		Message:          &bidTrace,
		ExecutionPayload: executionPayload,
		Signature:        []byte{0x1, 0x2},
	}

	json, err := versionedBuilderPayloadResponse.MarshalJSON()
	require.NoError(t, err)
	expectedStr := `{"version":"ecotone","message":{"slot":"0x1","parentHash":"0x0100000000000000000000000000000000000000000000000000000000000000","blockHash":"0x0200000000000000000000000000000000000000000000000000000000000000","builderAddress":"0x0300000000000000000000000000000000000000","proposerAddress":"0x0400000000000000000000000000000000000000","proposerFeeRecipient":"0x0500000000000000000000000000000000000000","gasLimit":"0x64","gasUsed":"0x5a","value":"0x1e"},"executionPayload":{"parentHash":"0x0203000000000000000000000000000000000000000000000000000000000000","feeRecipient":"0x0405000000000000000000000000000000000000","stateRoot":"0x0716000000000000000000000000000000000000000000000000000000000000","receiptsRoot":"0x0820000000000000000000000000000000000000000000000000000000000000","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","prevRandao":"0x0930000000000000000000000000000000000000000000000000000000000000","blockNumber":"0xa","gasLimit":"0x64","gasUsed":"0x64","timestamp":"0x69","extraData":"0x0042fafc","baseFeePerGas":"0x10","blockHash":"0x0930000000000000000000000000000000000000000000000000000000000000","transactions":[],"withdrawals":[],"blobGasUsed":"0x0","excessBlobGas":"0x0"},"signature":"0x0102"}`
	require.Equal(t, expectedStr, string(json))
}

func TestVersionedBuilderPayloadResponse_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"version":"canyon","message":{"slot":"0xa","parentHash":"0x0200000000000000000000000000000000000000000000000000000000000000","blockHash":"0x0300000000000000000000000000000000000000000000000000000000000000","builderAddress":"0x0300000000000000000000000000000000000000","proposerAddress":"0x0400000000000000000000000000000000000000","proposerFeeRecipient":"0x0500000000000000000000000000000000000000","gasLimit":"0x64","gasUsed":"0x5a","value":"0x1e"},"executionPayload":{"parentHash":"0x0203000000000000000000000000000000000000000000000000000000000000","feeRecipient":"0x0405000000000000000000000000000000000000","stateRoot":"0x0716000000000000000000000000000000000000000000000000000000000000","receiptsRoot":"0x0820000000000000000000000000000000000000000000000000000000000000","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","prevRandao":"0x0930000000000000000000000000000000000000000000000000000000000000","blockNumber":"0xa","gasLimit":"0x64","gasUsed":"0x64","timestamp":"0x69","extraData":"0x0042fafc","baseFeePerGas":"0x10","blockHash":"0x0930000000000000000000000000000000000000000000000000000000000000","transactions":[],"withdrawals":[],"blobGasUsed":"0x0","excessBlobGas":"0x0"},"signature":"0x0102"}`
	versionedBuilderPayloadResponse := VersionedBuilderPayloadResponse{}
	err := versionedBuilderPayloadResponse.UnmarshalJSON([]byte(jsonStr))
	require.NoError(t, err)

	bidTrace := BidTrace{
		Slot:                 uint64(10),
		ParentHash:           common.Hash{0x2},
		BlockHash:            common.Hash{0x3},
		BuilderAddress:       common.Address{0x3},
		ProposerAddress:      common.Address{0x4},
		ProposerFeeRecipient: common.Address{0x5},
		GasLimit:             uint64(100),
		GasUsed:              uint64(90),
		Value:                big.NewInt(30),
	}
	zero := uint64(0)
	executionPayload := &engine.ExecutableData{
		ParentHash:   common.Hash{0x02, 0x03},
		FeeRecipient: common.Address{0x04, 0x05},
		StateRoot:    common.Hash{0x07, 0x16},
		ReceiptsRoot: common.Hash{0x08, 0x20},
		Random:       common.Hash{0x09, 0x30},
		LogsBloom:    types.Bloom{}.Bytes(),
		Number:       uint64(10),
		GasLimit:     uint64(100),
		GasUsed:      uint64(100),
		Timestamp:    uint64(105),
		ExtraData:    hexutil.MustDecode("0x0042fafc"),

		BaseFeePerGas: big.NewInt(16),

		BlockHash:     common.Hash{0x09, 0x30},
		Transactions:  [][]byte{},
		Withdrawals:   types.Withdrawals{},
		BlobGasUsed:   &zero,
		ExcessBlobGas: &zero,
	}
	expectedVersionedBuilderPayloadResponse := VersionedBuilderPayloadResponse{
		Version:          SpecVersionCanyon,
		Message:          &bidTrace,
		ExecutionPayload: executionPayload,
		Signature:        []byte{0x1, 0x2},
	}
	require.Equal(t, expectedVersionedBuilderPayloadResponse, versionedBuilderPayloadResponse)
}
