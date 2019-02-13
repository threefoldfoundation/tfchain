package internal

import (
	"github.com/spf13/cobra"
	"github.com/threefoldtech/rivine/pkg/client"
)

// CommandLineClient extend for ERC20 commands
type CommandLineClient struct {
	*client.CommandLineClient

	ERC20Cmd *cobra.Command
}

func NewCommandLineClient(address, name, userAgent string) (*CommandLineClient, error) {
	client, err := client.NewCommandLineClient(address, name, userAgent)
	if err != nil {
		return nil, err
	}
	return &CommandLineClient{
		CommandLineClient: client,
	}, nil
}
