package builder

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	builderTypes "github.com/ethereum/go-ethereum/builder/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/params"
)

type IEthereumService interface {
	BuildBlock(attrs *builderTypes.PayloadAttributes) (*engine.ExecutionPayloadEnvelope, error)
	GetBlockByHash(hash common.Hash) *types.Block
	Config() *params.ChainConfig
	Synced() bool
}

type EthereumService struct {
	eth *eth.Ethereum
	cfg *Config
}

func NewEthereumService(eth *eth.Ethereum, config *Config) *EthereumService {
	return &EthereumService{
		eth: eth,
		cfg: config,
	}
}

func (s *EthereumService) BuildBlock(attrs *builderTypes.PayloadAttributes) (*engine.ExecutionPayloadEnvelope, error) {
	// Send a request to generate a full block in the background.
	// The result can be obtained via the returned channel.
	args := &miner.BuildPayloadArgs{
		Parent:       attrs.HeadHash,
		Timestamp:    uint64(attrs.Timestamp),
		FeeRecipient: attrs.SuggestedFeeRecipient, // TODO (builder): use builder key as fee recipient
		GasLimit:     &attrs.GasLimit,
		Random:       attrs.Random,
		Withdrawals:  attrs.Withdrawals,
		BeaconRoot:   attrs.ParentBeaconBlockRoot,
		Transactions: attrs.Transactions,
		NoTxPool:     attrs.NoTxPool,
	}

	payload, err := s.eth.Miner().BuildPayload(args)
	if err != nil {
		log.Error("Failed to build payload", "err", err)
		return nil, err
	}

	resCh := make(chan *engine.ExecutionPayloadEnvelope, 1)
	go func() {
		resCh <- payload.ResolveFull()
	}()

	timer := time.NewTimer(s.cfg.BlockTime)
	defer timer.Stop()

	select {
	case payload := <-resCh:
		if payload == nil {
			return nil, errors.New("received nil payload from sealing work")
		}
		return payload, nil
	case <-timer.C:
		payload.Cancel()
		log.Error("timeout waiting for block", "parent hash", attrs.HeadHash, "slot", attrs.Slot)
		return nil, errors.New("timeout waiting for block result")
	}
}

func (s *EthereumService) GetBlockByHash(hash common.Hash) *types.Block {
	return s.eth.BlockChain().GetBlockByHash(hash)
}

func (s *EthereumService) Config() *params.ChainConfig {
	return s.eth.BlockChain().Config()
}

func (s *EthereumService) Synced() bool {
	return s.eth.Synced()
}
