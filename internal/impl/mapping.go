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

	parser "github.com/zpatrick/go-parser"
)

var typeMaps = map[string]map[string]string{
	"[]byte": {
		"*big.Int": `new(big.Int).SetBytes({{ .Input }})`,
		"[32]byte": `BytesToBytes32({{ .Input }})`,
	},
	"string": {
		"common.Address": `common.HexToAddress({{ .Input }})`,
	},
	"*big.Int": {
		"[]byte": `{{ .Input }}.Bytes()`,
	},
	"[][]byte": {
		"[]*big.Int": `BytesToBigIntArray({{ .Input }})`,
	},
	"[]*big.Int": {
		"[][]byte": `BigIntArrayToBytes({{ .Input }})`,
	},
	"[32]byte": {
		"[]byte": `{{ .Input }}[:]`,
	},
	"common.Address": {
		"string": `{{ .Input }}.Hex()`,
	},
}

type TypeMap struct {
	Input    string
	Template string
}

func NewTypeMap(in, inType, outType string) *TypeMap {
	if inType == outType {
		return &TypeMap{
			Input:    in,
			Template: "{{ .Input }}",
		}
	}

	temp, ok := typeMaps[inType][outType]
	if !ok {
		return nil
	}
	return &TypeMap{
		Input:    in,
		Template: temp,
	}
}

func (t *TypeMap) String() string {
	implTemplate, err := template.New("type_map").Parse(t.Template)
	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		os.Exit(-1)
	}
	result := new(bytes.Buffer)
	implTemplate.Execute(result, t)
	return result.String()
}

func toRequestParam(f *parser.GoField, t *parser.GoType) string {
	typeMapping := NewTypeMap("r.Get"+f.Name+"()", f.Type, t.Type)
	if typeMapping == nil {
		panic("cannot find corresponding request type, from: " + f.Type + ", to: " + t.Type)
	}
	return typeMapping.String()
}

func toResponseParam(t *parser.GoType, f *parser.GoField) string {
	typeMapping := NewTypeMap("data", t.Type, f.Type)
	if typeMapping == nil {
		panic("cannot find corresponding response type, from: " + t.Type + ", to: " + f.Type)
	}
	return fmt.Sprintf("%v : %v", f.Name, typeMapping.String())
}
