package builder

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	builderTypes "github.com/ethereum/go-ethereum/builder/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

type testEthereumService struct {
	synced             bool
	testExecutableData *engine.ExecutionPayloadEnvelope
	testBlock          *types.Block
}

func (t *testEthereumService) BuildBlock(attrs *builderTypes.PayloadAttributes) (*engine.ExecutionPayloadEnvelope, error) {
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
		testPrivateKeyHex        = "2fc12ae741f29701f8e30f5de6350766c020cb80768a0ff01e6838ffd2431e11"
	)

	testPrivateKey, err := crypto.HexToECDSA(testPrivateKeyHex)
	require.NoError(t, err)
	feeRecipient := common.HexToAddress("0xabcf8e0d4e9587369b2301d0790347320302cc00")
	expectedGasLimit := core.CalcGasLimit(parentBlockGasLimit, validatorDesiredGasLimit)
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

	testPayloadAttributes := &builderTypes.PayloadAttributes{
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
		builderPrivateKey:           testPrivateKey,
		builderAddress:              crypto.PubkeyToAddress(testPrivateKey.PublicKey),
		proposerAddress:             crypto.PubkeyToAddress(testPrivateKey.PublicKey),
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

	msg := builderTypes.PayloadRequestV1{
		Slot:       testPayloadAttributes.Slot,
		ParentHash: testPayloadAttributes.HeadHash,
	}
	requestBytes, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	hash, err := BuilderSigningHash(builder.eth.Config(), requestBytes)
	require.NoError(t, err)
	signature, err := crypto.Sign(hash[:], testPrivateKey)
	require.NoError(t, err)
	blockResponse, err := builder.GetPayload(
		&builderTypes.BuilderPayloadRequest{
			Message:   msg,
			Signature: signature,
		},
	)
	require.NoError(t, err)

	expectedMessage := &builderTypes.BidTrace{
		Slot:                 uint64(25),
		BlockHash:            blockHash,
		ParentHash:           common.Hash{0x02, 0x03},
		BuilderAddress:       builderArgs.builderAddress,
		ProposerAddress:      builderArgs.proposerAddress,
		ProposerFeeRecipient: feeRecipient,
		GasLimit:             expectedGasLimit,
		GasUsed:              uint64(100),
		Value:                big.NewInt(10),
	}
	copy(expectedMessage.BlockHash[:], blockHash[:])
	require.NotNil(t, blockResponse)
	require.Equal(t, expectedMessage, blockResponse.Message)

	require.Equal(t, testExecutableData.ExecutionPayload, blockResponse.ExecutionPayload)
	expectedSignature, err := hexutil.Decode("0x72ce97e705ecbe69effb11ae63791fd45f058e5766be52ad6704516d5a984ebf00a22d8661201e21968fbf285700d6ab93accfe54acf2dfc4c52edc190a54fdb00")
	require.NoError(t, err)
	require.Equal(t, expectedSignature, blockResponse.Signature)
}
