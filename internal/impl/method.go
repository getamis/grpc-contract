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
	"html/template"
	"strings"

	parser "github.com/zpatrick/go-parser"
)

type Method struct {
	StructName string
	Name       string
	InputType  string
	OutputType string

	ContractMethod *parser.GoStructMethod
	Request        *parser.GoStruct
	Response       *parser.GoStruct
}

var methodTemplate = `func (s *server) {{ .Name }}(c context.Context, r {{ .InputType }}) (*{{ .OutputType }}, error) {
	{{ PrintBody }}
}`

var constMethodBodyTemplate = `data, err := s.contract.{{ .Name }}(
		&bind.CallOpts{
			Pending: true,
			Context: c,
		},{{ PrintArgs }}
	)
	result := &{{ .OutputType }}{ {{ PrintOutputArgs }}
	}
	return result, err`

var methodBodyTemplate = `tx, err := s.contract.{{ .Name }}(
		r.GetOpts().TransactOpts(),{{ PrintArgs }}
	)
	if tx == nil {
		return nil, err
	}
	return &TransactionResp{
		Hash: tx.Hash().Hex(),
	}, err`

func NewMethod(m *parser.GoMethod, requestStruct *parser.GoStruct, responseStruct *parser.GoStruct, goFile *parser.GoFile) *Method {
	im := &Method{
		Name: m.Name,
	}

	im.InputType = m.Params[1].Type
	// Remove pointer star
	im.OutputType = m.Results[0].Type[1:]

	// Find the corresponding function in struct
	var goSM *parser.GoStructMethod

	// Build the method body
	if im.isConstant() {
		for _, sm := range goFile.StructMethods {
			if len(sm.Receivers) > 0 && strings.HasSuffix(sm.Receivers[0], "Caller") && sm.Name == im.Name {
				goSM = sm
			}
		}
	} else {
		for _, sm := range goFile.StructMethods {
			if len(sm.Receivers) > 0 && strings.HasSuffix(sm.Receivers[0], "Transactor") && sm.Name == im.Name {
				goSM = sm
			}
		}
	}

	if goSM == nil {
		fmt.Printf("Cannot find Method: %v\n", m.Name)
		return nil
	}

	im.ContractMethod = goSM
	im.Request = requestStruct
	im.Response = responseStruct
	return im
}

func (m Method) isConstant() bool {
	if m.OutputType == "TransactionResp" {
		return false
	}
	return true
}

func (m Method) String() string {
	tmpl, err := template.New("method").Funcs(template.FuncMap(
		map[string]interface{}{
			"PrintBody": func() (result string) {
				return m.printBody()
			},
		})).Parse(methodTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, m)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return result.String()
}

func (m Method) printBody() string {
	result := new(bytes.Buffer)

	tmpl := template.New("body")
	if m.isConstant() {
		tmpl, _ = tmpl.Funcs(template.FuncMap(
			map[string]interface{}{
				"PrintArgs": func() (result string) {
					args := ""
					for i := 0; i < len(m.Request.Fields); i++ {
						args += "\n\t\t" + toRequestParam(m.Request.Fields[i], m.ContractMethod.Params[i+1]) + ","
					}

					return args
				},
				"PrintOutputArgs": func() (result string) {
					args := ""
					for i := 0; i < len(m.Response.Fields); i++ {
						// TODO: may add nil protection
						args += "\n\t\t" + toResponseParam(m.ContractMethod.Results[i], m.Response.Fields[i]) + ","
					}
					return args
				},
			})).Parse(constMethodBodyTemplate)
		tmpl.Execute(result, m)
	} else {
		tmpl, _ = tmpl.New("body").Funcs(template.FuncMap(
			map[string]interface{}{
				"PrintArgs": func() (result string) {
					args := ""
					for i := 1; i < len(m.Request.Fields); i++ {
						args += "\n\t\t" + toRequestParam(m.Request.Fields[i], m.ContractMethod.Params[i]) + ","
					}
					return args
				},
			})).Parse(methodBodyTemplate)
		tmpl.Execute(result, m)
	}
	return result.String()
}
