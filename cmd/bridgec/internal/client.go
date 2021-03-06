package internal

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/cli"
	"github.com/threefoldtech/rivine/pkg/client"
	rivinec "github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/types"
)

// Config defines the configuration for the bridge (CLI) client.
type Config struct {
	// These values aren't used for validation,
	// but only in order to estimate progress with the syncing of your consensus.
	BlockFrequencyInSeconds int64
	GenesisBlockTimestamp   types.Timestamp
}

// CommandLineClient specific for bridge client
type CommandLineClient struct {
	*api.HTTPClient

	Config     *client.Config
	RootCmd    *cobra.Command
	ERC20Cmd   *cobra.Command
	TFChainCmd *cobra.Command
	PreRunE    func(*client.Config) (*client.Config, error)
}

// NewCommandLineClient creates a new CLI client, which can be run as it is,
// or be extended/modified to fit your needs.
func NewCommandLineClient(address, name, userAgent string) (*CommandLineClient, error) {
	if address == "" {
		address = "http://localhost:23111"
	}
	if name == "" {
		name = "R?v?ne"
	}
	client := new(CommandLineClient)
	client.HTTPClient = &api.HTTPClient{
		RootURL:   address,
		UserAgent: userAgent,
	}

	createRootCmd(os.Args[0], name, client)

	// create command tree
	client.RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information.",
		Run: rivinec.Wrap(func() {
			fmt.Printf("Bride Client %s", config.GetBlockchainInfo().ChainVersion)

			fmt.Println()
			fmt.Printf("Go Version   v%s\r\n", runtime.Version()[2:])
			fmt.Printf("GOOS         %s\r\n", runtime.GOOS)
			fmt.Printf("GOARCH       %s\r\n", runtime.GOARCH)
		}),
	})
	client.RootCmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: fmt.Sprintf("Stop the %s bridge", name),
		Long:  fmt.Sprintf("Stop the %s bridge.", name),
		Run: rivinec.Wrap(func() {
			err := client.Post("/bridge/stop", "")
			if err != nil {
				cli.Die("Could not stop bridge:", err)
			}
			fmt.Println("bridge stopped.")
		}),
	})

	client.ERC20Cmd = CreateERC20Cmd(client)
	client.RootCmd.AddCommand(client.ERC20Cmd)

	client.TFChainCmd = createTFChainCommand(client)
	client.RootCmd.AddCommand(client.TFChainCmd)

	// parse flags
	client.RootCmd.PersistentFlags().StringVarP(&client.HTTPClient.RootURL, "addr", "a",
		client.HTTPClient.RootURL, fmt.Sprintf(
			"which host/port to communicate with (i.e. the host/port %sd is listening on)",
			name))

	// return client
	return client, nil
}

// Run the CLI, logic dependend upon the command the user used.
func (cli *CommandLineClient) Run() error {
	return cli.RootCmd.Execute()
}

// preRunE checks that all preConditions match
func (cli *CommandLineClient) preRunE(*cobra.Command, []string) error {
	address, err := sanitizeURL(cli.HTTPClient.RootURL)
	if err != nil {
		return fmt.Errorf("invalid daemon RPC address %q: %v", cli.HTTPClient.RootURL, err)
	}
	cli.HTTPClient.RootURL = address

	if cli.Config == nil {
		var err error
		cli.Config, err = client.FetchConfigFromDaemon(cli.HTTPClient)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetching config from daemon failed: %v\r\n", err)
		}
	}
	if cli.PreRunE != nil {
		cli.Config, err = cli.PreRunE(cli.Config)
		if err != nil {
			return fmt.Errorf("user-defined pre-run callback failed: %v", err)
		}
	}
	if cli.Config == nil {
		return errors.New("cannot run command line client: no config is defined")
	}
	return nil
}
