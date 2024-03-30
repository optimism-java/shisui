package beacon

//go:generate sszgen --path types.go

type LightClientUpdateKey struct {
	StartPeriod uint64
	Count       uint64
}

type LightClientBootstrapKey struct {
	BlockHash []byte `ssz-size:"32"`
}

type LightClientFinalityUpdateKey struct {
	FinalizedSlot uint64
}

type LightClientOptimisticUpdateKey struct {
	OptimisticSlot uint64
}
