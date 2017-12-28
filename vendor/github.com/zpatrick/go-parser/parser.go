package parser

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
)

// ParseFiles parses files at the same time
func ParseFiles(paths []string) ([]*GoFile, error) {
	files := make([]*ast.File, len(paths))
	fsets := make([]*token.FileSet, len(paths))
	for i, p := range paths {
		// File: A File node represents a Go source file: https://golang.org/pkg/go/ast/#File
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			return nil, err
		}
		files[i] = file
		fsets[i] = fset
	}

	goFiles := make([]*GoFile, len(paths))
	for i, p := range paths {
		goFile, err := parseFile(p, files[i], fsets[i], files)
		if err != nil {
			return nil, err
		}
		goFiles[i] = goFile
	}
	return goFiles, nil
}

// ParseSingleFile parses a single file at the same time
func ParseSingleFile(path string) (*GoFile, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	return parseFile(path, file, fset, []*ast.File{file})
}

func parseFile(path string, file *ast.File, fset *token.FileSet, files []*ast.File) (*GoFile, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	if _, err = conf.Check(file.Name.Name, fset, files, info); err != nil {
		return nil, err
	}

	goFile := &GoFile{
		Path:    path,
		Package: file.Name.Name,
		Structs: []*GoStruct{},
	}

	// File.Decls: A list of the declarations in the file: https://golang.org/pkg/go/ast/#Decl
	for _, decl := range file.Decls {
		switch declType := decl.(type) {

		// GenDecl: represents an import, constant, type or variable declaration: https://golang.org/pkg/go/ast/#GenDecl
		case *ast.GenDecl:
			genDecl := declType

			// Specs: the Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec: https://golang.org/pkg/go/ast/#Spec
			for _, genSpec := range genDecl.Specs {
				switch genSpecType := genSpec.(type) {

				// TypeSpec: A TypeSpec node represents a type declaration: https://golang.org/pkg/go/ast/#TypeSpec
				case *ast.TypeSpec:
					typeSpec := genSpecType

					// typeSpec.Type: an Expr (expression) node: https://golang.org/pkg/go/ast/#Expr
					switch typeSpecType := typeSpec.Type.(type) {

					// StructType: A StructType node represents a struct type: https://golang.org/pkg/go/ast/#StructType
					case (*ast.StructType):
						structType := typeSpecType
						goStruct := buildGoStruct(source, goFile, info, typeSpec, structType)
						goFile.Structs = append(goFile.Structs, goStruct)
					// InterfaceType: An InterfaceType node represents an interface type. https://golang.org/pkg/go/ast/#InterfaceType
					case (*ast.InterfaceType):
						interfaceType := typeSpecType
						goInterface := buildGoInterface(source, goFile, info, typeSpec, interfaceType)
						goFile.Interfaces = append(goFile.Interfaces, goInterface)
					default:
						// a not-implemented typeSpec.Type.(type), ignore
					}
					// ImportSpec: An ImportSpec node represents a single package import. https://golang.org/pkg/go/ast/#ImportSpec
				case *ast.ImportSpec:
					importSpec := genSpec.(*ast.ImportSpec)
					goImport := buildGoImport(importSpec, goFile)
					goFile.Imports = append(goFile.Imports, goImport)
				default:
					// a not-implemented genSpec.(type), ignore
				}
			}
		case *ast.FuncDecl:
			funcDecl := declType
			goStructMethod := buildStructMethod(info, funcDecl, source)
			goFile.StructMethods = append(goFile.StructMethods, goStructMethod)

		default:
			// a not-implemented decl.(type), ignore
		}
	}

	return goFile, nil
}

func buildGoImport(spec *ast.ImportSpec, file *GoFile) *GoImport {
	name := ""
	if spec.Name != nil {
		name = spec.Name.Name
	}

	path := ""
	if spec.Path != nil {
		path = spec.Path.Value
	}

	return &GoImport{
		Name: name,
		Path: path,
		File: file,
	}
}

func buildGoInterface(source []byte, file *GoFile, info *types.Info, typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) *GoInterface {
	goInterface := &GoInterface{
		File:    file,
		Name:    typeSpec.Name.Name,
		Methods: buildMethodList(info, interfaceType.Methods.List, source),
	}

	return goInterface
}

func buildMethodList(info *types.Info, fieldList []*ast.Field, source []byte) []*GoMethod {
	methods := []*GoMethod{}

	for _, field := range fieldList {
		name := getNames(field)[0]

		fType, ok := field.Type.(*ast.FuncType)
		if !ok {
			// method was not a function
			continue
		}

		goMethod := &GoMethod{
			Name:    name,
			Params:  buildTypeList(info, fType.Params, source),
			Results: buildTypeList(info, fType.Results, source),
		}

		methods = append(methods, goMethod)
	}

	return methods
}

