package discover

import (
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func TestBasicStorage(t *testing.T) {
	zeroNodeId := uint256.NewInt(0).Bytes32()
	storage, err := NewPortalStorage(math.MaxUint64, enode.ID(zeroNodeId), "./")
	defer storage.Close()
	assert.NoError(t, err)

	contentKey := []byte("test")
	content := []byte("value")

	_, err = storage.Get(contentKey, storage.ContentId(contentKey))
	assert.Equal(t, ContentNotFound, err)

	err = storage.Put(contentKey, content)
	assert.NoError(t, err)
}
