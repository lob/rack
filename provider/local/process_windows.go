package local

import (
	"fmt"
	"io"

	"github.com/lob/rack/pkg/structs"
)

func (p *Provider) processExec(app, pid, command string, rw io.ReadWriter, opts structs.ProcessExecOptions) (int, error) {
	return -1, fmt.Errorf("unimplemented")
}

func (p *Provider) processRun(app, service string, opts processStartOptions) (string, error) {
	return "", fmt.Errorf("unimplemented")
}
