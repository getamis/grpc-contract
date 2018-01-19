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
	auth.GasLimit = big.NewInt(m.GasLimit)
	auth.GasPrice = big.NewInt(m.GasPrice)
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
		b = append(b, i.Bytes())
	}
	return
}

// BytesToBigIntArray converts [][]byte to []*big.Int
func BytesToBigIntArray(b [][]byte) (ints []*big.Int) {
	for _, i := range b {
		ints = append(ints, new(big.Int).SetBytes(i))
	}
	return
}

// BytesToBytes32 converts []byte to [32]byte
func BytesToBytes32(b []byte) (bs [32]byte) {
	copy(bs[:], b[:32])
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
