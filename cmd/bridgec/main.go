package main

import (
	"fmt"
	"os"

	"github.com/threefoldfoundation/tfchain/cmd/bridgec/internal"
	"github.com/threefoldtech/rivine/pkg/cli"
	"github.com/threefoldtech/rivine/pkg/daemon"
)

func main() {
	// create Bridge cli
	cliClient, err := internal.NewCommandLineClient("", "Bridge", daemon.RivineUserAgent)
	if err != nil {
		panic(err)
	}

	// register root command
	cliClient.ERC20Cmd = createERC20Cmd(cliClient)
	cliClient.RootCmd.AddCommand(cliClient.ERC20Cmd)

	// start cli
	if err := cliClient.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "client exited with an error: ", err)
		// Since no commands return errors (all commands set Command.Run instead of
		// Command.RunE), Command.Execute() should only return an error on an
		// invalid command or flag. Therefore Command.Usage() was called (assuming
		// Command.SilenceUsage is false) and we should exit with exitCodeUsage.
		os.Exit(cli.ExitCodeUsage)
	}
}
