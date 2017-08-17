package btc

import "os"

const (
	NetworkTypeEnvKey = "BTC_NETWORK_TYPE"
	NetworkTypeMain   = "mainnet"
	NetworkTypeTest   = "testnet"

	SatoshiPerBtc = 100000000

	TxVersion    int32  = 1
	TxLockTime   uint32 = 0
	TxInSequence uint32 = 4294967295

	CoinBaseTxid = "0000000000000000000000000000000000000000000000000000000000000000"
)

type Btc float64

func (btc Btc) Float64() float64 {
	return float64(btc)
}

func (btc Btc) Satoshi() Satoshi {
	return Satoshi(btc * SatoshiPerBtc)
}

type Satoshi int64

func (satoshi Satoshi) Int64() int64 {
	return int64(satoshi)
}

func (satoshi Satoshi) Btc() Btc {
	return Btc(float64(satoshi) / SatoshiPerBtc)
}

func UseTestnet() error {
	return os.Setenv(NetworkTypeEnvKey, NetworkTypeTest)
}

func isTestNet() bool {
	if os.Getenv(NetworkTypeEnvKey) == NetworkTypeTest {
		return true
	}

	return false
}