func buildStructMethod(info *types.Info, funcDecl *ast.FuncDecl, source []byte) *GoStructMethod {
	return &GoStructMethod{
		Receivers: buildReceiverList(info, funcDecl.Recv, source),
		GoMethod: GoMethod{
			Name:    funcDecl.Name.Name,
			Params:  buildTypeList(info, funcDecl.Type.Params, source),
			Results: buildTypeList(info, funcDecl.Type.Results, source),
		},
	}
}

func buildReceiverList(info *types.Info, fieldList *ast.FieldList, source []byte) []string {
	receivers := []string{}

	if fieldList != nil {
		for _, t := range fieldList.List {
			receivers = append(receivers, getTypeString(t.Type, source))
		}
	}

	return receivers
}

func buildTypeList(info *types.Info, fieldList *ast.FieldList, source []byte) []*GoType {
	types := []*GoType{}

	if fieldList != nil {
		for _, t := range fieldList.List {
			goType := buildType(info, t.Type, source)

			for _, n := range getNames(t) {
				copyType := copyType(goType)
				copyType.Name = n
				types = append(types, copyType)
			}
		}
	}

	return types
}

func getNames(field *ast.Field) []string {
	if field.Names == nil || len(field.Names) == 0 {
		return []string{""}
	}

	result := []string{}
	for _, name := range field.Names {
		result = append(result, name.String())
	}

	return result
}

func getTypeString(expr ast.Expr, source []byte) string {
	return string(source[expr.Pos()-1 : expr.End()-1])
}

func getUnderlyingTypeString(info *types.Info, expr ast.Expr) string {
	if typeInfo := info.TypeOf(expr); typeInfo != nil {
		if underlying := typeInfo.Underlying(); underlying != nil {
			return underlying.String()
		}
	}

	return ""
}

func copyType(goType *GoType) *GoType {
	return &GoType{
		Type:       goType.Type,
		Inner:      goType.Inner,
		Name:       goType.Name,
		Underlying: goType.Underlying,
	}
}

func buildType(info *types.Info, expr ast.Expr, source []byte) *GoType {
	innerTypes := []*GoType{}
	typeString := getTypeString(expr, source)
	underlyingString := getUnderlyingTypeString(info, expr)

	switch specType := expr.(type) {
	case *ast.FuncType:
		innerTypes = append(innerTypes, buildTypeList(info, specType.Params, source)...)
		innerTypes = append(innerTypes, buildTypeList(info, specType.Results, source)...)
	case *ast.ArrayType:
		innerTypes = append(innerTypes, buildType(info, specType.Elt, source))
	case *ast.StructType:
		innerTypes = append(innerTypes, buildTypeList(info, specType.Fields, source)...)
	case *ast.MapType:
		innerTypes = append(innerTypes, buildType(info, specType.Key, source))
		innerTypes = append(innerTypes, buildType(info, specType.Value, source))
	case *ast.ChanType:
		innerTypes = append(innerTypes, buildType(info, specType.Value, source))
	case *ast.StarExpr:
		innerTypes = append(innerTypes, buildType(info, specType.X, source))
	case *ast.Ellipsis:
		innerTypes = append(innerTypes, buildType(info, specType.Elt, source))
	case *ast.InterfaceType:
		methods := buildMethodList(info, specType.Methods.List, source)
		for _, m := range methods {
			innerTypes = append(innerTypes, m.Params...)
			innerTypes = append(innerTypes, m.Results...)
		}

	case *ast.Ident:
	case *ast.SelectorExpr:
	default:
		fmt.Printf("Unexpected field type: `%s`,\n %#v\n", typeString, specType)
	}

	return &GoType{
		Type:       typeString,
		Underlying: underlyingString,
		Inner:      innerTypes,
	}
}

func buildGoStruct(source []byte, file *GoFile, info *types.Info, typeSpec *ast.TypeSpec, structType *ast.StructType) *GoStruct {
	goStruct := &GoStruct{
		File:   file,
		Name:   typeSpec.Name.Name,
		Fields: []*GoField{},
	}

	// Field: A Field declaration list in a struct type, a method list in an interface type,
	// or a parameter/result declaration in a signature: https://golang.org/pkg/go/ast/#Field
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			goField := &GoField{
				Struct: goStruct,
				Name:   name.String(),
				Type:   string(source[field.Type.Pos()-1 : field.Type.End()-1]),
			}

			if field.Tag != nil {
				goTag := &GoTag{
					Field: goField,
					Value: field.Tag.Value,
				}

				goField.Tag = goTag
			}

			goStruct.Fields = append(goStruct.Fields, goField)
		}
	}

	return goStruct
}
