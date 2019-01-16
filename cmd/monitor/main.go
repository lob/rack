package main

import (
	"github.com/lob/rack/provider"
	"github.com/lob/rack/pkg/structs"
)

var (
	Provider structs.Provider
)

func init() {
	p, err := provider.FromEnv()
	if err != nil {
		panic(err)
	}
	Provider = p
}

func main() {
	Provider.Initialize(structs.ProviderOptions{})

	go Provider.Workers()

	select {}
}
