package cli

import (
	"fmt"

	"github.com/convox/stdcli"
	"github.com/lob/rack/sdk"
)

func init() {
	registerWithoutProvider("update", "update the cli", Update, stdcli.CommandOptions{
		Flags:    []stdcli.Flag{flagRack},
		Validate: stdcli.ArgsMax(1),
	})
}

func Update(rack sdk.Interface, c *stdcli.Context) error {
	return fmt.Errorf("since this is a custom build, it should be updated manually")
}
