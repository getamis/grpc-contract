// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package impl

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/getamis/grpc-contract/internal/util"
	"golang.org/x/tools/imports"
)

type Utils struct {
	Package string
}

var UtilsTemplate string = `package {{ .Package }};

type TransactOptsFn func(m *TransactOpts) *bind.TransactOpts

// defaultTransactOpts
func defaultTransactOptsFn(m *TransactOpts) *bind.TransactOpts {
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

// AnyToTransaction converts data to types.Transaction
func AnyToTransaction(data *any.Any) (*types.Transaction, error){
	tx := &types.Transaction{}
	err := rlp.DecodeBytes(data.Value, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
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
`

func (c *Utils) Write(filepath, filename string) {
	implTemplate, err := template.New("utils").Parse(UtilsTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		os.Exit(-1)
	}
	result := new(bytes.Buffer)
	err = implTemplate.Execute(result, c)
	if err != nil {
		fmt.Printf("Failed to render template: %v\n", err)
		os.Exit(-1)
	}
	code, err := imports.Process(".", result.Bytes(), nil)
	if err != nil {
		fmt.Printf("Failed to process code: %v\n", err)
		os.Exit(-1)
	}
	util.WriteFile(string(code), filepath, filename)
}