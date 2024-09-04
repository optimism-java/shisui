package builder

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/builder/types"
	"gotest.tools/assert"
)

type mockBeaconNode struct {
	payloadAttributes types.PayloadAttributes
}

func (b *mockBeaconNode) Stop() {}

func (b *mockBeaconNode) SubscribeToPayloadAttributesEvents(payloadAttrC chan types.PayloadAttributes) {
	go func() {
		payloadAttrC <- b.payloadAttributes
	}()
}

func (b *mockBeaconNode) Start() error { return nil }

func newMockBeaconNode(payloadAttributes types.PayloadAttributes) *mockBeaconNode {
	return &mockBeaconNode{
		payloadAttributes: payloadAttributes,
	}
}

func TestMockBeaconNode(t *testing.T) {
	mockBeaconNode := newMockBeaconNode(types.PayloadAttributes{Slot: 123})
	payloadAttrC := make(chan types.PayloadAttributes, 1) // Use buffered channel

	mockBeaconNode.SubscribeToPayloadAttributesEvents(payloadAttrC)

	select {
	case payloadAttributes := <-payloadAttrC:
		assert.Equal(t, uint64(123), payloadAttributes.Slot)
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for payload attributes")
	}
}
