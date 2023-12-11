package discover

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

const nodeDataDir = "./"

func clear() {
	os.Remove(fmt.Sprintf("%s%s", nodeDataDir, sqliteName))
}

func genBytes(length int) []byte {
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = byte(i)
	}
	return res
}

func TestBasicStorage(t *testing.T) {
	zeroNodeId := uint256.NewInt(0).Bytes32()
	storage, err := NewPortalStorage(math.MaxUint32, enode.ID(zeroNodeId), nodeDataDir)
	assert.NoError(t, err)
	defer clear()
	defer storage.Close()

	contentKey := []byte("test")
	content := []byte("value")

	_, err = storage.Get(contentKey, storage.ContentId(contentKey))
	assert.Equal(t, ContentNotFound, err)

	pt := storage.Put(contentKey, content)
	assert.NoError(t, pt.Err())

	val, err := storage.Get(contentKey, storage.ContentId(contentKey))
	assert.NoError(t, err)
	assert.Equal(t, content, val)

	count, err := storage.ContentCount()
	assert.NoError(t, err)
	assert.Equal(t, count, uint64(1))

	size, err := storage.Size()
	assert.NoError(t, err)
	assert.True(t, size > 0)

	unusedSize, err := storage.UnusedSize()
	assert.NoError(t, err)

	usedSize, err := storage.UsedSize()
	assert.NoError(t, err)
	assert.True(t, usedSize == size-unusedSize)
}

func TestDBSize(t *testing.T) {
	zeroNodeId := uint256.NewInt(0).Bytes32()
	storage, err := NewPortalStorage(math.MaxUint32, enode.ID(zeroNodeId), nodeDataDir)
	assert.NoError(t, err)
	defer clear()
	defer storage.Close()

	numBytes := 10000

	size1, err := storage.Size()
	assert.NoError(t, err)
	putResult := storage.Put(uint256.NewInt(1).Bytes(), genBytes(numBytes))
	assert.Nil(t, putResult.Err())

	size2, err := storage.Size()
	assert.NoError(t, err)
	putResult = storage.Put(uint256.NewInt(2).Bytes(), genBytes(numBytes))
	assert.NoError(t, putResult.Err())

	size3, err := storage.Size()
	assert.NoError(t, err)
	putResult = storage.Put(uint256.NewInt(2).Bytes(), genBytes(numBytes))
	assert.NoError(t, putResult.Err())

	size4, err := storage.Size()
	assert.NoError(t, err)
	usedSize, err := storage.UsedSize()
	assert.NoError(t, err)

	assert.True(t, size2 > size1)
	assert.True(t, size3 > size2)
	assert.True(t, size4 == size3)
	assert.True(t, usedSize == size4)

	err = storage.del(storage.ContentId(uint256.NewInt(2).Bytes()))
	assert.NoError(t, err)
	err = storage.del(storage.ContentId(uint256.NewInt(1).Bytes()))
	assert.NoError(t, err)

	usedSize1, err := storage.UsedSize()
	assert.NoError(t, err)
	size5, err := storage.Size()
	assert.NoError(t, err)

	assert.True(t, size4 == size5)
	assert.True(t, usedSize1 < size5)

	err = storage.ReclaimSpace()
	assert.NoError(t, err)

	usedSize2, err := storage.UsedSize()
	assert.NoError(t, err)
	size6, err := storage.Size()
	assert.NoError(t, err)

	assert.Equal(t, size1, size6)
	assert.Equal(t, usedSize2, size6)

}

