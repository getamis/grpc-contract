package impl

import (
	"fmt"

	parser "github.com/zpatrick/go-parser"
)

// TODO: Add more mapping
var requestParams = map[string]map[string]string{
	"byte": {
		"byte": "int8(r.%v)",
	},
	"[]byte": {
		"*big.Int": "new(big.Int).SetBytes(r.%v)",
	},
	"string": {
		"string":         "r.%v",
		"common.Address": "common.HexToAddress(r.%v)",
	},
	"bool": {
		"bool": "r.%v",
	},
}

// TODO: Add more mapping
var responseParams = map[string]map[string]string{
	"string": {
		"string": "%v: data",
	},
}

func toRequestParam(f *parser.GoField, t *parser.GoType) string {
	expr, ok := requestParams[f.Type][t.Type]
	if !ok {
		panic("cannot find corresponding request type, from: " + f.Type + ", to: " + t.Type)
	}
	return fmt.Sprintf(expr, f.Name)
}

func toResponseParam(f *parser.GoField, t *parser.GoType) string {
	expr, ok := responseParams[f.Type][t.Type]
	if !ok {
		panic("cannot find corresponding response type, from: " + f.Type + ", to: " + t.Type)
	}
	return fmt.Sprintf(expr, f.Name)
}
