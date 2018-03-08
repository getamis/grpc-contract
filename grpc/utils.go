package grpc

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/getamis/sol2proto/pb"
)

type TransactOptsFn func(m *pb.TransactOpts) *bind.TransactOpts

// DefaultTransactOptsFn
func DefaultTransactOptsFn(m *pb.TransactOpts) *bind.TransactOpts {
	privateKey, err := crypto.ToECDSA(common.Hex2Bytes(m.PrivateKey))
	if err != nil {
		os.Exit(-1)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	if m.GasLimit < 0 {
		// get system suggested gas limit
		auth.GasLimit = 0
	} else {
		auth.GasLimit = uint64(m.GasLimit)
	}

	if m.GasPrice < 0 {
		// get system suggested gas price
		auth.GasPrice = nil
	} else {
		auth.GasPrice = big.NewInt(m.GasPrice)
	}

	if m.Nonce < 0 {
		// get system account nonce
		auth.Nonce = nil
	} else {
		auth.Nonce = big.NewInt(m.Nonce)
	}
	auth.Value = big.NewInt(m.Value)
	return auth
}

// BigIntArrayToBytes converts []*big.Int to [][]byte
func BigIntArrayToBytes(ints []*big.Int) (b [][]byte) {
	for _, i := range ints {
		if i == nil {
			b = append(b, nil)
		} else {
			b = append(b, i.Bytes())
		}
	}
	return
}

// BytesToBigIntArray converts [][]byte to []*big.Int
func BytesToBigIntArray(b [][]byte) (ints []*big.Int) {
	for _, i := range b {
		if i == nil {
			ints = append(ints, new(big.Int).SetInt64(0))
		} else {
			ints = append(ints, new(big.Int).SetBytes(i))
		}
	}
	return
}

// BytesToBytes32 converts []byte to [32]byte
func BytesToBytes32(b []byte) (bs [32]byte) {
	copyLen := len(b)
	if copyLen == 0 {
		return
	} else if copyLen > 32 {
		copyLen = 32
	}
	copy(bs[:], b[:copyLen])
	return
}

// BytesArrayToBytes32Array converts [][]byte to [][32]byte
func BytesArrayToBytes32Array(b [][]byte) (bs [][32]byte) {
	bs = make([][32]byte, len(b))
	for i := 0; i < len(b); i++ {
		bs[i] = BytesToBytes32(b[i])
	}
	return
}