func TestDBPruning(t *testing.T) {
	storageCapacity := uint64(100_000)

	zeroNodeId := uint256.NewInt(0).Bytes32()
	storage, err := NewPortalStorage(storageCapacity, enode.ID(zeroNodeId), nodeDataDir)
	assert.NoError(t, err)
	defer clear()
	defer storage.Close()

	furthestElement := uint256.NewInt(40)
	secondFurthest := uint256.NewInt(30)
	thirdFurthest := uint256.NewInt(20)

	numBytes := 10_000
	// test with private put method
	pt1 := storage.put(uint256.NewInt(1).Bytes(), genBytes(numBytes))
	assert.NoError(t, pt1.Err())
	pt2 := storage.put(thirdFurthest.Bytes(), genBytes(numBytes))
	assert.NoError(t, pt2.Err())
	pt3 := storage.put(uint256.NewInt(3).Bytes(), genBytes(numBytes))
	assert.NoError(t, pt3.Err())
	pt4 := storage.put(uint256.NewInt(10).Bytes(), genBytes(numBytes))
	assert.NoError(t, pt4.Err())
	pt5 := storage.put(uint256.NewInt(5).Bytes(), genBytes(numBytes))
	assert.NoError(t, pt5.Err())
	pt6 := storage.put(uint256.NewInt(11).Bytes(), genBytes(numBytes))
	assert.NoError(t, pt6.Err())
	pt7 := storage.put(furthestElement.Bytes(), genBytes(4000))
	assert.NoError(t, pt7.Err())
	pt8 := storage.put(secondFurthest.Bytes(), genBytes(3000))
	assert.NoError(t, pt8.Err())
	pt9 := storage.put(uint256.NewInt(2).Bytes(), genBytes(numBytes))
	assert.NoError(t, pt9.Err())

	res, _ := storage.GetLargestDistance()

	assert.Equal(t, res, uint256.NewInt(40))
	pt10 := storage.put(uint256.NewInt(4).Bytes(), genBytes(12000))
	assert.NoError(t, pt10.Err())

	assert.False(t, pt1.Pruned())
	assert.False(t, pt2.Pruned())
	assert.False(t, pt3.Pruned())
	assert.False(t, pt4.Pruned())
	assert.False(t, pt5.Pruned())
	assert.False(t, pt6.Pruned())
	assert.False(t, pt7.Pruned())
	assert.False(t, pt8.Pruned())
	assert.False(t, pt9.Pruned())
	assert.True(t, pt10.Pruned())

	assert.Equal(t, pt10.PrunedCount(), 2)
	usedSize, err := storage.UsedSize()
	assert.NoError(t, err)
	assert.True(t, usedSize < storage.storageCapacityInBytes)

	_, err = storage.Get(furthestElement.Bytes(), storage.ContentId(furthestElement.Bytes()))
	assert.Equal(t, ContentNotFound, err)

	_, err = storage.Get(secondFurthest.Bytes(), storage.ContentId(secondFurthest.Bytes()))
	assert.Equal(t, ContentNotFound, err)

	val, err := storage.Get(thirdFurthest.Bytes(), thirdFurthest.Bytes())
	assert.NoError(t, err)
	assert.NotNil(t, val)

}

func TestGetLargestDistance(t *testing.T) {
	storageCapacity := uint64(100_000)

	zeroNodeId := uint256.NewInt(0).Bytes32()
	storage, err := NewPortalStorage(storageCapacity, enode.ID(zeroNodeId), nodeDataDir)
	assert.NoError(t, err)
	defer clear()
	defer storage.Close()

	furthestElement := uint256.NewInt(40)
	secondFurthest := uint256.NewInt(30)

	pt7 := storage.put(furthestElement.Bytes(), genBytes(2000))
	assert.NoError(t, pt7.Err())

	val, err := storage.Get(furthestElement.Bytes(), furthestElement.Bytes())
	assert.NoError(t, err)
	assert.NotNil(t, val)
	pt8 := storage.put(secondFurthest.Bytes(), genBytes(2000))
	assert.NoError(t, pt8.Err())
	res, err := storage.GetLargestDistance()
	assert.NoError(t, err)
	assert.Equal(t, furthestElement, res)
}

func TestForcePruning(t *testing.T) {
	const startCap = uint64(14_159_872) // 100KB
	const endCapacity = uint64(500_000) // 40KB
	const amountOfItems = 10_000

	maxUint256:= uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	nodeId := uint256.MustFromHex("0x30994892f3e4889d99deb5340050510d1842778acc7a7948adffa475fed51d6e").Bytes()
	content := genBytes(1000)

	storage, err := NewPortalStorage(startCap, enode.ID(nodeId), nodeDataDir)
	assert.NoError(t, err)

	increment := uint256.NewInt(0).Div(maxUint256, uint256.NewInt(amountOfItems))
	remainder := uint256.NewInt(0).Mod(maxUint256, uint256.NewInt(amountOfItems))

	id := uint256.NewInt(0)
	putCount := 0
	// id < maxUint256 - remainder
	for id.Cmp(uint256.NewInt(0).Sub(maxUint256, remainder)) == -1 {
		res := storage.put(id.Bytes(), content)
		assert.NoError(t, res.Err())
		id = id.Add(id, increment)
		putCount++
	}

	storage.storageCapacityInBytes = endCapacity

	oldDistance, err := storage.GetLargestDistance()
	assert.NoError(t, err)
	newDistance, err := storage.EstimateNewRadius(oldDistance)
	assert.NoError(t, err)
	assert.NotEqual(t, oldDistance.Cmp(newDistance), -1)
	err = storage.ForcePrune(newDistance)
	assert.NoError(t, err)

	err = storage.ReclaimSpace()
	assert.NoError(t, err)
	size, err := storage.Size()
	assert.NoError(t, err)

	diff := math.Abs(float64(size - storage.storageCapacityInBytes))
	// uint64(float64(storage.storageCapacityInBytes) * 0.2)
	assert.True(t, diff < float64(storage.storageCapacityInBytes) * 0.3)

}
