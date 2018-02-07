package grpc

import (
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigIntArrayToBytes(t *testing.T) {
	inputs := []*big.Int{
		nil,
		big.NewInt(1),
		big.NewInt(10000),
	}
	output := BigIntArrayToBytes(inputs)

	// the nil should be updated to big0
	inputs[0] = big.NewInt(0)
	exp := BytesToBigIntArray(output)
	assert.Equal(t, inputs, exp, "should be equal, got:%v, exp:%v", output, exp)
}

func TestBytesToBytes32(t *testing.T) {
	inputs := [][]byte{
		nil,
		[]byte(""),
		[]byte(strings.Repeat("1", 4)),
		[]byte(strings.Repeat("1", 33)),
	}
	for _, s := range inputs {
		output := BytesToBytes32(s)
		inStr := string(s)
		outStr := string(output[:])
		if len(inStr) > len(outStr) {
			assert.True(t, strings.HasPrefix(inStr, outStr), "%v should has prefix:%v", inStr, outStr)
		} else {
			assert.True(t, strings.HasPrefix(outStr, inStr), "%v should has prefix:%v", outStr, inStr)
		}
	}
}
