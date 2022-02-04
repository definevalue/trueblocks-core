package ipfs

import (
	shell "github.com/ipfs/go-ipfs-api"
)

var nodeShell *shell.Shell

// TODO: investigate how to do it; we need to pin to Pinata
func Connect(nodeUrl string) *shell.Shell {
	if nodeShell == nil {
		nodeShell = shell.NewShell(nodeUrl)
	}

	return nodeShell
}
