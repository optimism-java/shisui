package builder

import (
	"math/big"
	"testing"
	"time"

	builderApiV1 "github.com/attestantio/go-builder-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/flashbots/go-boost-utils/bls"
	"github.com/flashbots/go-boost-utils/ssz"
	"github.com/flashbots/go-boost-utils/utils"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

type testEthereumService struct {
	synced             bool
	testExecutableData *engine.ExecutionPayloadEnvelope
	testBlock          *types.Block
}

func (t *testEthereumService) BuildBlock(attrs *BuilderPayloadAttributes) (*engine.ExecutionPayloadEnvelope, error) {
	return t.testExecutableData, nil
}

func (t *testEthereumService) GetBlockByHash(hash common.Hash) *types.Block { return t.testBlock }

func (t *testEthereumService) Config() *params.ChainConfig {
	zeroTime := uint64(0)
	config := params.TestChainConfig
	config.CanyonTime = &zeroTime
	config.EcotoneTime = &zeroTime
	return config
}

func (t *testEthereumService) Synced() bool { return t.synced }

func TestGetPayloadV1(t *testing.T) {
	const (
		validatorDesiredGasLimit = 30_000_000
		payloadAttributeGasLimit = 0
		parentBlockGasLimit      = 29_000_000
	)
	expectedGasLimit := core.CalcGasLimit(parentBlockGasLimit, validatorDesiredGasLimit)

	feeRecipient, err := utils.HexToAddress("0xabcf8e0d4e9587369b2301d0790347320302cc00")
	require.NoError(t, err)

	sk, err := bls.SecretKeyFromBytes(hexutil.MustDecode("0x31ee185dad1220a8c88ca5275e64cf5a5cb09cb621cb30df52c9bee8fbaaf8d7"))
	require.NoError(t, err)

	bDomain := ssz.ComputeDomain(ssz.DomainTypeAppBuilder, [4]byte{0x02, 0x0, 0x0, 0x0}, phase0.Root{})

	zero := uint64(0)
	blockHash := common.HexToHash("5fc0137650b887cdb47ae6426d5e0e368e315a2ad3f93df80a8d14f8e1ce239a")
	testExecutableData := &engine.ExecutionPayloadEnvelope{
		ExecutionPayload: &engine.ExecutableData{
			ParentHash:   common.Hash{0x02, 0x03},
			FeeRecipient: common.Address(feeRecipient),
			StateRoot:    common.Hash{0x07, 0x16},
			ReceiptsRoot: common.Hash{0x08, 0x20},
			LogsBloom:    types.Bloom{}.Bytes(),
			Number:       uint64(10),
			GasLimit:     expectedGasLimit,
			GasUsed:      uint64(100),
			Timestamp:    uint64(105),
			ExtraData:    hexutil.MustDecode("0x0042fafc"),

			BaseFeePerGas: big.NewInt(16),

			BlockHash:     blockHash,
			Transactions:  [][]byte{},
			Withdrawals:   types.Withdrawals{},
			BlobGasUsed:   &zero,
			ExcessBlobGas: &zero,
		},
		BlockValue:  big.NewInt(10),
		BlobsBundle: &engine.BlobsBundleV1{},
	}

	testBlock, err := engine.ExecutableDataToBlock(*testExecutableData.ExecutionPayload, nil, nil)
	require.NoError(t, err)

	testPayloadAttributes := &BuilderPayloadAttributes{
		Timestamp:             hexutil.Uint64(104),
		Random:                common.Hash{0x05, 0x10},
		SuggestedFeeRecipient: common.Address(feeRecipient),
		GasLimit:              uint64(payloadAttributeGasLimit),
		Slot:                  uint64(25),
		HeadHash:              common.Hash{0x02, 0x03},

		NoTxPool:     false,
		Transactions: []*types.Transaction{},
	}

	testEthService := &testEthereumService{
		synced:             true,
		testExecutableData: testExecutableData,
		testBlock:          testBlock,
	}

	mockBeaconNode := newMockBeaconNode(*testPayloadAttributes)
	builderArgs := BuilderArgs{
		sk:                          sk,
		builderSigningDomain:        bDomain,
		builderRetryInterval:        200 * time.Millisecond,
		blockTime:                   2 * time.Second,
		eth:                         testEthService,
		ignoreLatePayloadAttributes: false,
		beaconClient:                mockBeaconNode,
	}
	builder, err := NewBuilder(builderArgs)
	require.NoError(t, err)
	builder.Start()
	defer builder.Stop()
	time.Sleep(2 * time.Second)

	blockResponse, err := builder.GetPayload(
		PayloadRequestV1{
			Slot:       testPayloadAttributes.Slot,
			ParentHash: testPayloadAttributes.HeadHash,
		},
	)
	require.NoError(t, err)

	expectedMessage := &builderApiV1.BidTrace{
		Slot:                 uint64(25),
		BlockHash:            phase0.Hash32(blockHash),
		ParentHash:           phase0.Hash32{0x02, 0x03},
		BuilderPubkey:        builder.builderPublicKey,
		ProposerPubkey:       phase0.BLSPubKey{},
		ProposerFeeRecipient: feeRecipient,
		GasLimit:             expectedGasLimit,
		GasUsed:              uint64(100),
		Value:                &uint256.Int{0x0a},
	}
	copy(expectedMessage.BlockHash[:], blockHash[:])
	require.NotNil(t, blockResponse.Deneb)
	require.Equal(t, expectedMessage, blockResponse.Deneb.Message)

	expectedExecutionPayload := &deneb.ExecutionPayload{
		ParentHash:    phase0.Hash32{0x02, 0x03},
		FeeRecipient:  feeRecipient,
		StateRoot:     [32]byte(testExecutableData.ExecutionPayload.StateRoot),
		ReceiptsRoot:  [32]byte(testExecutableData.ExecutionPayload.ReceiptsRoot),
		LogsBloom:     [256]byte{},
		PrevRandao:    [32]byte(testExecutableData.ExecutionPayload.Random),
		BlockNumber:   testExecutableData.ExecutionPayload.Number,
		GasLimit:      testExecutableData.ExecutionPayload.GasLimit,
		GasUsed:       testExecutableData.ExecutionPayload.GasUsed,
		Timestamp:     testExecutableData.ExecutionPayload.Timestamp,
		ExtraData:     hexutil.MustDecode("0x0042fafc"),
		BaseFeePerGas: uint256.NewInt(16),
		BlockHash:     expectedMessage.BlockHash,
		Transactions:  []bellatrix.Transaction{},
		Withdrawals:   []*capella.Withdrawal{},
	}

	require.Equal(t, expectedExecutionPayload, blockResponse.Deneb.ExecutionPayload)

	expectedSignature, err := utils.HexToSignature("0xa2416ce0f5d65329c750ea9338f9b11280aa9b04180859a8a9e8082d1794b9fa17f7928bfab66ed20733d06e6f9ac2ec185ff98b91e5da9aaa6fbab8a410c9f9b76e5fde01d498ee2cfe9867949d2eead6a53e953ad910805f8dddae3ec3bbdf")
	require.NoError(t, err)
	require.Equal(t, expectedSignature, blockResponse.Deneb.Signature)
}
