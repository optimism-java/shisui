package miner

import (
	"github.com/ethereum/go-ethereum/metrics"
)

var (
	blockProfitHistogram   = metrics.NewRegisteredHistogram("miner/block/profit", nil, metrics.NewExpDecaySample(1028, 0.015))
	blockProfitGauge       = metrics.NewRegisteredGauge("miner/block/profit/gauge", nil)
	culmulativeProfitGauge = metrics.NewRegisteredGauge("miner/block/profit/culmulative", nil)

	buildBlockTimer = metrics.NewRegisteredTimer("miner/block/build", nil)

	gasUsedGauge        = metrics.NewRegisteredGauge("miner/block/gasused", nil)
	transactionNumGauge = metrics.NewRegisteredGauge("miner/block/txnum", nil)
)
