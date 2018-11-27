//go:generate go get golang.org/x/tools/cmd/goyacc
//go:generate goyacc -o parser/parser.go -v parser/y.output parser/parser.y
//go:generate goyacc -o builtins_parser/builtins_parser.go -p builtins -v builtins_parser/y.output builtins_parser/builtins_parser.y
//go:generate go build -o lsl_ext
package main

import (
	"fmt"
	"os"

    "./builtins_parser"
    "./parser"
)

// TODO: add commang line arguments and options

func main() {
    builtins, err := builtins_parser.ParseBuiltins()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

    for _, arg := range os.Args[1:] {
        script, err := parser.ParseFile(arg)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        script.Builtins = builtins
        script.ConnectTree()

        if len(script.Errors) == 0 {
            fmt.Printf("%s\n\n", script)
        } else {
            for _, err := range script.Errors {
                fmt.Printf("%s %s\n", err.At, err.Error)
            }
            os.Exit(1)
        }
    }
}
