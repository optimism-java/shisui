package builder

import (
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"time"

	builderTypes "github.com/ethereum/go-ethereum/builder/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
)

const (
	_PathGetPayload = "/eth/v1/builder/payload"
)

type Service struct {
	srv     *http.Server
	builder IBuilder
}

func (s *Service) Start() error {
	if s.srv != nil {
		log.Info("Service started")
		go s.srv.ListenAndServe()
	}

	s.builder.Start()

	return nil
}

func (s *Service) Stop() error {
	if s.srv != nil {
		s.srv.Close()
	}
	s.builder.Stop()
	return nil
}

func (s *Service) GetPayloadV1(request *builderTypes.BuilderPayloadRequest) (*builderTypes.VersionedBuilderPayloadResponse, error) {
	return s.builder.GetPayload(request)
}

func NewService(listenAddr string, builder IBuilder) *Service {
	var srv *http.Server

	router := mux.NewRouter()
	router.HandleFunc(_PathGetPayload, builder.handleGetPayload).Methods(http.MethodPost)

	srv = &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	return &Service{
		srv:     srv,
		builder: builder,
	}
}

func Register(stack *node.Node, backend *eth.Ethereum, cfg *Config) error {
	var beaconClient IBeaconClient
	if len(cfg.BeaconEndpoints) == 0 {
		beaconClient = &NilBeaconClient{}
	} else {
		beaconClient = NewOpBeaconClient(cfg.BeaconEndpoints[0])
	}

	ethereumService := NewEthereumService(backend, cfg)

	var builderRetryInterval time.Duration
	if cfg.RetryInterval != "" {
		d, err := time.ParseDuration(cfg.RetryInterval)
		if err != nil {
			return fmt.Errorf("error parsing builder retry interval - %v", err)
		}
		builderRetryInterval = d
	} else {
		builderRetryInterval = RetryIntervalDefault
	}

	builderPrivateKey, err := crypto.HexToECDSA(cfg.BuilderSigningKey)
	if err != nil {
		return fmt.Errorf("invalid builder private key: %w", err)
	}
	builderPublicKey := builderPrivateKey.Public()
	builderPublicKeyECDSA, ok := builderPublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("publicKey could not be converted to ECDSA")
	}
	builderAddress := crypto.PubkeyToAddress(*builderPublicKeyECDSA)

	var proposerAddress common.Address
	if common.IsHexAddress(cfg.ProposerAddress) {
		proposerAddress = common.HexToAddress(cfg.ProposerAddress)
	} else {
		log.Warn("proposer signing address is invalid or not set, proposer signature verification will be skipped")
	}

	builderArgs := BuilderArgs{
		builderPrivateKey:           builderPrivateKey,
		builderAddress:              builderAddress,
		proposerAddress:             proposerAddress,
		eth:                         ethereumService,
		builderRetryInterval:        builderRetryInterval,
		ignoreLatePayloadAttributes: cfg.IgnoreLatePayloadAttributes,
		beaconClient:                beaconClient,
		blockTime:                   cfg.BlockTime,
	}

	builderBackend, err := NewBuilder(builderArgs)
	if err != nil {
		return fmt.Errorf("failed to create builder backend: %w", err)
	}
	builderService := NewService(cfg.ListenAddr, builderBackend)

	stack.RegisterAPIs([]rpc.API{
		{
			Namespace:     "builder",
			Version:       "1.0",
			Service:       builderService,
			Public:        true,
			Authenticated: true,
		},
	})

	stack.RegisterLifecycle(builderService)

	return nil
}
