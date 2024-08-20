package builder

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

type mockBeaconNode struct {
	payloadAttributes BuilderPayloadAttributes
}

func (b *mockBeaconNode) Stop() {}

func (b *mockBeaconNode) SubscribeToPayloadAttributesEvents(payloadAttrC chan BuilderPayloadAttributes) {
	go func() {
		payloadAttrC <- b.payloadAttributes
	}()
}

func (b *mockBeaconNode) Start() error { return nil }

func newMockBeaconNode(payloadAttributes BuilderPayloadAttributes) *mockBeaconNode {
	return &mockBeaconNode{
		payloadAttributes: payloadAttributes,
	}
}

func TestMockBeaconNode(t *testing.T) {
	mockBeaconNode := newMockBeaconNode(BuilderPayloadAttributes{Slot: 123})
	payloadAttrC := make(chan BuilderPayloadAttributes, 1) // Use buffered channel

	mockBeaconNode.SubscribeToPayloadAttributesEvents(payloadAttrC)

	select {
	case payloadAttributes := <-payloadAttrC:
		assert.Equal(t, uint64(123), payloadAttributes.Slot)
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for payload attributes")
	}
}
