package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"os"
	"yet/cli_api"
)

//go:embed "clo.json"
var cloBytes []byte

func main() {
	nod.EnableStdOut()

	bytesBuffer := bytes.NewBuffer(cloBytes)

	defs, err := clo.Load(bytesBuffer, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"get": cli_api.GetHandler,
	})

	if err := defs.Serve(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
