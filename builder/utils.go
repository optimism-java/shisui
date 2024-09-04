package builder

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"time"

	builderTypes "github.com/ethereum/go-ethereum/builder/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type httpErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func handleError(w http.ResponseWriter, err error) {
	var errorMsg string
	var status int
	switch {
	case errors.Is(err, ErrIncorrectSlot):
		errorMsg = err.Error()
		status = http.StatusBadRequest
	case errors.Is(err, ErrNoPayloads):
		errorMsg = err.Error()
		status = http.StatusNotFound
	case errors.Is(err, ErrSlotFromPayload):
		errorMsg = err.Error()
		status = http.StatusInternalServerError
	case errors.Is(err, ErrSlotMismatch):
		errorMsg = err.Error()
		status = http.StatusBadRequest
	case errors.Is(err, ErrParentHashFromPayload):
		errorMsg = err.Error()
		status = http.StatusInternalServerError
	case errors.Is(err, ErrParentHashMismatch):
		errorMsg = err.Error()
		status = http.StatusBadRequest
	default:
		errorMsg = "error processing request"
		status = http.StatusInternalServerError
	}

	respondError(w, status, errorMsg)
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(httpErrorResp{code, message}); err != nil {
		http.Error(w, message, code)
	}
}

// runRetryLoop calls retry periodically with the provided interval respecting context cancellation
func runRetryLoop(ctx context.Context, interval time.Duration, retry func()) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			retry()
		}
	}
}

func SigningHash(domain [32]byte, chainID *big.Int, payloadBytes []byte) (common.Hash, error) {
	var msgInput [32 + 32 + 32]byte
	// domain: first 32 bytes
	copy(msgInput[:32], domain[:])
	// chain_id: second 32 bytes
	if chainID.BitLen() > 256 {
		return common.Hash{}, errors.New("chain_id is too large")
	}
	chainID.FillBytes(msgInput[32:64])
	// payload_hash: third 32 bytes, hash of encoded payload
	copy(msgInput[64:], crypto.Keccak256(payloadBytes))

	return crypto.Keccak256Hash(msgInput[:]), nil
}

func BuilderSigningHash(cfg *params.ChainConfig, payloadBytes []byte) (common.Hash, error) {
	return SigningHash(builderTypes.SigningDomainBuilderV1, cfg.ChainID, payloadBytes)
}
