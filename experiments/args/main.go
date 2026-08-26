package main

import (
    "fmt"
    "os"
    "github.com/alexflint/go-arg"
)

type args struct {
    Yaml string `arg:"required,positional"`
}

func (args) Version() string{
    return "v0.1.0"
}

func main() {
    var args args
    arg.MustParse(&args)

    if _, err := os.Stat(args.Yaml); err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
        os.Exit(1)
    }

    fmt.Println(args.Yaml)
}
