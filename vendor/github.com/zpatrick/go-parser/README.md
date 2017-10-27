# Go Parser

The `go-parser` package is a small wrapper around the [parser](https://golang.org/pkg/go/parser/) and [ast](https://golang.org/pkg/go/ast/) packages. 

# Installation
To install this package, run:

`go get github.com/zpatrick/go-parser`

# Motivation
If you've ever tried to implement code generation in Go, you're probably familiar with the  `parser` and `ast` packages. 
These packages analyze and parse Go source files. 
The main problems I have with these packages are:
* They aren't very intuitive
* The amount of type assertions required to use them causes code to look cluttered and confusing

This package is meant to fix both of those issues.
For example, the following two snippets of code perform the same function:

**Using go-parser**
```
    goFile, err := parser.ParseFile("user.go")
    if err != nil {
        log.Fatal(err)
    }

    for _, goStruct := range goFile.Structs {
        for _, goField := range goStruct.Fields {
            log.Println(goField.Name, goField.Type, goField.Tag)
        }
    }
```

**Using `ast` and `parser`**
```
    src, err := ioutil.ReadFile("user.go")
    if err != nil {
        log.Fatal(err)
    }

    file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
    if err != nil {
        log.Fatal(err)
    }

    for _, d := range file.Decls {
        if typeDecl, ok := d.(*ast.GenDecl); ok {
            for _, s := range typeDecl.Specs {
                if typeSpec, ok := s.(*ast.TypeSpec); ok {
                    if structDecl, ok := typeSpec.Type.(*ast.StructType); ok {
                        for _, field := range structDecl.Fields.List {
                            for _, name := range field.Names {
                                name := name.String()
                                _type := string(src[field.Type.Pos()-1 : field.Type.End()-1])
                                tag := ""
                                if field.Tag != nil{
                                    tag = field.Tag.Value
                                }
                                
                                fmt.Printf(name, _type, tag)
                            }
                        }
                    }
                }
            }
        }
    }
```

# Disclaimer & Contributing
I have been adding coverage on an as-needed basis. As such, this package is very much incomplete relative to the amount of coverage available in the `ast` and `parser` packages. However, merge requests are **strongly** encouraged! I've added a lot of comments to the parsing code to hopefully make it easier to read and contribute to. 

# License
This work is published under the MIT license.

Please see the `LICENSE` file for details.